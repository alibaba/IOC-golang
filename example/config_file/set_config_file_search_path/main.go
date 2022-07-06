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
	"github.com/alibaba/ioc-golang/config"
	configField "github.com/alibaba/ioc-golang/extension/config"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ConfigValue *configField.ConfigString `config:",config.app.config-value"`
}

func (a *App) Run() {
	fmt.Printf("Load '%s' from config file\n", a.ConfigValue.Value())
}

func main() {
	// config.WithSearchPath or set env 'IOC_GOLANG_CONFIG_SEARCH_PATH'
	if err := ioc.Load(config.WithSearchPath("./custom_config_path")); err != nil {
		panic(err)
	}

	app, err := GetAppIOCInterfaceSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
