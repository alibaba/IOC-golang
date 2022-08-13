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

package aop

import (
	"fmt"

	"github.com/alibaba/ioc-golang/aop"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=initAfterIsCalledAfterPanicTestInterceptor

type afterIsCalledAfterPanicTestInterceptor struct {
	AfterIsCalledNum int
}

func initAfterIsCalledAfterPanicTestInterceptor(a *afterIsCalledAfterPanicTestInterceptor) (*afterIsCalledAfterPanicTestInterceptor, error) {
	a.AfterIsCalledNum = 0
	return a, nil
}

func (p *afterIsCalledAfterPanicTestInterceptor) GetAfterIsCalledNum() int {
	return p.AfterIsCalledNum
}

func (p *afterIsCalledAfterPanicTestInterceptor) BeforeInvoke(ctx *aop.InvocationContext) {

}

func (p *afterIsCalledAfterPanicTestInterceptor) AfterInvoke(ctx *aop.InvocationContext) {
	p.AfterIsCalledNum++
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type PanicAfterCalledTestSubApp struct {
}

func (p *PanicAfterCalledTestSubApp) RunWithPanic(panicMsg string) {
	panic(panicMsg)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type PanicAfterCalledTestApp struct {
	PanicAfterCalledTestSubApp PanicAfterCalledTestSubAppIOCInterface `singleton:""`
}

func (p *PanicAfterCalledTestApp) RunWithPanic(panicMst string) (result string) {
	defer func() {
		if r := recover(); r != nil {
			result = fmt.Sprintf("%+v", r)
		}
	}()
	p.PanicAfterCalledTestSubApp.RunWithPanic(panicMst)
	return ""
}
