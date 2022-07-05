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
	"os"

	"github.com/go-redis/redis"

	"github.com/alibaba/ioc-golang"

	normalRedis "github.com/alibaba/ioc-golang/extension/state/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=Init
// +ioc:autowire:alias=AppAlias

type App struct {
	NormalRedis    normalRedis.RedisIOCInterface `normal:""`
	NormalDB1Redis normalRedis.RedisIOCInterface `normal:",db1-redis"`
	NormalDB2Redis normalRedis.RedisIOCInterface `normal:",db2-redis"`
	NormalDB3Redis normalRedis.RedisIOCInterface `normal:",address=127.0.0.1:6379&db=3"`
	NormalDB4Redis normalRedis.RedisIOCInterface `normal:",address=${REDIS_ADDRESS_EXPAND}&db=5"`
	NormalDB5Redis normalRedis.RedisIOCInterface `normal:",address=${autowire.normal.<github.com/alibaba/ioc-golang/extension/state/redis.Redis>.nested.address}&db=15"`

	privateClient *redis.Client
}

type Param struct {
	RedisAddr string
}

func (p *Param) Init(a *App) (*App, error) {
	privateClient := redis.NewClient(&redis.Options{
		Addr: p.RedisAddr,
	})
	a.privateClient = privateClient
	return a, nil
}

func (a *App) Run() {
	if _, err := a.NormalRedis.Set("mykey", "db0", -1).Result(); err != nil {
		panic(err)
	}

	if _, err := a.NormalDB1Redis.Set("mykey", "db1", -1).Result(); err != nil {
		panic(err)
	}

	if _, err := a.NormalDB2Redis.Set("mykey", "db2", -1).Result(); err != nil {
		panic(err)
	}

	if _, err := a.NormalDB3Redis.Set("mykey", "db3", -1).Result(); err != nil {
		panic(err)
	}
	if _, err := a.NormalDB4Redis.Set("mykey", "db15", -1).Result(); err != nil {
		panic(err)
	}
	if _, err := a.NormalDB5Redis.Set("mykey", "db15", -1).Result(); err != nil {
		panic(err)
	}

	val1, err := a.NormalRedis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client0 get ", val1)

	val2, err := a.NormalDB1Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client1 get ", val2)

	val3, err := a.NormalDB2Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client2 get ", val3)

	val4, err := a.NormalDB3Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client3 get ", val4)

	val5, err := a.NormalDB4Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client4 get ", val5)

	val6, err := a.NormalDB5Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client5 get ", val6)
}

func init() {
	err := os.Setenv("REDIS_ADDRESS_EXPAND", "localhost:6379")
	for err != nil {
		err = os.Setenv("REDIS_ADDRESS_EXPAND", "localhost:6379")
	}
}

func main() {
	if err := ioc.Load(
		config.WithSearchPath("../conf"),
		config.WithConfigName("ioc_golang")); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton(&Param{
		RedisAddr: "localhost:6379",
	})
	if err != nil {
		panic(err)
	}

	app.Run()
}
