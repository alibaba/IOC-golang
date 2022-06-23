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

	"github.com/alibaba/ioc-golang/config"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/test/docker_compose"
)

func (a *App) TestRun(t *testing.T) {
	_, err := a.NormalRedis.Set("mykey", "db0", -1).Result()
	assert.Nil(t, err)

	_, err = a.NormalDB1Redis.Set("mykey", "db1", -1).Result()
	assert.Nil(t, err)

	_, err = a.NormalDB2Redis.Set("mykey", "db2", -1).Result()
	assert.Nil(t, err)

	_, err = a.NormalDB3Redis.Set("mykey", "db3", -1).Result()
	assert.Nil(t, err)

	val1, err := a.NormalRedis.Get("mykey").Result()
	assert.Nil(t, err)
	assert.Equal(t, "db0", val1)

	val2, err := a.NormalDB1Redis.Get("mykey").Result()
	assert.Nil(t, err)
	assert.Equal(t, "db1", val2)

	val3, err := a.NormalDB2Redis.Get("mykey").Result()
	assert.Nil(t, err)
	assert.Equal(t, "db2", val3)

	val4, err := a.NormalDB3Redis.Get("mykey").Result()
	assert.Nil(t, err)
	assert.Equal(t, "db3", val4)
}

func TestRedisClient(t *testing.T) {
	assert.Nil(t, docker_compose.DockerComposeUp("../docker-compose/docker-compose.yaml", 0))
	if err := ioc.Load(
		config.WithSearchPath("../conf"),
		config.WithConfigName("ioc_golang"),
		config.WithConfigType("yaml")); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton(&Param{
		RedisAddr: "localhost:6379",
	})
	if err != nil {
		panic(err)
	}
	app.TestRun(t)

	assert.Nil(t, docker_compose.DockerComposeDown("../docker-compose/docker-compose.yaml"))
}
