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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"

	_ "github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/server/pkg/service"
)

func (a *App) TestRun(t *testing.T) {
	usr, err := a.ServiceStruct.GetUser("laurence", 23)
	assert.Nil(t, err)
	assert.NotNil(t, usr)
	assert.Equal(t, 1, usr.Id)
	assert.Equal(t, "laurence", usr.Name)
	assert.Equal(t, 23, usr.Age)
}

func TestRPCClient(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	app.TestRun(t)
}
