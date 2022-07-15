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
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func GetMethodUniqueKey(interfaceImplId, methodName string) string {
	return strings.Join([]string{interfaceImplId, methodName}, "-")
}

func ParseSDIDAndMethodFromUniqueKey(uniqueKey string) (string, string) {
	splitedUniqueKey := strings.Split(uniqueKey, "-")
	return strings.Join(splitedUniqueKey[:len(splitedUniqueKey)-1], "-"), splitedUniqueKey[len(splitedUniqueKey)-1]
}

func ReflectValues2String(values []reflect.Value, maxDepth, maxLength int) string {
	strs := ReflectValues2Strings(values, maxDepth, maxLength)
	return fmt.Sprintf("%+v", strs)
}

func ReflectValues2Strings(values []reflect.Value, maxDepth, maxLength int) []string {
	result := make([]string, 0)
	i := 0
	for ; i < len(values); i++ {
		if !values[i].IsValid() {
			result = append(result, "nil")
			continue
		}
		result = append(result, dumpSingletValue(values[i], maxDepth, maxLength))
	}
	return result
}

func dumpSingletValue(val reflect.Value, maxDepth, maxLength int) string {
	if !val.IsValid() {
		return "nil"
	}
	cfg := spew.NewDefaultConfig()
	cfg.DisablePointerAddresses = true
	cfg.MaxDepth = maxDepth
	cfg.SortKeys = true
	dumpedStr := cfg.Sdump(val.Interface())
	if len(dumpedStr) > maxLength {
		if maxDepth == 0 {
			return ""
		}
		dumpedStr = dumpSingletValue(val, maxDepth-1, maxLength)
	}
	return dumpedStr
}

func IsInvocationFailed(returnValues []reflect.Value) (bool, error) {
	if len(returnValues) == 0 {
		return false, nil
	}
	finalReturnValue := returnValues[len(returnValues)-1]
	if err, ok := finalReturnValue.Interface().(error); ok && err != nil {
		return true, err
	}
	return false, nil
}
