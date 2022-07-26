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

package autowire

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
	ServiceImpl1 Service `normal:"github.com/alibaba/ioc-golang/test/stress/autowire.ServiceImpl1"`

	// inject main.ServiceImpl2 pointer to Service interface with proxy wrapper
	ServiceImpl2 Service `normal:"github.com/alibaba/ioc-golang/test/stress/autowire.ServiceImpl2"`

	// inject ServiceImpl1 pointer to Service1 's own interface with proxy wrapper
	// this interface belongs to ServiceImpl1, there is no need to mark 'main.ServiceImpl1' in tag
	Service1OwnInterface ServiceImpl1IOCInterface `normal:""`

	// inject ServiceStruct struct pointer
	ServiceStruct *ServiceStruct `normal:""`
}

func (s *NormalApp) RunTest(t *testing.T) {
	// test creat by API
	createByAPIService1, err := GetServiceImpl1IOCInterface()
	assert.Nil(t, err)
	assert.Equal(t, expectString, createByAPIService1.GetHelloString(reqString))

	// test injection by API
	assert.Equal(t, expectString, s.Service1OwnInterface.GetHelloString(reqString))
	assert.Equal(t, expectString, s.ServiceStruct.GetString(reqString))
	assert.Equal(t, expectString, s.ServiceImpl1.GetHelloString(reqString))
	assert.Equal(t, expectString, s.ServiceImpl2.GetHelloString(reqString))
}
