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
	"reflect"
	"sync"

	"github.com/alibaba/ioc-golang/aop"

	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/aop/common"
	watchPB "github.com/alibaba/ioc-golang/extension/aop/watch/api/ioc_golang/aop/watch"
)

type interceptorImpl struct {
	watch sync.Map
}

func (w *interceptorImpl) BeforeInvoke(ctx *aop.InvocationContext) {
	if watchCtxInterface, ok := w.watch.Load(common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName)); ok {
		watchCtxInterface.(*context).beforeInvoke(ctx)
	}
}

func (w *interceptorImpl) AfterInvoke(ctx *aop.InvocationContext) {
	if watchCtxInterface, ok := w.watch.Load(common.GetMethodUniqueKey(ctx.SDID, ctx.MethodName)); ok {
		watchCtxInterface.(*context).afterInvoke(ctx)
	}
}

func (w *interceptorImpl) Watch(watchCtx *context) {
	w.watch.Store(common.GetMethodUniqueKey(watchCtx.SDID, watchCtx.MethodName), watchCtx)
}

func (w *interceptorImpl) UnWatch(watchCtx *context) {
	w.watch.Delete(common.GetMethodUniqueKey(watchCtx.SDID, watchCtx.MethodName))
}

type context struct {
	SDID              string
	MethodName        string
	Ch                chan *watchPB.WatchResponse
	FieldMatcher      *common.FieldMatcher
	watchGRRequestMap sync.Map
	maxDepth          int
	maxLength         int
}

func (w *context) beforeInvoke(ctx *aop.InvocationContext) {
	if w.FieldMatcher != nil && !w.FieldMatcher.Match(ctx.Params) {
		// doesn't match
		return
	}
	grID := goid.Get()
	w.watchGRRequestMap.Store(grID, ctx.Params)
}

func (w *context) afterInvoke(ctx *aop.InvocationContext) {
	paramValues, ok := w.watchGRRequestMap.Load(ctx.GrID)
	if !ok {
		return
	}
	invokeDetail := &watchPB.WatchResponse{
		Sdid:         w.SDID,
		MethodName:   w.MethodName,
		Params:       common.ReflectValues2Strings(paramValues.([]reflect.Value), w.maxDepth, w.maxLength),
		ReturnValues: common.ReflectValues2Strings(ctx.ReturnValues, w.maxDepth, w.maxLength),
	}
	w.Ch <- invokeDetail
	w.watchGRRequestMap.Delete(ctx.GrID)
}

var watchInterceptorSingleton *interceptorImpl

func getWatchInterceptorSingleton() *interceptorImpl {
	if watchInterceptorSingleton == nil {
		watchInterceptorSingleton = &interceptorImpl{}
	}
	return watchInterceptorSingleton
}
