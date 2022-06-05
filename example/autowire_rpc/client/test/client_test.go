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

package test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/example/autowire_rpc/client/test/dto"
	_ "github.com/alibaba/ioc-golang/example/autowire_rpc/client/test/service"
	"github.com/alibaba/ioc-golang/example/autowire_rpc/client/test/service/api"
	"github.com/alibaba/ioc-golang/extension/autowire/rpc/rpc_client"
)

func TestMain(m *testing.M) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	m.Run()
	os.Exit(0)
}

var param = &rpc_client.Param{
	Address: "127.0.0.1:2022",
}

func TestRPCClientGetByAPI(t *testing.T) {
	simpleClient, err := rpc_client.GetImpl("github.com/alibaba/ioc-golang/example/autowire_rpc/client/test/service/api.SimpleServiceIOCRPCClient", param)
	assert.Nil(t, err)
	testSimpleClient(t, simpleClient.(*api.SimpleServiceIOCRPCClient))

	simpleClient2, err := rpc_client.ImplClientStub(&api.SimpleServiceIOCRPCClient{}, param)
	assert.Nil(t, err)
	testSimpleClient(t, simpleClient2.(*api.SimpleServiceIOCRPCClient))
}

func testSimpleClient(t *testing.T, simpleClientImpl *api.SimpleServiceIOCRPCClient) {
	usr, err := simpleClientImpl.GetUser("laurence", 23)
	assert.Nil(t, err)
	assertUser(t, usr)
}

func assertUser(t *testing.T, usr *dto.User) {
	assert.NotNil(t, usr)
	assert.Equal(t, 1, usr.Id)
	assert.Equal(t, "laurence", usr.Name)
	assert.Equal(t, 23, usr.Age)
}

func assertComplexStruct(t *testing.T, complexStruct *dto.CustomStruct) {
	assert.NotNil(t, complexStruct)
	assert.Equal(t, int64(1), complexStruct.CustomStructId)
	assert.Equal(t, int64(1), *complexStruct.IdPtr)
	assert.Equal(t, "laurence", complexStruct.CustomStructName)
	assert.Equal(t, "laurence", *complexStruct.NamePtr)
	assertUser(t, complexStruct.SubStructPtr)
	assertUser(t, &complexStruct.SubStruct)
	assert.Equal(t, 2, len(complexStruct.SubStructSlice))
	for _, v := range complexStruct.SubStructSlice {
		assertUser(t, &v)
	}
	assert.Equal(t, 2, len(complexStruct.SubStructPtrSlice))
	for _, v := range complexStruct.SubStructPtrSlice {
		assertUser(t, v)
	}
	assert.Equal(t, 1, len(complexStruct.CustomSubStructPtrMap))
	for _, v := range complexStruct.CustomSubStructPtrMap {
		assertUser(t, v)
	}
	assert.Equal(t, 1, len(complexStruct.CustomIntMap))
	for _, v := range complexStruct.CustomIntMap {
		assert.Equal(t, 1, v)
	}
	assert.Equal(t, 1, len(complexStruct.CustomSubStructMap))
	for _, v := range complexStruct.CustomSubStructMap {
		assertUser(t, &v)
	}
	assert.Equal(t, 1, len(complexStruct.CustomStringMap))
	for _, v := range complexStruct.CustomStringMap {
		assert.Equal(t, "value", v)
	}
	assert.Equal(t, 2, len(complexStruct.StringSlice))

	// assert combination
	assertUser(t, &complexStruct.User)
}

func newUser() *dto.User {
	return &dto.User{
		Name: "laurence",
		Age:  23,
		Id:   1,
	}
}
func newComplexStruct() *dto.CustomStruct {
	idPtr := int64(1)
	name := "laurence"
	return &dto.CustomStruct{
		User:             *newUser(),
		CustomStructId:   1,
		IdPtr:            &idPtr,
		CustomStructName: name,
		NamePtr:          &name,
		CustomStringMap: map[string]string{
			"key": "value",
		},
		CustomIntMap: map[string]int{
			"key": 1,
		},
		CustomSubStructMap: map[string]dto.User{
			"key": *newUser(),
		},
		CustomSubStructPtrMap: map[string]*dto.User{
			"key": newUser(),
		},
		StringSlice: []string{
			"value1", "value2",
		},
		SubStructSlice: []dto.User{
			*newUser(), *newUser(),
		},
		SubStructPtrSlice: []*dto.User{
			newUser(), newUser(),
		},
		SubStruct:    *newUser(),
		SubStructPtr: newUser(),
	}
}

func TestComplexRPC(t *testing.T) {
	complexClient, err := rpc_client.ImplClientStub(&api.ComplexServiceIOCRPCClient{}, param)
	assert.Nil(t, err)
	complexClientImpl := complexClient.(*api.ComplexServiceIOCRPCClient)

	testRPCBasicType(t, complexClientImpl)
	testRPCWithoutSomething(t, complexClientImpl)
	testRPCWithCustomStruct(t, complexClientImpl)
	testRPCWithError(t, complexClientImpl)
	testRPCWithParamCustomMethod(t, complexClientImpl)
}

func testRPCWithParamCustomMethod(t *testing.T, client *api.ComplexServiceIOCRPCClient) {
	usr := client.RPCWithParamCustomMethod(*newComplexStruct())
	assertUser(t, &usr)
}

func testRPCWithError(t *testing.T, client *api.ComplexServiceIOCRPCClient) {
	usr, err := client.RPCWithError()
	assertUser(t, usr)
	assert.NotNil(t, err)
	assert.Equal(t, "custom error = custom", err.Error())
}

func testRPCWithCustomStruct(t *testing.T, client *api.ComplexServiceIOCRPCClient) {
	customStruct := newComplexStruct()
	rspCustomStruct, rspCustomStructPtr := client.RPCWithCustomValue(*customStruct, customStruct)
	assertComplexStruct(t, &rspCustomStruct)
	assertComplexStruct(t, rspCustomStructPtr)
}

func testRPCWithoutSomething(t *testing.T, client *api.ComplexServiceIOCRPCClient) {
	client.RPCWithoutParamAndReturnValue()
	client.RPCWithoutReturnValue(&dto.User{})
	usr, err := client.RPCWithoutParam()
	assert.Nil(t, err)
	assertUser(t, usr)
}

func testRPCBasicType(t *testing.T, client *api.ComplexServiceIOCRPCClient) {
	name := "laurence"
	age := 23
	age32 := int32(23)
	age64 := int64(23)
	ageF32 := float32(23)
	ageF64 := float64(23)
	rspName, rspAge, rspAge32, rspAge64, rspAgeF32, rspAgeF64, rspNamePtr, rspAgePtr, rspAge32Ptr, rspAge64Ptr,
		rspAgeF32Ptr, rspAgeF64Ptr := client.RPCBasicType(name, age, age32, age64, ageF32, ageF64, &name,
		&age, &age32, &age64, &ageF32, &ageF64)
	assert.Equal(t, name, rspName)
	assert.Equal(t, name, *rspNamePtr)
	assert.Equal(t, age, rspAge)
	assert.Equal(t, age, *rspAgePtr)
	assert.Equal(t, age32, rspAge32)
	assert.Equal(t, age32, *rspAge32Ptr)
	assert.Equal(t, age64, rspAge64)
	assert.Equal(t, age64, *rspAge64Ptr)
	assert.Equal(t, ageF32, rspAgeF32)
	assert.Equal(t, ageF32, *rspAgeF32Ptr)
	assert.Equal(t, ageF64, rspAgeF64)
	assert.Equal(t, ageF64, *rspAgeF64Ptr)
}
