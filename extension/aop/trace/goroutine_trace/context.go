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

package goroutine_trace

import (
	"fmt"

	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/aop"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:proxy=false
// +ioc:autowire:constructFunc=initGoRoutineTracingContext
// +ioc:autowire:paramType=GoRoutineTracingContextParams

type GoRoutineTracingContext struct {
	entranceMethodFullName string
	facadeCtx              GoRoutineTracingFacadeContext
	grID                   int64
}

type GoRoutineTracingContextParams struct {
	EntranceMethodFullName string
	FacadeCtx              GoRoutineTracingFacadeContext
	// todo: add tracing depth field, to control trace limit depth,
}

func (p *GoRoutineTracingContextParams) initGoRoutineTracingContext(c *GoRoutineTracingContext) (*GoRoutineTracingContext, error) {
	if p.FacadeCtx == nil {
		return c, fmt.Errorf("failed to create GoRoutineTracingContext, param field FacadeCtx is nil")
	}
	if p.EntranceMethodFullName == "" {
		return c, fmt.Errorf("failed to create GoRoutineTracingContext, param field EntranceMethodFullName is empty")
	}
	c.entranceMethodFullName = p.EntranceMethodFullName
	c.facadeCtx = p.FacadeCtx
	c.grID = goid.Get()
	return c, nil
}

func (c *GoRoutineTracingContext) GetFacadeCtx() GoRoutineTracingFacadeContext {
	return c.facadeCtx
}

type GoRoutineTracingFacadeContext interface {
	aop.Interceptor
	Type() string
}
