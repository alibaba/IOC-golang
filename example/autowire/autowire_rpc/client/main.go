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
	"github.com/alibaba/ioc-golang/example/autowire/autowire_rpc/server/pkg/service/api"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceStruct api.ServiceStructIOCRPCClient `rpc-client:",address=127.0.0.1:2022"`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		usr, err := a.ServiceStruct.GetUser("laurence", 23)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("get user = %+v\n", usr)
	}
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// 'App' is alias name
	// We can get instance by ths id
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
