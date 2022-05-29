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

package autowire

import (
	"fmt"
	"reflect"

	"github.com/alibaba/ioc-golang/autowire/util"
)

// Autowire
type Autowire interface {
	TagKey() string
	// IsSingleton means struct can be boot entrance, and only have one impl globally, only created once.
	IsSingleton() bool
	/*
		CanBeEntrance means the autowire sturct's param needs not to parse field tag value as param.
		By default, singleton can be boot entrance, and normal can't, because normal needs to try to parse tag value to find config key.
		But for grpc autowire, as singloton, can't be entrance because it needs to parse grpc type from cong tag.
	*/
	CanBeEntrance() bool
	Factory(sdID string) (interface{}, error)

	/*
		ParseSDID parse FieldInfo to struct descriptorId

		if field type is struct ptr like
		MyStruct *MyStruct `autowire-type:"MyStruct"`
		FieldInfo would be
		FieldInfo.FieldName == "MyStruct"
		FieldInfo.FieldType == "" // ATTENTION!!!
		FieldInfo.TagKey == "autowire-type"
		FieldInfo.TagValue == "MyStruct"
		You should make sure tag value contains ptr type if you use it.

		if field type is interface like
		MyStruct MyInterface ` `autowire-type:"MyStruct"`
		FieldInfo would be
		FieldInfo.FieldName == "MyStruct"
		FieldInfo.FieldType == "MyInterface"
		FieldInfo.TagKey == "autowire-type"
		FieldInfo.TagValue == "MyStruct"
	*/
	ParseSDID(field *FieldInfo) (string, error)
	ParseParam(sdID string, fi *FieldInfo) (interface{}, error)
	Construct(sdID string, impledPtr, param interface{}) (interface{}, error)
	GetAllStructDescriptors() map[string]*StructDescriptor
	InjectPosition() InjectPosition
}

var wrapperAutowireMap = make(map[string]WrapperAutowire)

func RegisterAutowire(autowire Autowire) {
	wrapperAutowireMap[autowire.TagKey()] = getWrappedAutowire(autowire, wrapperAutowireMap)
}

func GetAllWrapperAutowires() map[string]WrapperAutowire {
	return wrapperAutowireMap
}

// FieldInfo

type FieldInfo struct {
	FieldName        string
	FieldType        string
	FieldTypePkgPath string // PkgPath
	TagKey           string
	TagValue         string
}

// StructDescriptor

type StructDescriptor struct {
	Interface     interface{}
	Factory       func() interface{} // raw struct
	ParamFactory  func() interface{}
	ParamLoader   ParamLoader
	ConstructFunc func(impl interface{}, param interface{}) (interface{}, error) // injected
	DestroyFunc   func(impl interface{})
	Alias         string // alias of SDID

	impledStructPtr interface{} // impledStructPtr is only used to get name
	autowireType    string
}

func (ed *StructDescriptor) SetAutowireType(autowireType string) {
	ed.autowireType = autowireType
}

func (ed *StructDescriptor) AutowireType() string {
	return ed.autowireType
}

func (ed *StructDescriptor) ID() string {
	if ed.impledStructPtr == nil {
		ed.parse()
	}

	interfaceFullName := ed.populateInterfaceFullName()
	implFullName := ed.populateImplFullName()
	if interfaceFullName == implFullName {
		return interfaceFullName
	}

	return util.GetIdByNamePair(interfaceFullName, implFullName)
}

func (ed *StructDescriptor) populateInterfaceFullName() string {
	interfaceType := reflect.TypeOf(ed.Interface)
	interfacePkgPathName := interfaceType.Elem().PkgPath()
	interfaceName := interfaceType.Elem().Name()

	return fmt.Sprintf("%s.%s", interfacePkgPathName, interfaceName)
}

func (ed *StructDescriptor) populateImplFullName() string {
	implType := reflect.TypeOf(ed.impledStructPtr)
	implPkgPath := implType.Elem().PkgPath()
	implName := implType.Elem().Name()

	return fmt.Sprintf("%s.%s", implPkgPath, implName)
}

func (ed *StructDescriptor) parse() {
	ed.impledStructPtr = ed.Factory()
}

// ParamLoader is interface to load param
type ParamLoader interface {
	Load(sd *StructDescriptor, fi *FieldInfo) (interface{}, error)
}

// SDIDParser is interface to parse struct descriptor id
type SDIDParser interface {
	Parse(fi *FieldInfo) (string, error)
}

type InjectPosition int

const (
	AfterFactoryCalled     InjectPosition = 0
	AfterConstructorCalled InjectPosition = 1
)
