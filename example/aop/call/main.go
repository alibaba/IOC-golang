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

	"github.com/alibaba/ioc-golang/example/aop/call/dto"

	"github.com/alibaba/ioc-golang"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	// inject ServiceImpl1 pointer to Service1 's own interface with proxy wrapper
	// this interface belongs to ServiceImpl1
	Service1OwnInterface ServiceImpl1IOCInterface `singleton:""`

	UserService UserServiceIOCInterface `singleton:""`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(a.Service1OwnInterface.GetHelloString("laurence"))
		user, _ := a.UserService.CreateUser("laurence", 23)
		fmt.Println("create user", user)
	}
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
	Service2 ServiceImpl2IOCInterface `singleton:""`
}

func (s *ServiceImpl1) GetHelloString(name string) string {
	return s.Service2.GetHelloString(name)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl2 struct {
}

func (s *ServiceImpl2) GetHelloString(name string) string {
	return fmt.Sprintf("This is ServiceImpl2, hello %s", name)
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type UserService struct {
	mark string
}

func (s *UserService) SetMark(mark string) {
	s.mark = mark
}

func (s *UserService) CreateUser(name string, age int) (*dto.User, error) {
	return &dto.User{
		Id:   1,
		Name: name,
		Age:  age,
		Mark: s.mark,
	}, nil
}

func (s *UserService) ParseUserInfo(usr *dto.User) (string, int, string, error) {
	return usr.Name, usr.Age, usr.Mark, nil
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppIOCInterfaceSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
