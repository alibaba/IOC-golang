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

package call

import (
	"github.com/alibaba/ioc-golang/aop/common"
	logPB "github.com/alibaba/ioc-golang/extension/aop/log/api/ioc_golang/aop/log"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=debugLogContextParam
// +ioc:autowire:constructFunc=init
// +ioc:autowire:proxy=false

type debugLogContext struct {
	debugLogContextParam
}

type debugLogContextParam struct {
	sdid         string
	methodName   string
	ch           chan *logPB.LogResponse
	fieldMatcher *common.FieldMatcher
	//traceEnable  bool
}

func (p *debugLogContextParam) init(c *debugLogContext) (*debugLogContext, error) {
	c.debugLogContextParam = *p
	return c, nil
}
