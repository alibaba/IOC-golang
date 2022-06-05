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
	"github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/dto"
)

// +ioc:autowire=true
// +ioc:autowire:type=rpc

type ServiceStruct struct {
}

func (s *ServiceStruct) GetUser(name string, age int) (*dto.User, error) {
	return &dto.User{
		Id:   1,
		Name: name,
		Age:  age,
	}, nil
}
