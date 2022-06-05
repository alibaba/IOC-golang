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
	"os"
	"reflect"
	"runtime"

	"github.com/fatih/color"

	perrors "github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/autowire/util"
)

type WrapperAutowire interface {
	Autowire

	ImplWithoutParam(sdID string) (interface{}, error)
	ImplWithParam(sdID string, param interface{}) (interface{}, error)

	implWithField(info *FieldInfo) (interface{}, error)
}

func getWrappedAutowire(autowire Autowire, allAutowires map[string]WrapperAutowire) WrapperAutowire {
	return &WrapperAutowireImpl{
		Autowire:           autowire,
		allAutowires:       allAutowires,
		singletonImpledMap: map[string]interface{}{},
	}
}

type WrapperAutowireImpl struct {
	Autowire
	singletonImpledMap map[string]interface{}
	allAutowires       map[string]WrapperAutowire
}

// ImplWithParam is used to get impled struct with param
func (w *WrapperAutowireImpl) ImplWithParam(sdID string, param interface{}) (interface{}, error) {
	// 1. check singleton
	if singletonImpledPtr, ok := w.singletonImpledMap[sdID]; w.Autowire.IsSingleton() && ok {
		return singletonImpledPtr, nil
	}

	// 2. factory
	impledPtr, err := w.Autowire.Factory(sdID)
	if err != nil {
		return nil, err
	}

	if w.Autowire.InjectPosition() == AfterFactoryCalled {
		if err := w.inject(impledPtr, sdID); err != nil {
			return nil, err
		}
	}

	// 4. construct field
	impledPtr, err = w.Autowire.Construct(sdID, impledPtr, param)
	if err != nil {
		return nil, err
	}

	if w.Autowire.InjectPosition() == AfterConstructorCalled {
		if err := w.inject(impledPtr, sdID); err != nil {
			return nil, err
		}
	}

	// 5. record singleton ptr
	if w.Autowire.IsSingleton() {
		w.singletonImpledMap[sdID] = impledPtr
	}
	return impledPtr, nil
}

// ImplWithoutParam is used to create param from field without param
func (w *WrapperAutowireImpl) ImplWithoutParam(sdID string) (interface{}, error) {
	param, err := w.ParseParam(sdID, nil)
	if err != nil {
		if w.Autowire.IsSingleton() {
			// FIXME: ignore parse param error, because of singleton with empty param also try to find property from config file
			color.Red("[Wrapper Autowire] Parse param from config file with sdid %s failed, error: %s, continue with nil param.", sdID, err)
			return w.ImplWithParam(sdID, param)
		} else {
			return nil, err
		}
	}
	return w.ImplWithParam(sdID, param)
}

// ImplWithField is used to create param from field and call ImplWithParam
func (w *WrapperAutowireImpl) implWithField(fi *FieldInfo) (interface{}, error) {
	sdID, err := w.ParseSDID(fi)
	if err != nil {
		return nil, err
	}
	param, err := w.ParseParam(sdID, fi)
	if err != nil {
		if w.Autowire.IsSingleton() {
			// FIXME: ignore parse param error, because of singleton with empty param also try to find property from config file
			color.Red("[Wrapper Autowire] Parse param from config file with sdid %s failed, error: %s, continue with nil param.", sdID, err)
			return w.ImplWithParam(sdID, param)
		} else {
			return nil, err
		}
	}
	return w.ImplWithParam(sdID, param)
}

// inject do tag autowire and monkey inject
func (w *WrapperAutowireImpl) inject(impledPtr interface{}, sdId string) error {
	sd := w.Autowire.GetAllStructDescriptors()[sdId]

	// 1. reflect
	valueOf := reflect.ValueOf(impledPtr)
	valueOfElem := valueOf.Elem()
	typeOf := valueOfElem.Type()
	if typeOf.Kind() != reflect.Struct {
		// not struct, no needs to inject tag and monkey, just return
		return nil
	}

	// deal with struct
	// 2. tag inject
	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		field := typeOf.Field(i)
		var subImpledPtr interface{}
		tagKey := ""
		tagValue := ""
		for _, aw := range w.allAutowires {
			if val, ok := field.Tag.Lookup(aw.TagKey()); ok {
				fieldType := buildFiledTypeFullName(field.Type)
				fieldInfo := &FieldInfo{
					FieldName: field.Name,
					FieldType: fieldType,
					TagKey:    aw.TagKey(),
					TagValue:  val,
				}
				// create param from field info
				var err error
				subImpledPtr, err = aw.implWithField(fieldInfo)
				if err != nil {
					return err
				}
				tagKey = aw.TagKey()
				tagValue = val
				break // only one tag is support
			}
		}
		if tagKey == "" && tagValue == "" {
			continue
		}
		// set field
		subService := valueOfElem.Field(i)
		if !(subService.IsValid() && subService.CanSet()) {
			err := perrors.Errorf("Failed to autowire struct %s's impl %s service. It's field %s with tag '%s:\"%s\"', please check if the field is exported",
				sd.ID(), util.GetStructName(impledPtr), field.Type.Name(), tagKey, tagValue)
			return err
		}
		subService.Set(reflect.ValueOf(subImpledPtr))
	}
	// 3. monkey
	if monkeyFunction := GetMonkeyFunction(); (os.Getenv("GOARCH") == "amd64" || runtime.GOARCH == "amd64") && monkeyFunction != nil {
		// only amd64-os/amd64-go-arch-env with build flags '-gcflags="-N -l" -tags iocdebug' can inject monkey function
		monkeyFunction(impledPtr, sd.ID())
	}
	return nil
}

func buildFiledTypeFullName(fieldType reflect.Type) string {
	if fieldType.Kind() == reflect.Ptr {
		return fieldType.Elem().PkgPath() + "." + fieldType.Elem().Name()
	}
	return fieldType.PkgPath() + "." + fieldType.Name()
}
