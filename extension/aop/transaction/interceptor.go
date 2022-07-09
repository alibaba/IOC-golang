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
	"sync"

	"github.com/alibaba/ioc-golang/aop"
	"github.com/alibaba/ioc-golang/autowire"

	"github.com/alibaba/ioc-golang/aop/common"
)

type interceptorImpl struct {
	transactionGrIDMap sync.Map // transactionGrIDMap stores goroutine-id -> TxContext
}

func (t *interceptorImpl) BeforeInvoke(ctx *aop.InvocationContext) {
	// 1. if current invocation is already in transaction ?
	if _, ok := t.transactionGrIDMap.Load(ctx.GrID); ok {
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
		t.transactionGrIDMap.Store(ctx.GrID, newContext(common.CurrentCallingMethodName()))
		return
	}
	// not in transaction, don't want to start a transaction
}

func (t *interceptorImpl) AfterInvoke(ctx *aop.InvocationContext) {
	// if current goRoutine is in the transaction ?
	if val, ok := t.transactionGrIDMap.Load(ctx.GrID); ok {
		// this goRoutine is in the transaction
		txCtx := val.(*context)

		// if invocation failed
		invocationFailed, err := common.IsInvocationFailed(ctx.ReturnValues)

		// 1.1 if current invocation is the entrance of transaction ?
		// calculate level
		if common.TraceLevel(txCtx.entranceMethod) == 0 {
			// current invocation is the entrance of transaction
			t.transactionGrIDMap.Delete(ctx.GrID)
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

var transactionInterceptorSingleton *interceptorImpl

func getTransactionInterceptorSingleton() *interceptorImpl {
	if transactionInterceptorSingleton == nil {
		transactionInterceptorSingleton = &interceptorImpl{}
	}
	return transactionInterceptorSingleton
}
