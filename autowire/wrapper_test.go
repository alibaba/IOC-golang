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
	"testing"

	"github.com/alibaba/ioc-golang/autowire/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubImpl struct {
}

func (m *MockSubImpl) GetName() string {
	return "MockSubImpl"
}

type MockSubInterface interface {
	GetName() string
}

type MockImpl struct {
	SubImpl MockSubInterface `singleton:"github.com/alibaba/ioc-golang/autowire.MockSubImpl"`
}

var mockSDID = util.GetSDIDByStructPtr(&MockImpl{})
var mockSubSDID = util.GetSDIDByStructPtr(&MockSubImpl{})

func TestWrapperAutowireImpl_ImplWithParam(t *testing.T) {
	t.Run("test impl with param", func(t *testing.T) {
		mockAutowire := NewMockAutowire(t)
		mockAutowire.On("GetAllStructDescriptors").Return(func() map[string]*StructDescriptor {
			return map[string]*StructDescriptor{
				mockSubSDID: {
					Factory: func() interface{} {
						return &MockSubImpl{}
					},
				},
				mockSDID: {
					Factory: func() interface{} {
						return &MockImpl{}
					},
				},
			}
		})
		RegisterStructDescriptor(mockSubSDID, &StructDescriptor{
			Factory: func() interface{} {
				return &MockSubImpl{}
			},
		})
		RegisterStructDescriptor(mockSDID, &StructDescriptor{
			Factory: func() interface{} {
				return &MockImpl{}
			},
		})
		mockAutowire.On("ParseSDID", mock.MatchedBy(func(fi *FieldInfo) bool {
			return true
		})).Return(func(fi *FieldInfo) string {
			return fi.TagValue
		}, func(fi *FieldInfo) error {
			return nil
		})
		mockAutowire.On("ParseParam", mock.MatchedBy(func(sdID string) bool {
			return true
		}), mock.MatchedBy(func(fi *FieldInfo) bool {
			return true
		})).Return(func(sdID string, fi *FieldInfo) interface{} {
			return nil
		}, func(sdID string, fi *FieldInfo) error {
			return nil
		})

		mockAutowire.On("TagKey").Return(func() string {
			return "singleton"
		})
		mockAutowire.On("Factory", mock.MatchedBy(func(sdID string) bool {
			return true
		})).Return(func(sdID string) interface{} {
			if sdID == mockSDID {
				return &MockImpl{}
			} else {
				return &MockSubImpl{}
			}
		}, func(sdID string) error {
			return nil
		})

		mockAutowire.On("Construct", mock.MatchedBy(func(sdID string) bool {
			return true
		}), mock.MatchedBy(func(impl interface{}) bool {
			return true
		}), mock.MatchedBy(func(param interface{}) bool {
			return true
		})).Return(func(sdID string, impl, param interface{}) interface{} {
			return impl
		}, func(sdID string, impl, param interface{}) error {
			return nil
		})

		mockAutowire.On("IsSingleton").Return(func() bool {
			return true
		})
		mockAutowire.On("InjectPosition").Return(func() InjectPosition {
			return AfterFactoryCalled
		})
		wa := getWrappedAutowire(mockAutowire, GetAllWrapperAutowires())
		RegisterAutowire(wa)
		impl, err := wa.ImplWithParam(mockSDID, nil, false)
		assert.Nil(t, err)
		mockImpl, ok := impl.(*MockImpl)
		assert.True(t, ok)
		assert.NotNil(t, mockImpl.SubImpl)
		assert.Equal(t, "MockSubImpl", mockImpl.SubImpl.GetName())
	})
}

func TestGetWrappedAutowire(t *testing.T) {
	assert.NotNil(t, getWrappedAutowire(nil, nil))
}
