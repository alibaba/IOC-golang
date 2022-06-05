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

package rpc_client

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/autowire"
)

type mockImpl struct {
}

const mockImplName = "github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client.mockImpl"

func TestAutowire_RegisterAndGetAllStructDescriptors(t *testing.T) {
	t.Run("test register and get all struct descriptors", func(t *testing.T) {
		sd := &autowire.StructDescriptor{
			Factory: func() interface{} {
				return &mockImpl{}
			},
		}
		RegisterStructDescriptor(sd)
		a := &Autowire{}
		allStructDesc := a.GetAllStructDescriptors()
		assert.NotNil(t, allStructDesc)
		structDesc, ok := allStructDesc[mockImplName]
		assert.True(t, ok)
		assert.Equal(t, mockImplName, structDesc.ID())
	})
}

func TestAutowire_TagKey(t *testing.T) {
	t.Run("test rpc autowire tag", func(t *testing.T) {
		a := &Autowire{}
		assert.Equal(t, Name, a.TagKey())
	})
}
