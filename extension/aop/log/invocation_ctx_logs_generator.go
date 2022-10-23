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
	"fmt"
	"strings"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:proxy=false

type InvocationCtxLogsGenerator struct {
}

func (i *InvocationCtxLogsGenerator) GetParamsLogs(dumpedParams []string, isCall bool) string {
	result := ""
	paramOrResponse := "Param"
	if !isCall {
		paramOrResponse = "Response"
	}

	for idx, p := range dumpedParams {
		result += fmt.Sprintf("%s %d: %s\n\n", paramOrResponse, idx+1, p)
	}
	return strings.TrimSuffix(result, "\n\n")
}

func (i *InvocationCtxLogsGenerator) GetFunctionSignatureLogs(sdid, methodName string, isCall bool) string {
	onToPrint := "Call"
	if !isCall {
		onToPrint = "Response"
	}
	return fmt.Sprintf("========== On %s ==========\n%s.%s()", onToPrint, sdid, methodName)
}
