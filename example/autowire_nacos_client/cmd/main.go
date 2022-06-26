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
	"fmt"

	"github.com/alibaba/ioc-golang/config"

	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/extension/registry/nacos"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init
// +ioc:autowire:alias=AppAlias

type App struct {
	NormalNacosClient  nacos.NamingClientIOCInterface `normal:""`
	NormalNacosClient2 nacos.NamingClientIOCInterface `normal:",nacos-2"`

	createByAPINacosClient nacos.NamingClientIOCInterface
}

func (a *App) Run() {
	getAndSetService(a.NormalNacosClient, "normal-autowire-client-ioc-golang-debug-service")
	getAndSetService(a.NormalNacosClient2, "normal-autowire-client2-ioc-golang-debug-service")
	getAndSetService(a.createByAPINacosClient, "createByAPINacosClient-ioc-golang-debug-service")
}

type Param struct {
	NacosAddr string
	NacosPort int
}

func (p *Param) Init(a *App) (*App, error) {
	// create nacos client with api
	createByAPINacosClient, err := nacos.GetNamingClientIOCInterface(&nacos.Param{
		NacosClientParam: vo.NacosClientParam{
			ServerConfigs: []constant.ServerConfig{
				{
					IpAddr: p.NacosAddr,
					Port:   uint64(p.NacosPort),
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	a.createByAPINacosClient = createByAPINacosClient
	return a, nil
}

func getAndSetService(client nacos.NamingClientIOCInterface, serviceName string) {
	_, err := client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "127.0.0.1",
		Port:        1999,
		ServiceName: serviceName,
	})
	if err != nil {
		panic(err)
	}

	service, err := client.GetService(vo.GetServiceParam{
		ServiceName: serviceName,
	})
	if err != nil {
		return
	}
	fmt.Printf("\n\n==== Get service %s from nacos: %+v", serviceName, service)
}

func main() {
	if err := ioc.Load(
		config.WithSearchPath("../conf"),
		config.WithConfigName("ioc_golang")); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton(&Param{
		NacosAddr: "localhost",
		NacosPort: 8848,
	})
	if err != nil {
		panic(err)
	}

	app.Run()
}
