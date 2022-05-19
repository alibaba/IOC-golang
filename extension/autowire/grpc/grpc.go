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

package grpc

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		grpcAutowire := &Autowire{}
		grpcAutowire.Autowire = singleton.NewSingletonAutowire(&sdIDParser{}, &paramLoader{}, grpcAutowire)
		return grpcAutowire
	}())
}

const Name = "grpc"

type Autowire struct {
	autowire.Autowire
}

// TagKey re-write SingletonAutowire
func (a *Autowire) TagKey() string {
	return Name
}

func (a *Autowire) CanBeEntrance() bool {
	return false
}

func (a *Autowire) InjectPosition() autowire.InjectPosition {
	return autowire.AfterConstructorCalled
}

// GetAllStructDescriptors re-write SingletonAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return grpcStructDescriptorMap
}

var grpcStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	s.SetAutowireType(Name)
	grpcStructDescriptorMap[s.ID()] = s
}

func GetImpl(extensionId string) (interface{}, error) {
	return autowire.Impl(Name, extensionId, nil)
}
