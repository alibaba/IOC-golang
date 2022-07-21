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

package trace

import (
	"sort"
	"testing"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/stretchr/testify/assert"
)

func TestMonitorResponseItemsSorter(t *testing.T) {
	traceItems := make(traceSorter, 0)
	now := time.Now()
	traceItems = append(traceItems, &model.Trace{
		Spans: []*model.Span{
			{
				StartTime: now.Add(time.Millisecond),
			},
		},
	})
	traceItems = append(traceItems, &model.Trace{
		Spans: []*model.Span{
			{
				StartTime: now.Add(time.Millisecond * 2),
			},
		},
	})

	traceItems = append(traceItems, &model.Trace{
		Spans: []*model.Span{
			{
				StartTime: now,
			},
		},
	})

	sort.Sort(traceItems)
	assert.Equal(t, now.Unix(), traceItems[0].Spans[0].StartTime.Unix())
	assert.Equal(t, now.Add(time.Millisecond).Unix(), traceItems[1].Spans[0].StartTime.Unix())
	assert.Equal(t, now.Add(time.Millisecond*2).Unix(), traceItems[2].Spans[0].StartTime.Unix())
}
