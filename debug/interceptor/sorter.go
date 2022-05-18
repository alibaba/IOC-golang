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

package interceptor

import (
	"github.com/alibaba/ioc-golang/debug/api/ioc_golang/boot"
)

type MetadataSorter []*boot.ServiceMetadata

func (m MetadataSorter) Len() int {
	return len(m)
}

func (m MetadataSorter) Less(i, j int) bool {
	return m[i].InterfaceName+m[i].ImplementationName < m[j].InterfaceName+m[j].ImplementationName
}

func (m MetadataSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type MethodSorter []string

func (m MethodSorter) Len() int {
	return len(m)
}

func (m MethodSorter) Less(i, j int) bool {
	return m[i] < m[j]
}

func (m MethodSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
