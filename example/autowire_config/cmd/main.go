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

	ioc "github.com/alibaba/IOC-Golang"
	"github.com/alibaba/IOC-Golang/autowire/singleton"
	"github.com/alibaba/IOC-Golang/extension/config"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
type App struct {
	DemoConfigString *config.ConfigString `config:"ConfigString,autowire.config.demo-config.string-value"`
	DemoConfigInt    *config.ConfigInt    `config:"ConfigInt,autowire.config.demo-config.int-value"`
	DemoConfigMap    *config.ConfigMap    `config:"ConfigMap,autowire.config.demo-config.map-value"`
	DemoConfigSlice  *config.ConfigSlice  `config:"ConfigSlice,autowire.config.demo-config.slice-value"`
}

func (a *App) Run() {
	fmt.Println(a.DemoConfigString.Value())
	fmt.Println(a.DemoConfigInt.Value())
	fmt.Println(a.DemoConfigMap.Value())
	fmt.Println(a.DemoConfigSlice.Value())
}

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil {
		panic(err)
	}
	app := appInterface.(*App)
	app.Run()
}
