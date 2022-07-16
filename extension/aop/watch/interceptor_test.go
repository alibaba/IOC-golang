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

package watch

import (
	oriCtx "context"
	"reflect"
	"testing"
	"time"

	"github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"

	"github.com/alibaba/ioc-golang/aop"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/autowire/util"
)

func TestWatchInterceptor(t *testing.T) {
	watchInterceptor := getWatchInterceptorSingleton()
	sdid := util.GetSDIDByStructPtr(&common.ServiceFoo{})
	methodName := "Invoke"
	methodFullName := sdid + "." + methodName
	sendCh := make(chan *watch.WatchResponse, 10)
	controlCh := make(chan *watch.WatchResponse, 10)
	go func() {
		info := <-sendCh
		controlCh <- info
	}()
	watchInterceptor.Watch(newContext(sdid, methodName, 0, 0, sendCh, nil))

	service := &common.ServiceFoo{}
	ctx := oriCtx.Background()
	param := &common.RequestParam{
		User: &common.User{
			Name: "laurence",
		},
	}

	invocationCtx := aop.NewInvocationContext(nil, sdid, methodName, methodFullName, []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)})
	watchInterceptor.BeforeInvoke(invocationCtx)
	rsp, err := service.Invoke(ctx, param)
	invocationCtx.SetReturnValues([]reflect.Value{reflect.ValueOf(rsp), reflect.ValueOf(err)})
	watchInterceptor.AfterInvoke(invocationCtx)
	info := <-controlCh
	assert.Equal(t, sdid, info.Sdid)
	assert.Equal(t, "Invoke", info.MethodName)
}

func TestWatchInterceptorWithCondition(t *testing.T) {
	watchInterceptor := getWatchInterceptorSingleton()
	sdid := util.GetSDIDByStructPtr(&common.ServiceFoo{})
	methodName := "Invoke"
	methodFullName := sdid + "." + methodName
	sendCh := make(chan *watch.WatchResponse, 10)
	controlCh := make(chan *watch.WatchResponse, 10)
	go func() {
		for {
			info := <-sendCh
			controlCh <- info
		}
	}()
	watchCtx := newContext(sdid, methodName, 0, 0, sendCh,
		&common.FieldMatcher{
			FieldIndex: 1,
			MatchRule:  "User.Name=lizhixin",
		})
	watchInterceptor.Watch(watchCtx)

	service := &common.ServiceFoo{}
	ctx := oriCtx.Background()

	// not match
	param := &common.RequestParam{
		User: &common.User{
			Name: "laurence",
		},
	}
	invocationCtx := aop.NewInvocationContext(nil, sdid, methodName, methodFullName, []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)})
	watchInterceptor.BeforeInvoke(invocationCtx)
	rsp, err := service.Invoke(ctx, param)
	info := &watch.WatchResponse{}
	time.Sleep(time.Millisecond * 500)
	invocationCtx.SetReturnValues([]reflect.Value{reflect.ValueOf(rsp), reflect.ValueOf(err)})
	watchInterceptor.AfterInvoke(invocationCtx)
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.Sdid)

	// match
	param.User.Name = "lizhixin"
	invocationCtx = aop.NewInvocationContext(nil, sdid, methodName, methodFullName, []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)})
	watchInterceptor.BeforeInvoke(invocationCtx)
	rsp, err = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	invocationCtx.SetReturnValues([]reflect.Value{reflect.ValueOf(rsp), reflect.ValueOf(err)})
	watchInterceptor.AfterInvoke(invocationCtx)
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, util.GetSDIDByStructPtr(&common.ServiceFoo{}), info.Sdid)

	// not watch
	param.User.Name = "lizhixin"
	watchInterceptor.UnWatch(watchCtx)
	invocationCtx = aop.NewInvocationContext(nil, sdid, methodName, methodFullName, []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(param)})
	watchInterceptor.BeforeInvoke(invocationCtx)
	rsp, err = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	invocationCtx.SetReturnValues([]reflect.Value{reflect.ValueOf(rsp), reflect.ValueOf(err)})
	watchInterceptor.AfterInvoke(invocationCtx)
	info = &watch.WatchResponse{}
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.Sdid)
}
