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
	"reflect"
	"sync"

	"github.com/alibaba/ioc-golang/aop/interceptor"
	"github.com/alibaba/ioc-golang/autowire"

	"github.com/petermattis/goid"

	"github.com/alibaba/ioc-golang/aop/interceptor/common"
)

type Interceptor struct {
	transactionGrIDMap sync.Map // transactionGrIDMap stores goroutine-id -> TxContext
}

func (t *Interceptor) BeforeInvoke(ctx *interceptor.InvocationContext) {
	// 1. if current invocation is already in transaction ?
	grID := goid.Get()
	if _, ok := t.transactionGrIDMap.Load(grID); ok {
		// this goRoutine is already in transaction
		return
	}
	// not in transaction

	// 2. if current method want to start a transaction ?
	sd := autowire.GetStructDescriptor(ctx.SDID)
	if sd == nil {
		// todo: print logs
		return
	}
	if _, ok := sd.TransactionMethodsMap[ctx.MethodName]; ok {
		// current method wants to start a transaction
		t.transactionGrIDMap.Store(grID, newContext(common.CurrentCallingMethodName()))
		return
	}
	// not in transaction, don't want to start a transaction
}

func (t *Interceptor) AfterInvoke(ctx *interceptor.InvocationContext) {
	// if current goRoutine is in the transaction ?
	grID := goid.Get()
	if val, ok := t.transactionGrIDMap.Load(grID); ok {
		// this goRoutine is in the transaction
		txCtx := val.(*context)

		// if invocation failed
		invocationFailed, err := isInvocationFailed(ctx.ReturnValues)

		// 1.1 if current invocation is the entrance of transaction ?
		// calculate level
		if common.TraceLevel(txCtx.entranceMethod) == 0 {
			// current invocation is the entrance of transaction
			t.transactionGrIDMap.Delete(grID)
			// if the transaction failed ?
			if invocationFailed {
				txCtx.failed(err)
				return
			}
			txCtx.finish()
			return
		}
		// current invocation is not the entrance of transaction
		// if the invocation is success ?
		if !invocationFailed {
			// the invocation is success, try to add to context
			txCtx.addSuccessfullyCalledInvocationCtx(ctx)
			return
		}
		// the invocation failed
		return
	}
	// the goRoutine is not in the transaction
}

func isInvocationFailed(returnValues []reflect.Value) (bool, error) {
	if len(returnValues) == 0 {
		return false, nil
	}
	finalReturnValue := returnValues[len(returnValues)-1]
	if err, ok := finalReturnValue.Interface().(error); ok && err != nil {
		return true, err
	}
	return false, nil
}

var transactionInterceptorSingleton *Interceptor

func GetTransactionInterceptor() *Interceptor {
	if transactionInterceptorSingleton == nil {
		transactionInterceptorSingleton = &Interceptor{}
	}
	return transactionInterceptorSingleton
}
