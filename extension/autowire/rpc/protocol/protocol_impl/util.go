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
			//switch reflectType.Elem().Kind() {
			//case reflect.Int:
			//	num := int(rawArguments[idx].(float64))
			//	finalArgument[idx] = &num
			//case reflect.Int32:
			//	num := int32(rawArguments[idx].(float64))
			//	finalArgument[idx] = &num
			//case reflect.Int64:
			//	num := int64(rawArguments[idx].(float64))
			//	finalArgument[idx] = &num
			//case reflect.Float32:
			//	num := float32(rawArguments[idx].(float64))
			//	finalArgument[idx] = &num
			//case reflect.Float64:
			//	num := rawArguments[idx].(float64)
			//	finalArgument[idx] = &num
			//default:
			//	finalArgument[idx] = rawArguments[idx]
			//}
			finalArgument[idx] = rawArguments[idx]
		} else {
			switch reflectType.Kind() {
			case reflect.Int:
				num := rawArguments[idx].(*int)
				finalArgument[idx] = *num
			case reflect.Int32:
				num := rawArguments[idx].(*int32)
				finalArgument[idx] = *num
			case reflect.Int64:
				num := rawArguments[idx].(*int64)
				finalArgument[idx] = *num
			case reflect.Float32:
				num := rawArguments[idx].(*float32)
				finalArgument[idx] = *num
			case reflect.Float64:
				num := rawArguments[idx].(*float64)
				finalArgument[idx] = *num
			case reflect.String:
				num := rawArguments[idx].(*string)
				finalArgument[idx] = *num
			default:
				finalArgument[idx] = rawArguments[idx]
			}
		}
	}
	return finalArgument, nil
}
