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
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/extension/aop/monitor/api/ioc_golang/aop/monitor"
)

func TestMonitorResponseItemsSorter(t *testing.T) {
	monitorResponseItems := make(monitorResponseItemsSorter, 0)
	monitorResponseItems = append(monitorResponseItems, &monitor.MonitorResponseItem{
		Sdid:   "github.com/alibaba/ioc-golang/test.B",
		Method: "method2",
	})
	monitorResponseItems = append(monitorResponseItems, &monitor.MonitorResponseItem{
		Sdid:   "github.com/alibaba/ioc-golang/test.B",
		Method: "method1",
	})
	monitorResponseItems = append(monitorResponseItems, &monitor.MonitorResponseItem{
		Sdid:   "github.com/alibaba/ioc-golang/test.A",
		Method: "method1",
	})

	sort.Sort(monitorResponseItems)
	assert.Equal(t, "github.com/alibaba/ioc-golang/test.A", monitorResponseItems[0].Sdid)
	assert.Equal(t, "method1", monitorResponseItems[0].Method)
	assert.Equal(t, "github.com/alibaba/ioc-golang/test.B", monitorResponseItems[1].Sdid)
	assert.Equal(t, "method1", monitorResponseItems[1].Method)
	assert.Equal(t, "github.com/alibaba/ioc-golang/test.B", monitorResponseItems[2].Sdid)
	assert.Equal(t, "method2", monitorResponseItems[2].Method)
}
