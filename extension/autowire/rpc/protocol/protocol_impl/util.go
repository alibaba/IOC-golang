package protocol_impl

import (
	"encoding/json"
	"reflect"
)

// ParseArgs @data should be []interface{}'s marshal result, @argsType should be each object's reflect type
func ParseArgs(argsType []reflect.Type, data []byte) ([]interface{}, error) {
	// http req -> invocation
	rawArguments := make([]interface{}, 0)
	finalArgument := make([]interface{}, 0)
	for _, reflectType := range argsType {
		if reflectType.Kind() == reflect.Ptr {
			rawArguments = append(rawArguments, reflect.New(reflectType.Elem()).Interface())
			finalArgument = append(finalArgument, reflect.New(reflectType.Elem()).Interface())
		} else {
			rawArguments = append(rawArguments, reflect.New(reflectType).Interface())
			finalArgument = append(finalArgument, reflect.New(reflectType).Interface())
		}
	}
	if err := json.Unmarshal(data, &rawArguments); err != nil {
		return nil, err
	}
	for idx, reflectType := range argsType {
		if reflectType.Kind() == reflect.Ptr {
			finalArgument[idx] = rawArguments[idx]
		} else {
			finalArgument[idx] = reflect.ValueOf(rawArguments[idx]).Elem().Interface()
		}
	}
	return finalArgument, nil
}
