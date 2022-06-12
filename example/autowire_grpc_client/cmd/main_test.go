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

package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/example/autowire_grpc_client/api"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	name := "laurence"
	rsp, err := a.HelloServiceClient.SayHello(context.Background(), &api.HelloRequest{
		Name: name,
	})
	assert.Nil(t, err)

	assert.Equal(t, "Hello laurence", rsp.Reply)
	assert.Equal(t, "Hello laurence_service1_impl1", a.ExampleService1Impl1.Hello(name+"_service1_impl1"))
	assert.Equal(t, "Hello laurence_service2_impl1", a.ExampleService2Impl1.Hello(name+"_service2_impl1"))
	assert.Equal(t, "Hello laurence_service2_impl2", a.ExampleService2Impl2.Hello(name+"_service2_impl2"))
	assert.Equal(t, "Hello laurence_struct", a.ExampleStruct1.Hello(name+"_struct"))
}

func TestGRPC(t *testing.T) {
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", 0))
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetApp()
	if err != nil {
		panic(err)
	}

	app.TestRun(t)
	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
