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

	"github.com/alibaba/ioc-golang/config"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
)

func (a *App) TestRun(t *testing.T, assertProfile string) {
	assert.NotNil(t, a.ServiceImpl)
	assert.Equal(t, "This is service"+assertProfile+"Impl, hello laurence", a.ServiceImpl.GetHelloString("laurence"))
}

func TestAutowireProfileActiveDevInjection(t *testing.T) {
	assert.Nil(t, ioc.Load(config.WithProfilesActive("dev")))
	app, err := GetApp()
	assert.Nil(t, err)
	app.TestRun(t, "Dev")
}

func TestAutowireProfileActiveProAndDevInjection(t *testing.T) {
	assert.Nil(t, ioc.Load(config.WithProfilesActive("dev", "pro")))
	app, err := GetApp()
	assert.Nil(t, err)
	app.TestRun(t, "Pro")
}

func TestAutowireProfileNoneExistingFallbackInjection(t *testing.T) {
	assert.Nil(t, ioc.Load(config.WithProfilesActive("bad")))
	app, err := GetApp()
	assert.Nil(t, err)
	app.TestRun(t, "Default")
}

func TestAutowireProfileEmptyFallbackInjection(t *testing.T) {
	assert.Nil(t, ioc.Load())
	app, err := GetApp()
	assert.Nil(t, err)
	app.TestRun(t, "Default")
}
