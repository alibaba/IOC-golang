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
