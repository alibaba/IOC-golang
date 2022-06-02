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
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	normalNacos "github.com/alibaba/ioc-golang/extension/normal/nacos"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

const (
	ipFoo   = "127.0.0.1"
	portFoo = 1999
)

func (a *App) TestRun(t *testing.T) {
	testGetAndSetService(t, a.NormalNacosClient, "normal-autowire-client-ioc-golang-debug-service")
	testGetAndSetService(t, a.NormalNacosClient2, "normal-autowire-client2-ioc-golang-debug-service")
	testGetAndSetService(t, a.createByAPINacosClient, "createByAPINacosClient-ioc-golang-debug-service")
}

func testGetAndSetService(t *testing.T, client normalNacos.NacosClient, serviceName string) {
	_, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        1999,
		ServiceName: serviceName,
	})
	retries := 0
	for err != nil && retries < 3 {
		time.Sleep(time.Second * 10)
		_, err = client.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          "127.0.0.1",
			Port:        1999,
			ServiceName: serviceName,
		})
		retries++
	}
	if err != nil {
		panic(err)
	}

	service, err := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return
	}
	assert.Equal(t, serviceName, service.Name)
	assert.Equal(t, 1, len(service.Hosts))
	assert.Equal(t, ipFoo, service.Hosts[0].Ip)
	assert.Equal(t, uint64(portFoo), service.Hosts[0].Port)
}

func TestNacosClient(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		log.Println("Warning: Nacos image only support amd arch. Skip integration test")
		return
	}
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", time.Second*10))
	assert.Nil(t, ioc.Load())
	appInterface, err := singleton.GetImpl("AppAlias")
	assert.Nil(t, err)
	app := appInterface.(*App)
	app.TestRun(t)

	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
