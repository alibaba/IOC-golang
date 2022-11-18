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
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/example/aop/log/app"
	aopLog "github.com/alibaba/ioc-golang/extension/aop/log"
)

func TestGlobalLoggerReadonlyConfiguration(t *testing.T) {
	assert.Nil(t, ioc.Load(
		config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".global-logger-read-only", true),
		config.AddProperty(common.IOCGolangAOPConfigPrefix+"."+aopLog.Name+".level", "debug"),
	))
	_, err := app.GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	assert.Equal(t, logrus.InfoLevel, logrus.GetLevel())
}
