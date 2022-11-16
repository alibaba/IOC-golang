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

package test_default

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/example/aop/log/app"
	aopLog "github.com/alibaba/ioc-golang/extension/aop/log"
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestDefaultLogAOPLevelConfiguration(t *testing.T) {
	assert.Nil(t, ioc.Load())
	_, err := app.GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
	logInterceptorSingleton, _ := aopLog.GetlogInterceptorIOCInterfaceSingleton(nil)
	assert.Equal(t, logrus.InfoLevel, logInterceptorSingleton.GetInvocationCtxLogger().GetLevel())

	testLogCommand(t)
}

func testLogCommand(t *testing.T) {
	assert.Nil(t, ioc.Load())
	application, err := app.GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	go func() {
		application.Run()
	}()
	time.Sleep(time.Second * 1)

	output, err := iocli_command.Run([]string{"logs", "singleton", "github.com/alibaba/ioc-golang/example/aop/log/app.ServiceImpl1", "GetHelloString"}, time.Second*4)
	assert.Nil(t, err)
	assertContainsLevel(t, output, []string{"debug", "info", "warning", "error"}, "This is ServiceImpl1, hello laurence")
	assertContainsLevel(t, output, []string{"debug", "info", "warning", "error"}, "This is ServiceImpl2, hello laurence")

	output, err = iocli_command.Run([]string{"logs", "singleton", "github.com/alibaba/ioc-golang/example/aop/log/app.ServiceImpl1", "GetHelloString", "--level", "error"}, time.Second*4)
	assert.Nil(t, err)
	assertContainsLevel(t, output, []string{"error"}, "This is ServiceImpl1, hello laurence")
	assertContainsLevel(t, output, []string{"error"}, "This is ServiceImpl2, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug", "info", "warning"}, "This is ServiceImpl1, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug", "info", "warning"}, "This is ServiceImpl2, hello laurence")

	output, err = iocli_command.Run([]string{"logs", "singleton", "github.com/alibaba/ioc-golang/example/aop/log/app.ServiceImpl1", "GetHelloString", "--level", "info"}, time.Second*4)
	assert.Nil(t, err)
	assertContainsLevel(t, output, []string{"error", "warning", "info"}, "This is ServiceImpl1, hello laurence")
	assertContainsLevel(t, output, []string{"error", "warning", "info"}, "This is ServiceImpl2, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug"}, "This is ServiceImpl1, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug"}, "This is ServiceImpl2, hello laurence")

	output, err = iocli_command.Run([]string{"log", "singleton", "github.com/alibaba/ioc-golang/example/aop/log/app.ServiceImpl1", "GetHelloString", "--invocation"}, time.Second*4)
	assert.Nil(t, err)
	assertContainsLevel(t, output, []string{"debug", "info", "warning", "error"}, "This is ServiceImpl1, hello laurence")
	assertContainsLevel(t, output, []string{"debug", "info", "warning", "error"}, "This is ServiceImpl2, hello laurence")
	assertContainsInvocationCtx(t, output, `Response 1:`)
	assertContainsInvocationCtx(t, output, `Param 1:`)

	output, err = iocli_command.Run([]string{"log", "singleton", "github.com/alibaba/ioc-golang/example/aop/log/app.ServiceImpl1", "GetHelloString", "--level", "error", "--invocation"}, time.Second*4)
	assert.Nil(t, err)
	assertContainsLevel(t, output, []string{"error"}, "This is ServiceImpl1, hello laurence")
	assertContainsLevel(t, output, []string{"error"}, "This is ServiceImpl2, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug", "info", "warning"}, "This is ServiceImpl1, hello laurence")
	assertNotContainsLevel(t, output, []string{"debug", "info", "warning"}, "This is ServiceImpl2, hello laurence")
	assertContainsInvocationCtx(t, output, `Response 1`)
	assertContainsInvocationCtx(t, output, `Param 1`)

}

func assertContainsLevel(t *testing.T, output string, levels []string, msg string) {
	for _, level := range levels {
		assert.True(t, strings.Contains(output, fmt.Sprintf(`level=%s msg="%s"`, level, msg)))
	}
}

func assertNotContainsLevel(t *testing.T, output string, levels []string, msg string) {
	for _, level := range levels {
		assert.False(t, strings.Contains(output, fmt.Sprintf(`level=%s msg="%s"`, level, msg)))
	}
}

func assertContainsInvocationCtx(t *testing.T, output, msg string) {
	assert.True(t, strings.Contains(output, "[AOP Function Call] ========== On Call =========="))
	assert.True(t, strings.Contains(output, "[AOP Function Response] ========== On Response =========="))
	assert.True(t, strings.Contains(output, msg))
}
