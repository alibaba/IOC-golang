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

package list

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/aop/common"
)

func TestListServiceImpl(t *testing.T) {
	debugMetadata := aop.GetAllInterfaceMetadata()
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct1"] = &common.StructMetadata{
		MethodMetadata: map[string]*common.MethodMetadata{},
	}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct2"] = &common.StructMetadata{
		MethodMetadata: map[string]*common.MethodMetadata{},
	}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct1"].MethodMetadata["MockMethod1"] = &common.MethodMetadata{}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct1"].MethodMetadata["MockMethod2"] = &common.MethodMetadata{}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct1"].MethodMetadata["MockMethod3"] = &common.MethodMetadata{}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct2"].MethodMetadata["MockMethod1"] = &common.MethodMetadata{}
	debugMetadata["github.com/alibaba/ioc-golang/aop/test.Struct2"].MethodMetadata["MockMethod2"] = &common.MethodMetadata{}
	mockService := newMockServiceImpl(debugMetadata)
	result, err := mockService.List(nil, nil)
	assert.Nil(t, err)
	serviceMetadatas := result.GetServiceMetadata()
	assert.NotNil(t, serviceMetadatas)
	assert.Equal(t, 2, len(serviceMetadatas))
	assert.Equal(t, "github.com/alibaba/ioc-golang/aop/test.Struct1", serviceMetadatas[0].ImplementationName)
	assert.Equal(t, 3, len(serviceMetadatas[0].Methods))
	assert.Equal(t, "MockMethod1", serviceMetadatas[0].Methods[0])
	assert.Equal(t, "MockMethod2", serviceMetadatas[0].Methods[1])
	assert.Equal(t, "MockMethod3", serviceMetadatas[0].Methods[2])

	assert.Equal(t, "github.com/alibaba/ioc-golang/aop/test.Struct2", serviceMetadatas[1].ImplementationName)
	assert.Equal(t, 2, len(serviceMetadatas[1].Methods))
	assert.Equal(t, "MockMethod1", serviceMetadatas[1].Methods[0])
	assert.Equal(t, "MockMethod2", serviceMetadatas[1].Methods[1])
}
