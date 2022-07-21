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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/extension/aop/list/api/ioc_golang/aop/list"
)

func TestMetadataSorter(t *testing.T) {
	structsMetadatas := make(metadataSorter, 0)
	structsMetadatas = append(structsMetadatas, &list.ServiceMetadata{
		InterfaceName:      "interface2",
		ImplementationName: "impl2",
	})
	structsMetadatas = append(structsMetadatas, &list.ServiceMetadata{
		InterfaceName:      "interface2",
		ImplementationName: "impl1",
	})
	structsMetadatas = append(structsMetadatas, &list.ServiceMetadata{
		InterfaceName:      "interface1",
		ImplementationName: "impl1",
	})
	sort.Sort(structsMetadatas)
	assert.Equal(t, "interface1", structsMetadatas[0].InterfaceName)
	assert.Equal(t, "impl1", structsMetadatas[0].ImplementationName)
	assert.Equal(t, "interface2", structsMetadatas[1].InterfaceName)
	assert.Equal(t, "impl1", structsMetadatas[1].ImplementationName)
	assert.Equal(t, "interface2", structsMetadatas[2].InterfaceName)
	assert.Equal(t, "impl2", structsMetadatas[2].ImplementationName)
}

func TestMethodSorter(t *testing.T) {
	methodMetadatas := make(methodSorter, 0)
	methodMetadatas = append(methodMetadatas, "method1")
	methodMetadatas = append(methodMetadatas, "method2")
	methodMetadatas = append(methodMetadatas, "method3")
	sort.Sort(methodMetadatas)
	assert.Equal(t, "method1", methodMetadatas[0])
	assert.Equal(t, "method2", methodMetadatas[1])
	assert.Equal(t, "method3", methodMetadatas[2])
}
