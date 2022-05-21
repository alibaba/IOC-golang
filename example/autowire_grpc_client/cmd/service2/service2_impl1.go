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

package service2

import (
	"context"

	"github.com/alibaba/ioc-golang/example/autowire_grpc_client/api"
)

// +ioc:autowire=true
// +ioc:autowire:interface=Service2
// +ioc:autowire:type=singleton

type Impl1 struct {
	HelloServiceClient api.HelloServiceClient `grpc:"hello-service"`
}

func (i *Impl1) Hello(string) string {
	rsp, err := i.HelloServiceClient.SayHello(context.Background(), &api.HelloRequest{
		Name: "laurence_service2_impl1",
	})
	if err != nil {
		panic(err)
	}
	return rsp.Reply
}
