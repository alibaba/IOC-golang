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
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
	aopLog "github.com/alibaba/ioc-golang/extension/aop/log"
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestAOPConcurrent(t *testing.T) {
	assert.Nil(t, ioc.Load(config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".invocation-aop-log.disable", true)))
	closeCh := make(chan struct{})
	go func() {
		output, err := iocli_command.Run([]string{"monitor"}, time.Second*6)
		assert.Nil(t, err)
		t.Log(output)
		assert.True(t, strings.Contains(output, `github.com/alibaba/ioc-golang/test/stress/aop.NormalApp.RunTest()
Total: 100000, Success: 100000, Fail: 0, AvgRT: `))
		assert.True(t, strings.Contains(output, `us, FailRate: 0.00%
github.com/alibaba/ioc-golang/test/stress/aop.ServiceImpl1.GetHelloString()
Total: 100000, Success: 100000, Fail: 0, AvgRT: `))
		close(closeCh)
	}()
	time.Sleep(time.Second * 1)
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				normalApp, err := GetNormalAppIOCInterface()
				assert.Nil(t, err)
				normalApp.RunTest(t)
			}
		}()
	}
	wg.Wait()
	<-closeCh
}

func TestAOPRecursive(t *testing.T) {
	assert.Nil(t, ioc.Load(config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".invocation-aop-log.disable", true)))
	closeCh := make(chan struct{})
	go func() {
		output, err := iocli_command.Run([]string{"monitor"}, time.Second*6)
		assert.Nil(t, err)
		t.Log(output)
		assert.True(t, strings.Contains(output, `github.com/alibaba/ioc-golang/test/stress/aop.RecursiveApp.RunTest()
Total: 901, Success: 901, Fail: 0, AvgRT: `))
		assert.True(t, strings.Contains(output, `us, FailRate: 0.00%
github.com/alibaba/ioc-golang/test/stress/aop.ServiceImpl1.GetHelloString()
Total: 2, Success: 2, Fail: 0, AvgRT: `))
		close(closeCh)
	}()
	time.Sleep(time.Second * 1)
	recApp, err := GetRecursiveAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	recApp.RunTest(t)
	<-closeCh

	recApp.Reset()
	closeCh = make(chan struct{})
	go func() {
		output, err := iocli_command.Run([]string{"trace"}, time.Second*6)
		assert.Nil(t, err)
		assert.Equal(t, 901, strings.Count(output, ", OperationName: github.com/alibaba/ioc-golang/test/stress/aop.(*recursiveApp_).RunTest, StartTime: "))
		close(closeCh)
	}()
	time.Sleep(time.Second * 1)
	recApp.RunTest(t)
	<-closeCh
}
