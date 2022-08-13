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

package main

import (
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/example/aop/dynamic_plugin/complex/service2"
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func (a *App) TestWithPlugin(t *testing.T) {
	assert.Equal(t, "plugin", a.Service2OwnInterface.GetName())
	// default because plugin doesn't have construct function, which would be fixed soon
	assert.Equal(t, "default", a.Service2OwnInterface.GetService1Normal().GetName())
	// assert old singleton sub injection is reserved
	assert.Equal(t, "service1 singleton", a.Service2OwnInterface.GetService1Singleton().GetName())

	a.Service2OwnInterface.SetData("value")
	// set data bug is fixed by plugin, the set works successfully
	assert.Equal(t, "value", a.Service2OwnInterface.LoadData())
}

func (a *App) Init(t *testing.T) {
	// origin default name
	assert.Equal(t, "default", a.Service2OwnInterface.GetName())
	assert.Equal(t, "default", a.Service2OwnInterface.GetService1Normal().GetName())
	assert.Equal(t, "default", a.Service2OwnInterface.GetService1Singleton().GetName())

	// init name
	a.Service2OwnInterface.SetName("service2")
	a.Service2OwnInterface.GetService1Singleton().SetName("service1 singleton")
	a.Service2OwnInterface.GetService1Normal().SetName("service1 normal")
	a.Service2OwnInterface.SetData("value")

	// bug occurs as expected
	assert.Equal(t, "", a.Service2OwnInterface.LoadData())
}

func TestApp(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppSingleton()
	assert.Nil(t, err)

	app.Init(t)

	assert.Nil(t, exec.Command("go", "build", "-buildmode=plugin", "-o", "./plugin_1/service2.so", "./plugin_1").Run())
	sdid := util.GetSDIDByStructPtr(&service2.Service2{})
	output, err := iocli_command.Run([]string{"goplugin", "update", "singleton", sdid, "./plugin_1/service2.so", "Service2Plugin"}, time.Second)
	assert.Nil(t, err)
	time.Sleep(time.Second)
	assert.True(t, strings.Contains(output, "Update plugin success!"))

	app.TestWithPlugin(t)
}
