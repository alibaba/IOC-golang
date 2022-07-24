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

	"github.com/alibaba/ioc-golang"

	normalRedis "github.com/alibaba/ioc-golang/extension/state/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	NormalDB1Redis normalRedis.RedisIOCInterface `normal:",db1-redis"`
}

func (a *App) Run() {
	if _, err := a.NormalDB1Redis.Set("mykey", "db1", -1).Result(); err != nil {
		panic(err)
	}

	val2, err := a.NormalDB1Redis.Get("mykey").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("client1 get ", val2)
}

/*
Using API to config the following config properties
```yaml
autowire:
  normal:
    github.com/alibaba/ioc-golang/extension/state/redis.Redis:
      db1-redis:
        param:
          address: localhost:6379
          db: 1
```
*/

func main() {
	if err := ioc.Load(
		config.AddProperty("autowire.normal.<github.com/alibaba/ioc-golang/extension/state/redis.Redis>.db1-redis.param.address", "localhost:6379"),
		config.AddProperty("autowire.normal.<github.com/alibaba/ioc-golang/extension/state/redis.Redis>.db1-redis.param.db", 1)); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}

	app.Run()
}
