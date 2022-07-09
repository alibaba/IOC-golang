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

package monitor

import (
	"github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

type monitorResponseItemsSorter []*monitor.MonitorResponseItem

func (m monitorResponseItemsSorter) Len() int {
	return len(m)
}

func (m monitorResponseItemsSorter) Less(i, j int) bool {
	return m[i].Sdid+m[i].Method < m[j].Sdid+m[j].Method
}

func (m monitorResponseItemsSorter) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
