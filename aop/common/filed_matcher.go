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
	"encoding/json"
	"reflect"
	"strings"
)

type FieldMatcher struct {
	FieldIndex int
	MatchRule  string // A.B.C=xxx
}

func (f *FieldMatcher) Match(values []reflect.Value) bool {
	if len(values) < f.FieldIndex {
		return false
	}
	targetVal := values[f.FieldIndex]
	data, err := json.Marshal(targetVal.Interface())
	if err != nil {
		return false
	}
	anyTypeMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &anyTypeMap); err != nil {
		return false
	}
	rules := strings.Split(f.MatchRule, "=")
	paths := rules[0]
	expectedValue := rules[1]
	splitedPaths := strings.Split(paths, ".")
	for i, p := range splitedPaths {
		subInterface, ok := anyTypeMap[p]
		if !ok {
			return false
		}
		if i == len(splitedPaths)-1 {
			// final must be string
			realStr, ok := subInterface.(string)
			if !ok {
				return false
			}
			if realStr != expectedValue {
				return false
			}
		} else {
			// not final, subInterface should be map[string]interface{}
			anyTypeMap, ok = subInterface.(map[string]interface{})
			if !ok {
				return false
			}
		}
	}
	return true
}
