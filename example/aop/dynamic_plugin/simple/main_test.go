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
	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func (a *App) TestRun(t *testing.T) {
	assert.Equal(t, "This is Plugin, hello laurence", a.Service1OwnInterface.GetHelloString("laurence"))
}

func TestApp(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppSingleton()
	assert.Nil(t, err)
	assert.Nil(t, exec.Command("go", "build", "-buildmode=plugin", "-o", "./plugin_1/service.so", "./plugin_1").Run())
	sdid := util.GetSDIDByStructPtr(&ServiceImpl1{})
	output, err := iocli_command.Run([]string{"goplugin", "update", "singleton", sdid, "./plugin_1/service.so", "Service1Plugin"}, time.Second)
	assert.Nil(t, err)
	time.Sleep(time.Second)
	assert.True(t, strings.Contains(output, "Update plugin success!"))
	app.TestRun(t)
}
