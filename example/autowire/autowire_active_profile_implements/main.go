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
	"github.com/alibaba/ioc-golang/example/autowire/autowire_active_profile_implements/service"
	_ "github.com/alibaba/ioc-golang/example/autowire/autowire_active_profile_implements/service/impl"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:type=normal

type App struct {
	ServiceImpl service.Service `singleton:""` // inject implement ptr with active profile like 'dev'
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(a.ServiceImpl.GetHelloString("laurence"))
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
