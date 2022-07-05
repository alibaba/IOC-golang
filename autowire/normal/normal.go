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

package normal

import (
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/base"
	"github.com/alibaba/ioc-golang/autowire/param_loader"
	"github.com/alibaba/ioc-golang/autowire/sdid_parser"
)

func init() {
	autowire.RegisterAutowire(NewNormalAutowire(nil, nil, nil))
}

const Name = "normal"

// NewNormalAutowire create a normal autowire based autowire, e.g. config, base.facade can be re-write to outer autowire
func NewNormalAutowire(sp autowire.SDIDParser, pl autowire.ParamLoader, facade autowire.Autowire) autowire.Autowire {
	if sp == nil {
		sp = sdid_parser.GetDefaultSDIDParser()
	}
	if pl == nil {
		pl = param_loader.GetDefaultParamLoader()
	}
	normalAutowire := &NormalAutowire{}
	if facade == nil {
		facade = normalAutowire
	}
	normalAutowire.AutowireBase = base.New(facade, sp, pl)
	return normalAutowire
}

type NormalAutowire struct {
	base.AutowireBase
}

func (n *NormalAutowire) IsSingleton() bool {
	return false
}

// GetAllStructDescriptors should be re-write by facade
func (n *NormalAutowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return normalEntryDescriptorMap
}

// TagKey should be re-writed by facade autowire
func (n *NormalAutowire) TagKey() string {
	return Name
}

func (s *NormalAutowire) CanBeEntrance() bool {
	return false
}

var normalEntryDescriptorMap = make(map[string]*autowire.StructDescriptor)

// developer APIs

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	sdID := s.ID()
	normalEntryDescriptorMap[sdID] = s
	autowire.RegisterStructDescriptor(sdID, s)
	if s.Alias != "" {
		autowire.RegisterAlias(s.Alias, sdID)
	}
}

func GetImpl(sdID string, param interface{}) (interface{}, error) {
	return autowire.Impl(Name, sdID, param)
}

func GetImplWithProxy(sdID string, param interface{}) (interface{}, error) {
	return autowire.ImplWithProxy(Name, sdID, param)
}
