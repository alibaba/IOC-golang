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

package monitor

import (
	"github.com/alibaba/ioc-golang/aop"
)

type interceptorImpl struct {
	monitorContext *context
}

func (w *interceptorImpl) BeforeInvoke(ctx *aop.InvocationContext) {
	if w.monitorContext != nil {
		w.monitorContext.beforeInvoke(ctx)
	}
}

func (w *interceptorImpl) AfterInvoke(ctx *aop.InvocationContext) {
	if w.monitorContext != nil {
		w.monitorContext.afterInvoke(ctx)
	}
}

func (w *interceptorImpl) Monitor(monitorCtx *context) {
	w.monitorContext = monitorCtx
}

func (w *interceptorImpl) StopMonitor() {
	w.monitorContext.destroy()
	w.monitorContext = nil
}

var monitorInterceptorSingleton *interceptorImpl

func getMonitorInterceptorSingleton() *interceptorImpl {
	if monitorInterceptorSingleton == nil {
		monitorInterceptorSingleton = &interceptorImpl{}
	}
	return monitorInterceptorSingleton
}
