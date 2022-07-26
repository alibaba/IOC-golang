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
)

const reqString = "laurence"
const expectString = "hello laurence"

// +ioc:autowire=true
// +ioc:autowire:type=normal

type NormalApp struct {
	// inject main.ServiceImpl1 pointer to Service interface with proxy wrapper
	ServiceImpl1 Service `normal:"github.com/alibaba/ioc-golang/test/stress/aop.ServiceImpl1"`
}

func (s *NormalApp) RunTest(t *testing.T) {
	// test creat by API
	createByAPIService1, err := GetServiceImpl1IOCInterface()
	assert.Nil(t, err)
	assert.Equal(t, expectString, createByAPIService1.GetHelloString(reqString))
}
