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
	"testing"

	"github.com/alibaba/ioc-golang/autowire/util"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/extension/state/redis"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	redisClientGetByNormalAPIImpl, err := redis.GetRedisIOCInterface(&redis.Param{
		Address: "localhost:6379",
		DB:      "0",
	})
	assert.Nil(t, err)
	_, err = redisClientGetByNormalAPIImpl.Set("myKey", "myValue", -1).Result()
	assert.Nil(t, err)

	redisClientGetByRedisExtension, err := redis.GetRedisIOCInterface(&redis.Param{
		Address: "localhost:6379",
		DB:      "0",
	})
	assert.Nil(t, err)
	val, err := redisClientGetByRedisExtension.Get("myKey").Result()
	assert.Nil(t, err)
	assert.Equal(t, "myValue", val)
}

func TestGetAPI(t *testing.T) {
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", 0))
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

func TestGetSingletonWithProxyAndWithoutProxy(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	assert.Nil(t, err)
	assert.Equal(t, util.GetSDIDByStructPtr(new(App)), util.GetSDIDByStructPtr(app))

	appWithProxy, err := GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	assert.Equal(t, util.GetSDIDByStructPtr(new(app_)), util.GetSDIDByStructPtr(appWithProxy))

	app, err = GetAppSingleton()
	assert.Nil(t, err)
	assert.Equal(t, util.GetSDIDByStructPtr(new(App)), util.GetSDIDByStructPtr(app))
}
