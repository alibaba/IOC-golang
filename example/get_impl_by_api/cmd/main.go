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

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/extension/normal/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:alias=appalias

type App struct {
}

func (a *App) Run() {
	redisClientGetyByNormalAPI, err := normal.GetImpl("github.com/alibaba/ioc-golang/extension/normal/redis.Impl", &redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
	if err != nil {
		panic(err)
	}
	redisClientGetByNormalAPIImpl := redisClientGetyByNormalAPI.(redis.Redis)
	_, err = redisClientGetByNormalAPIImpl.Set("myKey", "myValue", -1)
	if err != nil {
		panic(err)
	}
	fmt.Sprintln("redisClientByNormalAPIImpl set  myKey:myValue")

	redisClientGetByRedisExtension, err := redis.GetRedis(&redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
	if err != nil {
		panic(err)
	}
	val, err := redisClientGetByRedisExtension.Get("myKey")
	if err != nil {
		panic(err)
	}
	fmt.Println("get val = ", val)
}

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetApp()
	if err != nil {
		panic(err)
	}
	app.Run()
}
