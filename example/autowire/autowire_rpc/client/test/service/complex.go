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

package service

import (
	"fmt"

	"github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/client/test/dto"
)

// +ioc:autowire=true
// +ioc:autowire:type=rpc

type ComplexService struct {
}

func (s *ComplexService) RPCBasicType(name string, age int, age32 int32, age64 int64, ageF32 float32, ageF64 float64, namePtr *string, agePtr *int, age32Ptr *int32, age64Ptr *int64, ageF32Ptr *float32, ageF64Ptr *float64) (string, int, int32, int64, float32, float64, *string, *int, *int32, *int64, *float32, *float64) {
	return name, age, age32, age64, ageF32, ageF64, namePtr, agePtr, age32Ptr, age64Ptr, ageF32Ptr, ageF64Ptr
}

func (s *ComplexService) RPCWithoutParamAndReturnValue() {
}

func (s *ComplexService) RPCWithoutParam() (*dto.User, error) {
	return &dto.User{
		Id:   1,
		Name: "laurence",
		Age:  23,
	}, nil
}

func (s *ComplexService) RPCWithoutReturnValue(user *dto.User) {
}

func (s *ComplexService) RPCWithCustomValue(customStruct dto.CustomStruct, customStruct2 *dto.CustomStruct) (dto.CustomStruct, *dto.CustomStruct) {
	return customStruct, customStruct2
}

func (s *ComplexService) RPCWithError() (*dto.User, error) {
	errMsg := "custom"
	return &dto.User{
		Id:   1,
		Name: "laurence",
		Age:  23,
	}, fmt.Errorf("custom error = %s", errMsg)
}

func (s *ComplexService) RPCWithParamCustomMethod(customStruct dto.CustomStruct) dto.User {
	return customStruct.GetUser()
}
