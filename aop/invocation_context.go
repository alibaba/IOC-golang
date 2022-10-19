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

package aop

import (
	"reflect"

	"github.com/google/uuid"

	"github.com/petermattis/goid"
)

type InvocationContext struct {
	ID              uuid.UUID
	ProxyServicePtr interface{}
	SDID            string
	MethodName      string
	MethodFullName  string
	Params          []reflect.Value
	ReturnValues    []reflect.Value
	GrID            int64
	Metadata        map[string]interface{}
}

func (c *InvocationContext) SetReturnValues(returnValues []reflect.Value) {
	c.ReturnValues = returnValues
}

func NewInvocationContext(proxyServicePtr interface{}, sdid, methodName, methodFullName string, params []reflect.Value) *InvocationContext {
	grID := goid.Get()
	newInvocationCtx := &InvocationContext{
		ID:              uuid.New(),
		ProxyServicePtr: proxyServicePtr,
		SDID:            sdid,
		Metadata:        make(map[string]interface{}),
		MethodName:      methodName,
		Params:          params,
		GrID:            grID,
		MethodFullName:  methodFullName,
	}
	invocationContextMap[grID] = newInvocationCtx
	return newInvocationCtx
}

// invocationContextMap is thread safe, because every gr can only read or write their own key
var invocationContextMap = make(map[int64]*InvocationContext)

func GetCurrentInvocationCtx() *InvocationContext {
	return invocationContextMap[goid.Get()]
}
