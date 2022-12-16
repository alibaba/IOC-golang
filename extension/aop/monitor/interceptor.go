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

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy:autoInjection=false

type interceptorImpl struct {
	monitorContext contextIOCInterface
}

func (w *interceptorImpl) BeforeInvoke(ctx *aop.InvocationContext) {
	if w.monitorContext != nil {
		w.monitorContext.BeforeInvoke(ctx)
	}
}

func (w *interceptorImpl) AfterInvoke(ctx *aop.InvocationContext) {
	if w.monitorContext != nil {
		w.monitorContext.AfterInvoke(ctx)
	}
}

func (w *interceptorImpl) Monitor(monitorCtx contextIOCInterface) {
	w.monitorContext = monitorCtx
}

func (w *interceptorImpl) StopMonitor() {
	w.monitorContext.Destroy()
	w.monitorContext = nil
}
