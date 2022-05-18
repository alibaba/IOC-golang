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
	"testing"
	"time"

	"github.com/alibaba/ioc-golang/debug/common"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

type User struct {
	Name string
}

type RequestParam struct {
	User *User
}

type Response struct {
	Name string
}

type ServiceFoo struct {
}

func (s *ServiceFoo) Invoke(ctx context.Context, param *RequestParam) (*Response, error) {
	return &Response{
		Name: param.User.Name,
	}, nil
}

func TestWatchInterceptor(t *testing.T) {
	watchInterceptor := GetWatchInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan *boot.WatchResponse, 10)
	controlCh := make(chan *boot.WatchResponse, 10)
	go func() {
		info := <-sendCh
		controlCh <- info
		info = <-sendCh
		controlCh <- info
	}()

	watchInterceptor.Watch(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true), &WatchContext{
		Ch: sendCh,
	})
	watchInterceptor.Watch(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, false), &WatchContext{
		Ch: sendCh,
	})

	service := &ServiceFoo{}
	ctx := context.Background()
	param := &RequestParam{
		User: &User{
			Name: "laurence",
		},
	}

	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err := service.Invoke(ctx, param)
	info := <-controlCh
	assert.Equal(t, info.InterfaceName, "Service")
	assert.Equal(t, info.ImplementationName, "ServiceFoo")
	assert.Equal(t, info.MethodName, "Invoke")

	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, false),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	info = <-controlCh
	assert.Equal(t, "Service", info.InterfaceName)
	assert.Equal(t, "ServiceFoo", info.ImplementationName)
	assert.Equal(t, "Invoke", info.MethodName)
}

func TestWatchInterceptorWithCondition(t *testing.T) {
	watchInterceptor := GetWatchInterceptor()
	interfaceImplId := "Service-ServiceFoo"
	methodName := "Invoke"
	sendCh := make(chan *boot.WatchResponse, 10)
	controlCh := make(chan *boot.WatchResponse, 10)
	go func() {
		for {
			info := <-sendCh
			controlCh <- info
		}
	}()
	watchInterceptor.Watch(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true), &WatchContext{
		Ch: sendCh,
		FieldMatcher: &FieldMatcher{
			FieldIndex: 2,
			MatchRule:  "User.Name=lizhixin",
		},
	})

	service := &ServiceFoo{}
	ctx := context.Background()

	// not match
	param := &RequestParam{
		User: &User{
			Name: "laurence",
		},
	}
	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err := service.Invoke(ctx, param)
	info := &boot.WatchResponse{}
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.InterfaceName)
	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, false),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.InterfaceName)

	// match
	param.User.Name = "lizhixin"
	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	rsp, err = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "Service", info.InterfaceName)
	assert.Equal(t, "ServiceFoo", info.ImplementationName)
	assert.Equal(t, "Invoke", info.MethodName)
	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, false),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(rsp), reflect.ValueOf(err)})
	time.Sleep(time.Millisecond * 500)
	info = &boot.WatchResponse{}
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.InterfaceName)

	// not match
	param.User.Name = "lizhixin"
	watchInterceptor.UnWatch(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true))
	_, _ = watchInterceptor.Invoke(common.NewInterceptorContext(context.Background(), interfaceImplId, methodName, true),
		[]reflect.Value{reflect.ValueOf(service), reflect.ValueOf(ctx), reflect.ValueOf(param)})
	_, _ = service.Invoke(ctx, param)
	time.Sleep(time.Millisecond * 500)
	info = &boot.WatchResponse{}
	select {
	case info = <-controlCh:
	default:
	}
	assert.Equal(t, "", info.InterfaceName)
}
