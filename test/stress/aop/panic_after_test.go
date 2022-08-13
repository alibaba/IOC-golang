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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/aop"
)

const Name = "afterIsCalledAfterPanicTestInterceptor"

func init() {
	// register a custom aop interceptor
	aop.RegisterAOP(aop.AOP{
		Name: Name,
		InterceptorFactory: func() aop.Interceptor {
			i, _ := GetafterIsCalledAfterPanicTestInterceptorIOCInterfaceSingleton()
			return i
		},
	})
}

func TestAOPAfterWithPanic(t *testing.T) {
	assert.Nil(t, ioc.Load())
	panicApp, err := GetPanicAfterCalledTestAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	result := panicApp.RunWithPanic("test")
	assert.Equal(t, "test", result)

	i, _ := GetafterIsCalledAfterPanicTestInterceptorIOCInterfaceSingleton()
	// the invocation have two layers, the internal layer is called after panic
	// assert the internal layer is not skipped because of panic.
	assert.Equal(t, 2, i.GetAfterIsCalledNum())
}
