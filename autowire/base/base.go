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

package base

import (
	"github.com/fatih/color"

	perrors "github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/autowire"
)

type FacadeAutowire interface {
	GetAllStructDescriptors() map[string]*autowire.StructDescriptor
	TagKey() string
}

// New return new AutowireBase
func New(facadeAutowire FacadeAutowire, sp autowire.SDIDParser, pl autowire.ParamLoader) AutowireBase {
	return AutowireBase{
		facadeAutowire: facadeAutowire,
		sdIDParser:     sp,
		paramLoader:    pl,
	}
}

type AutowireBase struct {
	facadeAutowire FacadeAutowire
	sdIDParser     autowire.SDIDParser
	paramLoader    autowire.ParamLoader
}

func (a *AutowireBase) Factory(sdID string) (interface{}, error) {
	allStructDescriptor := a.facadeAutowire.GetAllStructDescriptors()
	if allStructDescriptor == nil {
		return nil, perrors.New("struct descriptor map is empty.")
	}
	sd, ok := allStructDescriptor[sdID]
	if !ok {
		return nil, perrors.Errorf("struct ID %s struct descriptor not found ", sdID)
	}
	return sd.Factory(), nil
}

func (a *AutowireBase) Construct(sdID string, impledPtr, param interface{}) (interface{}, error) {
	allStructDescriptor := a.facadeAutowire.GetAllStructDescriptors()
	if allStructDescriptor == nil {
		return nil, perrors.New("struct descriptor map is empty.")
	}
	sd, ok := allStructDescriptor[sdID]
	if !ok {
		return nil, perrors.Errorf("struct ID %s struct descriptor not found ", sdID)
	}
	if sd.ConstructFunc != nil {
		return sd.ConstructFunc(impledPtr, param)
	}
	return impledPtr, nil
}

func (a *AutowireBase) ParseSDID(field *autowire.FieldInfo) (string, error) {
	return a.sdIDParser.Parse(field)
}

func (a *AutowireBase) ParseParam(sdID string, fi *autowire.FieldInfo) (interface{}, error) {
	allStructDescriptor := a.facadeAutowire.GetAllStructDescriptors()
	if allStructDescriptor == nil {
		return nil, perrors.New("struct descriptor map is empty.")
	}
	sd, ok := allStructDescriptor[sdID]
	if !ok {
		return nil, perrors.Errorf("struct ID %s struct descriptor not found ", sdID)
	}
	if sd.ParamFactory == nil {
		// doesn't register param factory, do not load param, return with success
		return nil, nil
	}
	if sd.ParamLoader != nil {
		// try to use sd ParamLoader
		param, err := sd.ParamLoader.Load(sd, fi)
		if err == nil {
			return param, nil
		} else {
			// log warning, given pl load failed, fall back to default
			color.Red("[Autoware Base] Load SD %s param with defined sd.ParamLoader error: %s\n"+
				"Try load by autowire %s's default paramloader", sd.ID(), err, a.facadeAutowire.TagKey())
		}
	}
	// use autowire defined paramLoader as fall back
	return a.paramLoader.Load(sd, fi)
}

func (a *AutowireBase) InjectPosition() autowire.InjectPosition {
	return autowire.AfterFactoryCalled
}
