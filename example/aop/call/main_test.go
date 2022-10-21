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
	"strings"
	"testing"
	"time"

	"github.com/alibaba/ioc-golang"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestCallAOP(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	go func() {
		app.Run()
	}()
	time.Sleep(time.Second * 1)

	output, err := iocli_command.Run([]string{"call", "singleton", "github.com/alibaba/ioc-golang/example/aop/call.UserService",
		`CreateUser`, `--params`, `["laurence",22]`}, time.Second*1)
	assert.Nil(t, err)
	t.Log(output)
	assert.True(t, strings.Contains(output, `Call singleton: github.com/alibaba/ioc-golang/example/aop/call.UserService.CreateUser() success!
Param = ["laurence",22]
Return values = [{"Id":1,"Name":"laurence","Age":22,"Mark":""},null]`))

	output, err = iocli_command.Run([]string{"call", "singleton", "github.com/alibaba/ioc-golang/example/aop/call.UserService",
		"ParseUserInfo", `--params`, `[{"Id":1,"Name":"laurence","Age":24}]`}, time.Second*1)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, `Call singleton: github.com/alibaba/ioc-golang/example/aop/call.UserService.ParseUserInfo() success!
Param = [{"Id":1,"Name":"laurence","Age":24}]
Return values = ["laurence",24,"",null]`))
}
