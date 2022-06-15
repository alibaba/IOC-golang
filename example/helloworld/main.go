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
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceImpl1  Service        `singleton:"main.ServiceImpl1"` // inject Service 's ServiceImpl1 implementation
	ServiceImpl2  Service        `singleton:"main.ServiceImpl2"` // inject Service 's ServiceImpl2 implementation
	ServiceStruct *ServiceStruct `singleton:""`                  // inject ServiceStruct struct pointer
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		a.ServiceImpl1.Hello()
		a.ServiceImpl2.Hello()

		fmt.Println(a.ServiceStruct.GetString("laurence"))
	}
}

type Service interface {
	Hello()
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
}

func (s *ServiceImpl1) Hello() {
	fmt.Println("This is ServiceImpl1, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl2 struct {
}

func (s *ServiceImpl2) Hello() {
	fmt.Println("This is ServiceImpl2, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceStruct struct {
}

func (s *ServiceStruct) GetString(name string) string {
	return fmt.Sprintf("Hello %s", name)
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetApp()
	if err != nil {
		panic(err)
	}
	app.Run()
}
