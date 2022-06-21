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
	"os"
	"path/filepath"

	conf "github.com/alibaba/ioc-golang/config"
	"github.com/go-redis/redis"

	"github.com/alibaba/ioc-golang"
	normalRedis "github.com/alibaba/ioc-golang/extension/normal/redis"
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

	privateClient *redis.Client
}

type Param struct {
	RedisAddr string
}

func init() {
	_ = os.Setenv("REDIS_ADDRESS_EXPAND", "localhost:6379")
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

	if _, err := a.NormalDB4Redis.Set("mykey", "db5", -1).Result(); err != nil {
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
}

func main() {
	wd, _ := os.Getwd()
	absPathOpt := conf.WithAbsPath(filepath.Join(wd, "./example/autowire_redis_client/conf/ioc_golang.yaml"))

	if err := ioc.Load(absPathOpt); err != nil {
		panic(err)
	}
	app, err := GetApp(&Param{
		RedisAddr: "localhost:6379",
	})
	if err != nil {
		panic(err)
	}

	app.Run()
}
