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
	"github.com/alibaba/ioc-golang/extension/config_center/nacos"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init

type App struct {
	NormalNacosClient  nacos.ConfigClientIOCInterface `normal:""`
	NormalNacosClient2 nacos.ConfigClientIOCInterface `normal:",nacos-2"`

	createByAPINacosConfigCenterClient nacos.ConfigClientIOCInterface
}

func (a *App) Run() {
	setAndGetConfig(a.NormalNacosClient, "data1", "group1", "mycontent1")
	setAndGetConfig(a.NormalNacosClient2, "data2", "group2", "mycontent2")
	setAndGetConfig(a.createByAPINacosConfigCenterClient, "data3", "group3", "mycontent3")
}

type Param struct {
	NacosAddr string
	NacosPort int
}

func (p *Param) Init(a *App) (*App, error) {
	// create nacos client with api
	createByAPINacosConfigCenterClient, err := nacos.GetConfigClientIOCInterface(&nacos.Param{
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
	a.createByAPINacosConfigCenterClient = createByAPINacosConfigCenterClient
	return a, nil
}

func setAndGetConfig(client nacos.ConfigClientIOCInterface, dataID, group, content string) {
	_, err := client.PublishConfig(vo.ConfigParam{
		DataId:  dataID,
		Group:   group,
		Content: content,
	})
	if err != nil {
		panic(err)
	}

	configContent, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\n==== Set and Get config %s from nacos", configContent)
}

func main() {
	if err := ioc.Load(
		config.WithSearchPath("../conf")); err != nil {
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
