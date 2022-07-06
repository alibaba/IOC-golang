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

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/config"
	configField "github.com/alibaba/ioc-golang/extension/config"
	"github.com/alibaba/ioc-golang/extension/state/redis"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ConfigValue              *configField.ConfigString `config:",config.app.config-value"`
	ConfigValueFromEnv       *configField.ConfigString `config:",config.app.config-value-from-env"`
	NestedConfigValue        *configField.ConfigString `config:",config.app.nested-config-value"`
	NestedConfigValueFromEnv *configField.ConfigString `config:",config.app.nested-config-value-from-env"`

	RedisClient redis.RedisIOCInterface `singleton:""`
}

func (a *App) Run() {
	fmt.Printf("Load '%s' from config file\n", a.ConfigValue.Value())
	fmt.Printf("Load '%s' from env\n", a.ConfigValueFromEnv.Value())
	fmt.Printf("Load nested value '%s' from config file\n", a.NestedConfigValue.Value())
	fmt.Printf("Load nested value '%s' from env\n", a.NestedConfigValueFromEnv.Value())
	if err := a.RedisClient.Ping().Err(); err != nil {
		panic(err)
	}
}

func main() {
	// start
	os.Setenv("MY_CONFIG_ENV_KEY", "myEnvValue")
	os.Setenv("REDIS_ADDRESS", "127.0.0.1:6379")
	if err := ioc.Load(
		config.WithSearchPath("./redis_config", "./my_custom_config_path"),
		config.WithProfilesActive("pro")); err != nil {
		panic(err)
	}

	app, err := GetAppIOCInterfaceSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
