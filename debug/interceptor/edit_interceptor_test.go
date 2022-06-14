/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package interceptor

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

const serviceFooStructID = "github.com/alibaba/ioc-golang/debug/interceptore.ServiceFoo"

func TestEditInterceptorWithCondition(t *testing.T) {
	editInterceptor := GetEditInterceptor()
	interfaceImplId := serviceFooStructID
	methodName := "Invoke"
	sendCh := make(chan *boot.WatchResponse, 10)
	recvCh := make(chan *EditData, 10)
	controlSendCh := make(chan *boot.WatchResponse, 10)
	controlRecvCh := make(chan *EditData, 10)
	go func() {
		for {
			info := <-sendCh
			controlSendCh <- info
		}
	}()
	go func() {
		for {
			info := <-controlRecvCh
			recvCh <- info
		}
	}()
	editInterceptor.WatchEdit(interfaceImplId, methodName, true, &EditContext{
		SendCh: sendCh,
		RecvCh: recvCh,
		FieldMatcher: &FieldMatcher{
			FieldIndex: 1,
			MatchRule:  "User.Name=lizhixin",
		},
	})

	service := &ServiceFoo{}
	ctx := context.Background()

	param := &RequestParam{
		User: &User{
			Name: "lizhixin",
		},
	}
	go func() {
		controlRecvCh <- &EditData{
			FieldIndex: 1,
			FieldPath:  "User.Name",
			Value:      "laurence",
		}
	}()

	editInterceptor.Invoke(interfaceImplId, methodName, true,
		[]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)})

	rsp, err := service.Invoke(ctx, param)

	time.Sleep(time.Millisecond * 500)
	info := &boot.WatchResponse{}
	select {
	case info = <-controlSendCh:
	default:
	}
	assert.Equal(t, serviceFooStructID, info.ImplementationName)
	assert.Equal(t, "Invoke", info.MethodName)
	assert.Equal(t, true, info.IsParam)
	assert.True(t, strings.Contains(info.Params[1], "lizhixin"))

	assert.Nil(t, err)
	assert.Equal(t, "laurence", rsp.Name)
}
