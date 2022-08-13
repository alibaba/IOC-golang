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
	"time"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/example/aop/dynamic_plugin/complex/service2"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	Service2OwnInterface service2.Service2IOCInterface `singleton:""`
}

func (a *App) Run() {
	// origin default name
	fmt.Println(a.Service2OwnInterface.GetName())
	fmt.Println(a.Service2OwnInterface.GetService1Normal().GetName())
	fmt.Println(a.Service2OwnInterface.GetService1Singleton().GetName())

	// init name
	a.Service2OwnInterface.SetName("service2")
	a.Service2OwnInterface.GetService1Singleton().SetName("service1 singleton")
	a.Service2OwnInterface.GetService1Normal().SetName("service1 normal")

	for {
		time.Sleep(time.Second * 3)
		// print name
		fmt.Println(a.Service2OwnInterface.GetName())
		fmt.Println(a.Service2OwnInterface.GetService1Normal().GetName())
		fmt.Println(a.Service2OwnInterface.GetService1Singleton().GetName())

		// set data
		a.Service2OwnInterface.SetData("value")
		fmt.Println("get loaded data = ", a.Service2OwnInterface.LoadData())
	}
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
