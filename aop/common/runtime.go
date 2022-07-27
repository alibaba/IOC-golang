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

package common

import (
	"runtime"
	"strings"
)

const (
	ProxyMethodPrefix = "github.com/alibaba/ioc-golang/aop."
)

func CurrentCallingMethodName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip, pc)
	return runtime.FuncForPC(pc[0]).Name()
}

func IsTraceEntrance(entranceName string) bool {
	pc := make([]uintptr, 500)
	n := runtime.Callers(0, pc)
	foundEntrance := false
	level := int64(0)

	for i := n - 1; i >= 0; i-- {
		fName := runtime.FuncForPC(pc[i]).Name()
		if foundEntrance {
			if strings.HasPrefix(fName, ProxyMethodPrefix) {
				level++
			}
			if level == 2 {
				return false
			}
			continue
		}
		if fName == entranceName {
			foundEntrance = true
		}
	}

	return level-1 == 0
}
