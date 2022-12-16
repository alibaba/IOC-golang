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
	"fmt"
	"testing"

	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/example/autowire/autowire_allimpls/service"
	"github.com/alibaba/ioc-golang/extension/autowire/allimpls"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
)

func (a *App) TestRun(t *testing.T) {
	outputMap := make(map[string]struct{})
	assert.Equal(t, 4, len(a.ServiceImpls))
	for _, s := range a.ServiceImpls {
		outputMap[s.GetHelloString("laurence")] = struct{}{}
	}
	_, ok := outputMap["This is ServiceImpl1, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl2, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl3, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl4, hello laurence"]
	assert.True(t, ok)
}

func TestAutowireAllImplsInjection(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetAppSingleton()
	assert.Nil(t, err)
	app.TestRun(t)
}

func TestAutowireAllImplsAPI(t *testing.T) {
	assert.Nil(t, ioc.Load())
	outputMap := make(map[string]struct{})
	allServiceImpls, err := allimpls.GetImpl(util.GetSDIDByStructPtr(new(service.Service)))
	if err != nil {
		panic(err)
	}
	assert.Equal(t, 4, len(allServiceImpls.([]service.Service)))
	for _, s := range allServiceImpls.([]service.Service) {
		fmt.Println(s.GetHelloString("laurence"))
		outputMap[s.GetHelloString("laurence")] = struct{}{}
	}
	_, ok := outputMap["This is ServiceImpl1, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl2, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl3, hello laurence"]
	assert.True(t, ok)
	_, ok = outputMap["This is ServiceImpl4, hello laurence"]
	assert.True(t, ok)
	assert.Nil(t, err)
}
