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

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/debug/interceptor"
)

type Context struct {
	successfullyInvokedRollbackAbleMethodsSnapshot []InvokedMethodSnapshot
	entranceMethod                                 string
}

func (c *Context) Finish() {

}

func (c *Context) Failed(err error) {
	for i := len(c.successfullyInvokedRollbackAbleMethodsSnapshot) - 1; i >= 0; i-- {
		snapshot := c.successfullyInvokedRollbackAbleMethodsSnapshot[i]
		snapshot.Rollback(err)
	}
}

func (c *Context) AddSuccessfullyCalledInvocationCtx(ctx *interceptor.InvocationContext) {
	sd := autowire.GetStructDescriptor(ctx.SDID)
	if sd == nil {
		// todo: print logs
		return
	}
	if rollbackMethodName, ok := sd.TransactionMethodsMap[ctx.MethodName]; ok && rollbackMethodName != "" {
		c.successfullyInvokedRollbackAbleMethodsSnapshot = append(c.successfullyInvokedRollbackAbleMethodsSnapshot, InvokedMethodSnapshot{
			invocationCtx:      ctx,
			rollbackMethodName: rollbackMethodName,
		})
	}
}

type InvokedMethodSnapshot struct {
	invocationCtx      *interceptor.InvocationContext
	rollbackMethodName string
}

func (snapshot *InvokedMethodSnapshot) Rollback(err error) {
	valueOf := reflect.ValueOf(snapshot.invocationCtx.ProxyServicePtr)
	valueOfElem := valueOf.Elem()
	funcRaw := valueOfElem.FieldByName(snapshot.rollbackMethodName + "_")
	rollbackParam := snapshot.invocationCtx.Params
	rollbackParam = append(rollbackParam, reflect.ValueOf(err))
	funcRaw.Call(rollbackParam)
}
