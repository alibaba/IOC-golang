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

	"github.com/alibaba/ioc-golang/config"

	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/extension/config_center/nacos"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

const (
	ipFoo   = "127.0.0.1"
	portFoo = 1999
)

func (a *App) TestRun(t *testing.T) {
	testSetAndGetConfig(t, a.NormalNacosClient, "data1", "group1", "mycontent1")
	testSetAndGetConfig(t, a.NormalNacosClient2, "data2", "group2", "mycontent2")
	testSetAndGetConfig(t, a.createByAPINacosConfigCenterClient, "data3", "group3", "mycontent3")
}

func testSetAndGetConfig(t *testing.T, client nacos.ConfigClientIOCInterface, dataID, group, content string) {
	_, err := client.PublishConfig(vo.ConfigParam{
		DataId:  dataID,
		Group:   group,
		Content: content,
	})
	assert.Nil(t, err)

	configContent, err := client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	assert.Nil(t, err)
	assert.Equal(t, content, configContent)
}

func TestNacosClient(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		log.Println("Warning: Nacos image only support amd arch. Skip integration test")
		return
	}
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", time.Second*10))
	assert.Nil(t, ioc.Load(
		config.WithSearchPath("../conf")))
	app, err := GetAppSingleton(&Param{
		NacosPort: 8848,
		NacosAddr: "localhost",
	})
	assert.Nil(t, err)
	app.TestRun(t)

	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
