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

package util

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	emptyString = ""
)

func GetStructName(v interface{}) string {
	if v == nil {
		return ""
	}
	typeOfInterface := GetTypeFromInterface(v)
	return typeOfInterface.Name()
}

func GetSDIDByStructPtr(v interface{}) string {
	if v == nil {
		return ""
	}
	return GetSDIDByReflectType(GetTypeFromInterface(v))
}

func GetSDIDByReflectType(typeOfInterface reflect.Type) string {
	return fmt.Sprintf("%s.%s", typeOfInterface.PkgPath(), typeOfInterface.Name())
}

func IsProxyStructPtr(v interface{}) bool {
	if v == nil {
		return false
	}
	typeOfInterface := GetTypeFromInterface(v)
	return strings.HasSuffix(typeOfInterface.Name(), "_")
}

func GetProxySDIDByStructPtr(v interface{}) string {
	if v == nil {
		return ""
	}
	typeOfInterface := GetTypeFromInterface(v)
	return fmt.Sprintf("%s.%s_", typeOfInterface.PkgPath(), strings.ToLower(string(typeOfInterface.Name()[0]))+typeOfInterface.Name()[1:])
}

func GetTypeFromInterface(v interface{}) reflect.Type {
	valueOfInterface := reflect.ValueOf(v)
	valueOfElemInterface := valueOfInterface.Elem()
	return valueOfElemInterface.Type()
}

func ToCamelCase(src string) string {
	if src == emptyString {
		return emptyString
	}

	return strings.ToLower(src[:1]) + src[1:]
}

func ToSnakeCase(src string) string {
	if src == emptyString {
		return src
	}
	srcLen := len(src)
	result := make([]byte, 0, srcLen*2)
	caseSymbol := false
	for i := 0; i < srcLen; i++ {
		char := src[i]
		if i > 0 && char >= 'A' && char <= 'Z' && caseSymbol { // _xxx || yyy__zzz
			result = append(result, '_')
		}
		caseSymbol = char != '_'

		result = append(result, char)
	}

	snakeCase := strings.ToLower(string(result))

	return snakeCase
}

func IsPointerField(fieldType reflect.Type) bool {
	return fieldType.Kind() == reflect.Ptr
}

func IsSliceField(fieldType reflect.Type) bool {
	return fieldType.Kind() == reflect.Slice
}

func ToRPCClientStubInterfaceSDID(clientStubSDID string) (string, error) {
	splitedClientStubSDID := strings.Split(clientStubSDID, ".")
	if len(splitedClientStubSDID) < 2 {
		return "", fmt.Errorf("invalid client stub sdid %s", clientStubSDID)
	}
	splitedClientStubSDID[len(splitedClientStubSDID)-1] = ToFirstCharUpper(splitedClientStubSDID[len(splitedClientStubSDID)-1])
	return strings.Join(splitedClientStubSDID, "."), nil
}

func ToRPCServiceSDID(clientStubInterfaceSDID string) string {
	trimedSuffix := strings.TrimSuffix(clientStubInterfaceSDID, "IOCRPCClient")
	return strings.ReplaceAll(trimedSuffix, "/api.", ".")
}

func ToFirstCharLower(input string) string {
	return strings.ToLower(string(input[0])) + input[1:]
}

func ToFirstCharUpper(input string) string {
	return strings.ToUpper(string(input[0])) + input[1:]
}

func ToRPCClientStubSDID(clientStubInterfaceSDID string) (string, error) {
	splitedClientStubInterfaceSDID := strings.Split(clientStubInterfaceSDID, ".")
	if len(splitedClientStubInterfaceSDID) < 2 {
		return "", fmt.Errorf("invalid client stub interface sdid %s", clientStubInterfaceSDID)
	}
	splitedClientStubInterfaceSDID[len(splitedClientStubInterfaceSDID)-1] = ToFirstCharLower(splitedClientStubInterfaceSDID[len(splitedClientStubInterfaceSDID)-1])
	return strings.Join(splitedClientStubInterfaceSDID, "."), nil
}
