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
	"testing"

	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/alibaba/IOC-Golang/autowire"
	"github.com/alibaba/IOC-Golang/autowire/mocks"
)

const (
	mockSDID          = "MockInterface-MockImpl"
	mockName          = "mockName"
	modifiedMockName  = "modifiedMockName"
	mockFieldType     = "Interface"
	mockFieldTagValue = "Impl"
)

type paramLoader struct {
}

func (p *paramLoader) Load(sd *autowire.StructDescriber, fi *autowire.FieldInfo) (interface{}, error) {
	return &param{
		name: modifiedMockName,
	}, nil
}

type param struct {
	name string
}

type MockImpl struct {
	name string
}

func (m *MockImpl) Name() string {
	return m.name
}

func newMockImpl() *MockImpl {
	return &MockImpl{name: mockName}
}

func TestAutowireBase_Construct(t *testing.T) {
	t.Run("test call construct success", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID: {
					ConstructFunc: func(impl interface{}, p interface{}) (interface{}, error) {
						mockImpl := impl.(*MockImpl)
						paramImpl := p.(*param)
						mockImpl.name = paramImpl.name
						return mockImpl, nil
					},
				},
			}
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		mockImpl := newMockImpl()
		param := &param{name: modifiedMockName}
		got, err := a.Construct(mockSDID, mockImpl, param)
		assert.Nil(t, err)
		impl, ok := got.(*MockImpl)
		assert.True(t, ok)
		assert.Equal(t, modifiedMockName, impl.Name())
	})

	t.Run("test call construct without constructFunc defined", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID: {},
			}
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		mockImpl := newMockImpl()
		param := &param{name: modifiedMockName}
		got, err := a.Construct(mockSDID, mockImpl, param)
		assert.Nil(t, err)
		impl, ok := got.(*MockImpl)
		assert.True(t, ok)
		assert.Equal(t, mockName, impl.Name())
	})

	t.Run("test call construct with sdID not found", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID + "foo": {
					Factory: func() interface{} {
						return newMockImpl()
					},
				},
			}
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		mockImpl := newMockImpl()
		param := &param{name: modifiedMockName}
		_, err := a.Construct(mockSDID, mockImpl, param)
		assert.Equal(t, perrors.Errorf("struct ID %s struct describer not found ", mockSDID).Error(), err.Error())
	})

	t.Run("test call construct with struct describer map empty", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return nil
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		mockImpl := newMockImpl()
		param := &param{name: modifiedMockName}
		_, err := a.Construct(mockSDID, mockImpl, param)
		assert.Equal(t, "struct describer map is empty.", err.Error())
	})
}

func TestAutowireBase_Factory(t *testing.T) {
	t.Run("test call factory success", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID: {
					Factory: func() interface{} {
						return newMockImpl()
					},
				},
			}
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		got, err := a.Factory(mockSDID)
		assert.Nil(t, err)
		impl, ok := got.(*MockImpl)
		assert.True(t, ok)
		assert.Equal(t, mockName, impl.Name())
	})

	t.Run("test call factory with struct describer map empty ", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return nil
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		_, err := a.Factory(mockSDID)
		assert.Equal(t, "struct describer map is empty.", err.Error())
	})

	t.Run("test call factory with struct describer ID not found ", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID + "foo": {
					Factory: func() interface{} {
						return newMockImpl()
					},
				},
			}
		})
		a := &AutowireBase{
			facadeAutowire: facadeAutowire,
		}
		_, err := a.Factory(mockSDID)
		assert.Equal(t, perrors.Errorf("struct ID %s struct describer not found ", mockSDID).Error(), err.Error())
	})
}

func TestAutowireBase_InjectPosition(t *testing.T) {
	a := &AutowireBase{}
	assert.Equal(t, autowire.AfterFactoryCalled, a.InjectPosition())
}

func TestAutowireBase_ParseParam(t *testing.T) {
	defaultParamLoader := mocks.NewParamLoader(t)
	defaultParamLoader.On("Load", mock.MatchedBy(func(sd *autowire.StructDescriber) bool {
		return true
	}), mock.MatchedBy(func(fi *autowire.FieldInfo) bool {
		return true
	})).Return(func(sd *autowire.StructDescriber, fi *autowire.FieldInfo) interface{} {
		return &param{
			name: mockName,
		}
	}, func(sd *autowire.StructDescriber, fi *autowire.FieldInfo) error {
		return nil
	})

	t.Run("test parse param success with param loader defined in sd", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID: {
					ParamFactory: func() interface{} {
						return &param{
							mockName,
						}
					},
					ParamLoader: &paramLoader{},
				},
			}
		})

		a := &AutowireBase{
			paramLoader:    defaultParamLoader,
			facadeAutowire: facadeAutowire,
		}
		got, err := a.ParseParam(mockSDID, nil)
		assert.Nil(t, err)
		paramImpl, ok := got.(*param)
		assert.True(t, ok)
		assert.Equal(t, modifiedMockName, paramImpl.name)
	})

	t.Run("test parse param success with default param loader", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{
				mockSDID: {
					ParamFactory: func() interface{} {
						return &param{
							mockName,
						}
					},
				},
			}
		})

		a := &AutowireBase{
			paramLoader:    defaultParamLoader,
			facadeAutowire: facadeAutowire,
		}
		got, err := a.ParseParam(mockSDID, nil)
		assert.Nil(t, err)
		paramImpl, ok := got.(*param)
		assert.True(t, ok)
		assert.Equal(t, mockName, paramImpl.name)
	})

	t.Run("test parse param success with empty sd map", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return nil
		})

		a := &AutowireBase{
			paramLoader:    defaultParamLoader,
			facadeAutowire: facadeAutowire,
		}
		_, err := a.ParseParam(mockSDID, nil)
		assert.Equal(t, "struct describer map is empty.", err.Error())
	})

	t.Run("test parse param success with sdID not found", func(t *testing.T) {
		facadeAutowire := mocks.NewFacadeAutowire(t)
		facadeAutowire.On("GetAllStructDescribers").Return(func() map[string]*autowire.StructDescriber {
			return map[string]*autowire.StructDescriber{}
		})

		a := &AutowireBase{
			paramLoader:    defaultParamLoader,
			facadeAutowire: facadeAutowire,
		}
		_, err := a.ParseParam(mockSDID, nil)
		assert.Equal(t, perrors.Errorf("struct ID %s struct describer not found ", mockSDID).Error(), err.Error())
	})
}

func TestAutowireBase_ParseSDID(t *testing.T) {
	t.Run("test parse sdid", func(t *testing.T) {
		a := &AutowireBase{
			sdIDParser: func() autowire.SDIDParser {
				mockSDIDParser := mocks.NewSDIDParser(t)
				mockSDIDParser.On("Parse", mock.MatchedBy(func(fi *autowire.FieldInfo) bool {
					return true
				})).Return(func(fi *autowire.FieldInfo) string {
					return fi.FieldType + "-" + fi.TagValue
				}, func(fi *autowire.FieldInfo) error {
					return nil
				})
				return mockSDIDParser
			}(),
		}
		field := &autowire.FieldInfo{
			FieldType: mockFieldType,
			TagValue:  mockFieldTagValue,
		}
		got, err := a.ParseSDID(field)
		assert.Nil(t, err)
		assert.Equal(t, mockFieldType+"-"+mockFieldTagValue, got)
	})
}

func TestNew(t *testing.T) {
	assert.NotNil(t, New(autowire.NewMockAutowire(t), mocks.NewSDIDParser(t), mocks.NewParamLoader(t)))
}
