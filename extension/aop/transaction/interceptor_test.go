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

package transaction

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/alibaba/ioc-golang/aop/common"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/util"
)

type bizStruct struct {
}

const transactionMethodFullName = "github.com/alibaba/ioc-golang/extension/aop/transaction.(*bizStruct).TestMethodTransaction()"
const transactionMethodName = "TestMethodTransaction"
const mockGRID = int64(1)

func TestBeforeInvoke(t *testing.T) {
	// 1. register mock descriptor
	contextStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &context{}
		},
		ParamFactory: func() interface{} {
			return &contextParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(*contextParam)
			assert.Equal(t, transactionMethodFullName, param.entranceMethod)
			return param.init(i.(*context))
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(contextStructDescriptor)

	// 2. create to test object
	impl, err := GetinterceptorImplSingleton()
	assert.Nil(t, err)

	// 3. register a proxy biz struct
	bizStructSD := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &bizStruct{}
		},
		Metadata: map[string]interface{}{
			"aop": map[string]interface{}{
				"transaction": map[string]string{
					transactionMethodName:     "",
					"TestMethodDoTransaction": "TestMethodDoTransactionRollback",
				},
			},
		},
	}
	normal.RegisterStructDescriptor(bizStructSD)

	impl.BeforeInvoke(&aop.InvocationContext{
		MethodFullName: transactionMethodFullName,
		MethodName:     transactionMethodName,
		SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
		GrID:           mockGRID,
	})

	record, ok := impl.transactionGrIDMap.Load(mockGRID)
	assert.True(t, ok)
	assert.Equal(t, transactionMethodFullName, record.(contextIOCInterface).getEntranceMethod())
}

func TestAfterInvoke(t *testing.T) {
	t.Run("transaction success with normal returns", func(t *testing.T) {
		// 1. create to test object
		impl, err := GetinterceptorImplSingleton()
		assert.Nil(t, err)

		// 2. mock context
		ctx := newMockContextIOCInterface(t)
		currentMethodName := common.CurrentCallingMethodName(2)
		ctx.On("getEntranceMethod").Return(currentMethodName)
		ctx.On("finish").Once()

		impl.transactionGrIDMap.Store(mockGRID, ctx)

		aop.GetMockProxyFunctionLayer()(func() {
			impl.AfterInvoke(&aop.InvocationContext{
				MethodFullName: transactionMethodFullName,
				MethodName:     transactionMethodName,
				SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
				GrID:           mockGRID,
				ReturnValues: []reflect.Value{
					reflect.ValueOf("value"),
					reflect.ValueOf("stringValue"),
				},
			})
		})
	})

	t.Run("transaction success without return values", func(t *testing.T) {
		// 1. create to test object
		impl, err := GetinterceptorImplSingleton()
		assert.Nil(t, err)

		// 2. mock context
		ctx := newMockContextIOCInterface(t)
		currentMethodName := common.CurrentCallingMethodName(2)
		ctx.On("getEntranceMethod").Return(currentMethodName)
		ctx.On("finish").Once()

		impl.transactionGrIDMap.Store(mockGRID, ctx)

		aop.GetMockProxyFunctionLayer()(func() {
			impl.AfterInvoke(&aop.InvocationContext{
				MethodFullName: transactionMethodFullName,
				MethodName:     transactionMethodName,
				SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
				GrID:           mockGRID,
			})
		})
	})

	t.Run("transaction success with nil error returns", func(t *testing.T) {
		// 1. create to test object
		impl, err := GetinterceptorImplSingleton()
		assert.Nil(t, err)

		// 2. mock context
		ctx := newMockContextIOCInterface(t)
		currentMethodName := common.CurrentCallingMethodName(2)
		ctx.On("getEntranceMethod").Return(currentMethodName)
		ctx.On("finish").Once()

		impl.transactionGrIDMap.Store(mockGRID, ctx)

		aop.GetMockProxyFunctionLayer()(func() {
			impl.AfterInvoke(&aop.InvocationContext{
				MethodFullName: transactionMethodFullName,
				MethodName:     transactionMethodName,
				SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
				GrID:           mockGRID,
				ReturnValues: []reflect.Value{
					reflect.ValueOf("value"),
					reflect.ValueOf(new(error)),
				},
			})
		})
	})

	t.Run("transaction failed with entrance method", func(t *testing.T) {
		// 1. create to test object
		impl, err := GetinterceptorImplSingleton()
		assert.Nil(t, err)

		// 2. mock context
		expectErr := fmt.Errorf("error")

		ctx := newMockContextIOCInterface(t)
		currentMethodName := common.CurrentCallingMethodName(2)
		ctx.On("getEntranceMethod").Return(currentMethodName)
		ctx.On("failed", mock.MatchedBy(func(err error) bool {
			return err.Error() == expectErr.Error()
		})).Once()

		impl.transactionGrIDMap.Store(mockGRID, ctx)

		aop.GetMockProxyFunctionLayer()(func() {
			impl.AfterInvoke(&aop.InvocationContext{
				MethodFullName: transactionMethodFullName,
				MethodName:     transactionMethodName,
				SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
				GrID:           mockGRID,
				ReturnValues: []reflect.Value{
					reflect.ValueOf("value"),
					reflect.ValueOf(expectErr),
				},
			})
		})
	})

	t.Run("transaction success with not entrance method", func(t *testing.T) {
		// 1. create to test object
		impl, err := GetinterceptorImplSingleton()
		assert.Nil(t, err)

		// 2. mock context
		ctx := newMockContextIOCInterface(t)
		currentMethodName := common.CurrentCallingMethodName(2)
		ctx.On("getEntranceMethod").Return(currentMethodName + "notEntrance")
		ctx.On("addSuccessfullyCalledInvocationCtx", mock.MatchedBy(func(ctx *aop.InvocationContext) bool {
			return ctx.GrID == mockGRID
		})).Once()

		impl.transactionGrIDMap.Store(mockGRID, ctx)

		aop.GetMockProxyFunctionLayer()(func() {
			impl.AfterInvoke(&aop.InvocationContext{
				MethodFullName: transactionMethodFullName,
				MethodName:     transactionMethodName,
				SDID:           util.GetSDIDByStructPtr(&bizStruct{}),
				GrID:           mockGRID,
				ReturnValues: []reflect.Value{
					reflect.ValueOf("value"),
				},
			})
		})
	})
}
