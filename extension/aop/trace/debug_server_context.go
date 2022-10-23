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
	traceCommon "github.com/alibaba/ioc-golang/extension/aop/trace/common"

	"github.com/alibaba/ioc-golang/aop/common"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=debugServerTraceByMethodContextParam
// +ioc:autowire:constructFunc=initDebugServerTraceByMethodContext
// +ioc:autowire:proxy=false

type debugServerTraceByMethodContext struct {
	methodName   string
	sdid         string
	fieldMatcher *common.FieldMatcher
	maxDepth     int64
	maxLength    int64
}

type debugServerTraceByMethodContextParam struct {
	sdid         string
	method       string
	fieldMatcher *common.FieldMatcher
	maxDepth     int64
	maxLength    int64
}

func (p *debugServerTraceByMethodContextParam) initDebugServerTraceByMethodContext(c *debugServerTraceByMethodContext) (*debugServerTraceByMethodContext, error) {
	if p.maxDepth == 0 {
		p.maxDepth = traceCommon.DefaultRecordValuesDepth
	}
	if p.maxLength == 0 {
		p.maxLength = traceCommon.DefaultRecordValuesLength
	}
	return &debugServerTraceByMethodContext{
		sdid:         p.sdid,
		methodName:   p.method,
		fieldMatcher: p.fieldMatcher,
		maxLength:    p.maxLength,
		maxDepth:     p.maxDepth,
	}, nil
}
