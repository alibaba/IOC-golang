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

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	// create table
	assert.Nil(t, a.MyDB.Model(&MyDataDO{}).AutoMigrate(&MyDataDO{}))
	toInsertMyData := &MyDataDO{
		Value: "first value",
	}
	assert.Nil(t, a.MyDB.Model(&MyDataDO{}).Create(toInsertMyData).Error())
	myDataDOs := make([]MyDataDO, 0)
	assert.Nil(t, a.MyDB.Model(&MyDataDO{}).Where("id = ?", 1).Find(&myDataDOs).Error())
	assert.Equal(t, 1, len(myDataDOs))
	assert.Equal(t, int32(1), myDataDOs[0].Id)
	assert.Equal(t, "first value", myDataDOs[0].Value)
}

func TestGORM(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		log.Println("Warning: Mysql image only support amd arch. Skip integration test")
		return
	}
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", time.Second*10))
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
