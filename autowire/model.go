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
	"reflect"

	"github.com/alibaba/ioc-golang/autowire/util"
)

// Autowire
type Autowire interface {
	TagKey() string
	// IsSingleton means struct can be boot entrance, and only have one impl globally, only created once.
	IsSingleton() bool
	/*
		CanBeEntrance means the struct is loaded at the start of application.
		By default, only rpc-server can be entrance.
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
	FieldName         string
	FieldType         string
	TagKey            string
	TagValue          string
	FieldReflectType  reflect.Type
	FieldReflectValue reflect.Value
}

// StructDescriptor

type StructDescriptor struct {
	Factory       func() interface{} // raw struct
	ParamFactory  func() interface{}
	ParamLoader   ParamLoader
	ConstructFunc func(impl interface{}, param interface{}) (interface{}, error) // injected
	DestroyFunc   func(impl interface{})
	Alias         string // alias of SDID
	DisableProxy  bool   // disable proxy and aop
	Metadata      Metadata

	impledStructPtr interface{} // impledStructPtr is only used to get name
}

func (ed *StructDescriptor) ID() string {
	return util.GetSDIDByStructPtr(ed.getStructPtr())
}

func (ed *StructDescriptor) getStructPtr() interface{} {
	if ed.impledStructPtr == nil {
		ed.impledStructPtr = ed.Factory()
	}
	return ed.impledStructPtr
}

// ParamLoader is interface to load param
type ParamLoader interface {
	Load(sd *StructDescriptor, fi *FieldInfo) (interface{}, error)
}

// SDIDParser is interface to parse struct descriptor id
type SDIDParser interface {
	Parse(fi *FieldInfo) (string, error)
}

// Metadata is SD metadata
type Metadata map[string]interface{}

type InjectPosition int

const (
	AfterFactoryCalled     InjectPosition = 0
	AfterConstructorCalled InjectPosition = 1
)
