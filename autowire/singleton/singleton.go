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

package singleton

import (
	"github.com/alibaba/IOC-Golang/autowire"
	"github.com/alibaba/IOC-Golang/autowire/base"
	"github.com/alibaba/IOC-Golang/autowire/param_loader"
	"github.com/alibaba/IOC-Golang/autowire/sdid_parser"
)

func init() {
	autowire.RegisterAutowire(NewSingletonAutowire(nil, nil, nil))
}

const Name = "singleton"

var singletonStructDescribersMap = make(map[string]*autowire.StructDescriber)

// autowire APIs

// NewSingletonAutowire create a singleton autowire based autowire, e.g. grpc, base.facade can be re-write to outer autowire
func NewSingletonAutowire(sp autowire.SDIDParser, pl autowire.ParamLoader, facade autowire.Autowire) autowire.Autowire {
	if sp == nil {
		sp = sdid_parser.GetDefaultSDIDParser()
	}
	if pl == nil {
		pl = param_loader.GetDefaultParamLoader()
	}
	singletonAutowire := &SingletonAutowire{
		paramLoader: pl,
		sdIDParser:  sp,
	}
	if facade == nil {
		facade = singletonAutowire
	}
	singletonAutowire.AutowireBase = base.New(facade, sp, pl)
	return singletonAutowire

}

type SingletonAutowire struct {
	base.AutowireBase
	paramLoader autowire.ParamLoader
	sdIDParser  autowire.SDIDParser
}

// GetAllStructDescribers should be re-write by facade
func (s *SingletonAutowire) GetAllStructDescribers() map[string]*autowire.StructDescriber {
	return singletonStructDescribersMap
}

// TagKey should be re-writed by facade autowire
func (s *SingletonAutowire) TagKey() string {
	return Name
}

func (s *SingletonAutowire) IsSingleton() bool {
	return true
}

func (s *SingletonAutowire) CanBeEntrance() bool {
	return true
}

// developer APIs

func RegisterStructDescriber(sd *autowire.StructDescriber) {
	sd.SetAutowireType(Name)
	singletonStructDescribersMap[sd.ID()] = sd
}

func GetImpl(sdId string) (interface{}, error) {
	return autowire.Impl(Name, sdId, nil)
}
