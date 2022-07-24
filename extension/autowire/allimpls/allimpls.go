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

package allimpls

import (
	"fmt"
	"reflect"

	"github.com/alibaba/ioc-golang/autowire/normal"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/autowire/util"
	"github.com/alibaba/ioc-golang/logger"
)

func init() {
	autowire.RegisterAutowire(func() autowire.Autowire {
		allImplsAutowire := &Autowire{sdIDParser: getSDIDParserSingleton()}
		allImplsAutowire.Autowire = singleton.NewSingletonAutowire(getSDIDParserSingleton(), nil, allImplsAutowire)
		return allImplsAutowire
	}())
}

const Name = "allimpls"

type Autowire struct {
	autowire.Autowire
	sdIDParser *sdIDParser
}

// TagKey re-write SingletonAutowire
func (a *Autowire) TagKey() string {
	return Name
}

func (a *Autowire) IsSingleton() bool {
	return false
}

// GetAllStructDescriptors re-write SingletonAutowire
func (a *Autowire) GetAllStructDescriptors() map[string]*autowire.StructDescriptor {
	return allImplsStructDescriptorMap
}

func (a *Autowire) Factory(fieldInterfaceID string) (interface{}, error) {
	return a.sdIDParser.newFieldInterfaceSliceValue(fieldInterfaceID), nil
}

func (a *Autowire) Construct(fieldInterfaceID string, sliceValue, _ interface{}) (interface{}, error) {
	implSDs, ok := intefaceSDIDImplsSDMap[fieldInterfaceID]
	if !ok {
		logger.Red("[Autowire allimpls] interface ID %s's implementations not found", fieldInterfaceID)
		return nil, fmt.Errorf("[Autowire allimpls] interface ID %s implementations not found", fieldInterfaceID)
	}
	for _, implSD := range implSDs {
		var impl interface{}
		var err error
		autowireType := singleton.Name
		if autowireTypFromMetadata := parseAllImplsItemAutowireTypeFromSDMetadata(implSD.Metadata); autowireTypFromMetadata != "" {
			autowireType = autowireTypFromMetadata
		}
		if implSD.DisableProxy {
			impl, err = autowire.Impl(autowireType, implSD.ID(), nil)
		} else {
			impl, err = autowire.ImplWithProxy(autowireType, implSD.ID(), nil)
		}

		if err != nil {
			return impl, err
		}
		sliceValue = reflect.Append(sliceValue.(reflect.Value), reflect.ValueOf(impl))
	}
	return sliceValue.(reflect.Value).Interface(), nil
}

var allImplsStructDescriptorMap = make(map[string]*autowire.StructDescriptor)
var intefaceSDIDImplsSDMap = make(map[string][]*autowire.StructDescriptor)
var impledInterfaceSDIDTypeMap = make(map[string]reflect.Type)

func RegisterStructDescriptor(s *autowire.StructDescriptor) {
	sdID := s.ID()
	allImplsStructDescriptorMap[sdID] = s

	// check and register sd to item impl autowire type
	if err := checkAndRegisterSDToAutowireType(s); err != nil {
		logger.Red(err.Error())
		return
	}

	// register sd to autowire base layer
	autowire.RegisterStructDescriptor(sdID, s)
	if s.Alias != "" {
		// register alias if necessary
		autowire.RegisterAlias(s.Alias, sdID)
	}
	allImpledIntefaces := parseAllImpledIntefacesFromSDMetadata(s.Metadata)
	for _, impledInteface := range allImpledIntefaces {
		interfaceSDID := util.GetSDIDByStructPtr(impledInteface)
		interfaceType := util.GetTypeFromInterface(impledInteface)
		// 1. record current impl sd to interface -> sd mapping
		existingImpls, ok := intefaceSDIDImplsSDMap[interfaceSDID]
		if !ok {
			existingImpls = make([]*autowire.StructDescriptor, 0)
		}
		intefaceSDIDImplsSDMap[interfaceSDID] = append(existingImpls, s)

		// 2. record interface type for reflection
		impledInterfaceSDIDTypeMap[interfaceSDID] = interfaceType

		// 3. create interface empty struct descriptor
		allImplsStructDescriptorMap[interfaceSDID] = &autowire.StructDescriptor{}
	}
}

func GetImpl(key string) (interface{}, error) {
	return autowire.Impl(Name, key, nil)
}

func checkAndRegisterSDToAutowireType(s *autowire.StructDescriptor) error {
	autowireTypFromMetadata := parseAllImplsItemAutowireTypeFromSDMetadata(s.Metadata)
	if autowireTypFromMetadata == "" {
		// default singleton autowire type
		singleton.RegisterStructDescriptor(s)
		return nil
	}

	switch autowireTypFromMetadata {
	case singleton.Name:
		singleton.RegisterStructDescriptor(s)
	case normal.Name:
		normal.RegisterStructDescriptor(s)
	default:
		return fmt.Errorf("[Autowire allimpls] Found invalid item impl autowire type %s, now only singleton(default) and normal are supported", autowireTypFromMetadata)
	}
	return nil
}
