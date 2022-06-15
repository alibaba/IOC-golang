package common

import (
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func GetMethodUniqueKey(interfaceImplId, methodName string) string {
	return strings.Join([]string{interfaceImplId, methodName}, "-")
}

func ReflectValues2Strings(values []reflect.Value) []string {
	result := make([]string, 0)
	i := 0
	for ; i < len(values); i++ {
		if !values[i].IsValid() {
			result = append(result, "nil")
			continue
		}
		result = append(result, spew.Sdump(values[i].Interface()))
	}
	return result
}
