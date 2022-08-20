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

package config

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		configAutowire := &Autowire{}
		configAutowire.Autowire = normal.NewNormalAutowire(nil, &paramLoader{}, configAutowire)
		return configAutowire
	}())
}

const Name = "config"

type Autowire struct {
	autowire.Autowire
}

// TagKey re-write NormalAutowire
func (a *Autowire) TagKey() string {
	return Name
}

// GetAllStructDescriptors re-write NormalAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return configStructDescriptorMap
}

var configStructDescriptorMap = make(map[string]*autowire.StructDescriptor)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	configStructDescriptorMap[s.ID()] = s
	autowire.RegisterStructDescriptor(s)
}

func GetImpl(key string, configPrefix string) (interface{}, error) {
	return autowire.Impl(Name, key, configPrefix)
}
