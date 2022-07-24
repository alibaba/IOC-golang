# IOC-golang: A golang dependency injection framework

```
  ___    ___     ____                           _                         
 |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                |___/                               |___/ 
```

[![IOC-golang CI](https://github.com/alibaba/IOC-golang/actions/workflows/github-actions.yml/badge.svg)](https://github.com/alibaba/IOC-golang/actions/workflows/github-actions.yml)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

### English | [中文](./README_CN.md)

![demo gif](https://raw.githubusercontent.com/ioc-golang/ioc-golang-website/main/resources/video/ioc-golang-demo.gif)

IOC-golang is a powerful golang dependency injection framework that provides a complete implementation of IoC containers. Its capabilities are as follows:

- [Dependency Injection](https://ioc-golang.github.io/docs/getting-started/tutorial/)

  Supports dependency injection of any structure and interface, we also support object life cycle management mechanism.

  Can take over object creation, parameter injection, factory methods. Customizable object parameter source.

- [Struct Proxy](./example/aop)

  Based on the idea of AOP, we provide struct proxy layer for all struct registered to ioc-golang. In the scene of interface oriented development, we can use many devlops features based on the extenablility of this proxy AOP layer. Such as interface listing, param value watching, method level tracing, performance badpoint analysis, fault injection, method level tracing in distributed system and so on.

- [Automatic struct descriptor codes generation capability](./iocli)

  We provide a code generation tool, and developers can annotate the structure through annotations, so as to easily generate structure registration code.

- [Scalability](./extension)

  Support the extension of struct to be injected, the extension of autowire model, and the extension of the debug AOP layer.

- [Many pre-defined components](./example)

  Provides pre-defined objects and middleware sdk for injection directly.

## Project Structure

- **aop:** Debug module: Provide debugging API, provide debugging injection layer basic  implementation and extendable API.
- **autowire:** Provides two basic injection models: singleton model and multi-instance model
- **config:** Configuration loading module, responsible for parsing ion-golang's configuration files
- **extension:** Component extension directory: Provides preset implementation structures based on various domain. Such as database, cache, pubs.
- **example:** example repository
- **iocli:** code generation/program debugging tool

## Quick start

### Install code generation tools

```shell
% go install github.com/alibaba/ioc-golang/iocli@latest
% iocli
hello
````

### Dependency Injection Tutorial

We will develop a project with the following topology, This tutorial can show:

1. Registry codes generation
2. Interface injection
3. Struct pointer injection
4. Get object by API
5. Debug capability, list interface, implementations and methods; watch real-time param and return value.

![ioc-golang-quickstart-structure](https://raw.githubusercontent.com/ioc-golang/ioc-golang-website/main/resources/img/ioc-golang-quickstart-structure.png)


All the code the user needs to write: main.go

```go
package main

import (
	"fmt"
	"time"

	"github.com/alibaba/ioc-golang"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	// inject main.ServiceImpl1 pointer to Service interface with proxy wrapper
	ServiceImpl1 Service `singleton:"main.ServiceImpl1"`

	// inject main.ServiceImpl2 pointer to Service interface with proxy wrapper
	ServiceImpl2 Service `singleton:"main.ServiceImpl2"`

	// inject ServiceImpl1 pointer to Service1 's own interface with proxy wrapper
	// this interface belongs to ServiceImpl1, there is no need to mark 'main.ServiceImpl1' in tag
	Service1OwnInterface ServiceImpl1IOCInterface `singleton:""`

	// inject ServiceStruct struct pointer
	ServiceStruct *ServiceStruct `singleton:""`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		fmt.Println(a.ServiceImpl1.GetHelloString("laurence"))
		fmt.Println(a.ServiceImpl2.GetHelloString("laurence"))

		fmt.Println(a.Service1OwnInterface.GetHelloString("laurence"))
		
		fmt.Println(a.ServiceStruct.GetString("laurence"))
	}
}

type Service interface {
	GetHelloString(string) string
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceImpl1 struct {
}

func (s *ServiceImpl1) GetHelloString(name string) string {
	return fmt.Sprintf("This is ServiceImpl1, hello %s", name)
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

type ServiceStruct struct {
}

func (s *ServiceStruct) GetString(name string) string {
	return fmt.Sprintf("This is ServiceStruct, hello %s", name)
}

func main() {
	// start to load all structs
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// Get Struct
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}


```
The proxy wrapped layer mentioned above, is a proxy layer injected by ioc-golang by default, when developer want to inject an object to interface field, or get with interface by API. Inject to interface is recommended by us. Every object injected with proxy wrapped layer would have devops feature.

After writing, you can exec the following cli command to init go mod and generate codes.  (mac may require sudo due to permissions during code generation)

```bash
% go mod init ioc-golang-demo
% export GOPROXY="https://goproxy.cn"
% go mod tidy
% go get github.com/alibaba/ioc-golang@master
% sudo iocli gen
````

It will be generated in the current directory: zz_generated.ioc.go, developers **do not need to care about this file**, 'GetAppSingleton' method mentioned above is defined in generated code.

```go
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli

package main

import (
        autowire "github.com/alibaba/ioc-golang/autowire"
        normal "github.com/alibaba/ioc-golang/autowire/normal"
        "github.com/alibaba/ioc-golang/autowire/singleton"
        util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
        normal.RegisterStructDescriptor(&autowire.StructDescriptor{
                Factory: func() interface{} {
                        return &app_{}
                },
        })
        singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
                Factory: func() interface{} {
                        return &App{}
                },
        })
  ...
func GetServiceStructIOCInterface() (ServiceStructIOCInterface, error) {
        i, err := singleton.GetImplWithProxy(util.GetSDIDByStructPtr(new(ServiceStruct)), nil)
        if err != nil {
                return nil, err
        }
        impl := i.(ServiceStructIOCInterface)
        return impl, nil
}


```

See the file tree:

```bash
% tree
.
├── go.mod
├── go.sum
├── main.go
└── zz_generated.ioc.go

0 directories, 4 files
````

#### Execute program

`go run .`

Console printout:

```sh
  ___    ___     ____                           _                         
 |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                |___/                               |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golang config
[Config] Load default config file from ../conf/ioc_golang.yaml
[Config] Load ioc-golang config file failed. open /Users/laurence/Desktop/workplace/alibaba/conf/ioc_golang.yaml: no such file or directory
 The load procedure is continue
[Boot] Start to load debug
[Debug] Debug port is set to default :1999
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type normal
[Autowire Struct Descriptor] Found type normal registered SD main.serviceStruct_
[Autowire Struct Descriptor] Found type normal registered SD main.app_
[Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl1_
[Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl2_
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl1
[Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl2
[Autowire Struct Descriptor] Found type singleton registered SD main.ServiceStruct
[Debug] Debug server listening at :1999
This is ServiceImpl1, hello laurence
This is ServiceImpl2, hello laurence
This is ServiceImpl1, hello laurence
This is ServiceStruct, hello laurence
...
```

It shows that the injection is successful and the program runs normally.

**Debug the app**

Following logs can be found in console output:

```bash
[Debug] Debug server listening at :1999
```

Open a new console, use iocli 's debug feature to list all structs with proxy layer, and their methods. Default port is 1999.

```
% iocli list
main.ServiceImpl1
[GetHelloString]

main.ServiceImpl2
[GetHelloString]
```

Watch real-time param and return value. We take  main.ServiceImpl 's 'GetHelloString' method as an example. The method would be called twice every 3s :

```bash
% iocli watch main.ServiceImpl1 GetHelloString
========== On Call ==========
main.ServiceImpl1.GetHelloString()
Param 1: (string) (len=8) "laurence"

========== On Response ==========
main.ServiceImpl1.GetHelloString()
Response 1: (string) (len=36) "This is ServiceImpl1, hello laurence"

========== On Call ==========
main.ServiceImpl1.GetHelloString()
Param 1: (string) (len=8) "laurence"

========== On Response ==========
main.ServiceImpl1.GetHelloString()
Response 1: (string) (len=36) "This is ServiceImpl1, hello laurence"
...
```



### Annotation Analysis

````go
// +ioc:autowire=true
The code generation tool recognizes objects marked with the +ioc:autowire=true annotation

// +ioc:autowire:type=singleton
The marker autowire model is the singleton
````

###  More

[Docs](https://ioc-golang.github.io/cn)

More code generation annotations can be viewed at [iocli](https://github.com/alibaba/IOC-golang/tree/master/iocli).

You can go to [ioc-golang/example](https://github.com/alibaba/IOC-golang/tree/master/example) for more examples and advanced usage.

You can go to [E-commercial system demo based on ioc-golang](https://github.com/ioc-golang/shopping-system) to refer to applications system on distributed scene.

### License

IOC-golang developed by Alibaba and licensed under the Apache License (Version 2.0).
See the NOTICE file for more information.

### Connect with us

Welcome to join dingtalk group 44638289 if you are interested with the project.

<div align="center">
	<img src="https://github.com/ioc-golang/ioc-golang-website/blob/main/resources/img/dingtalk_group.png?raw=true" width="30%">
</div>

### Star me please ⭐

If you think this project is interesting, or helpful to you, please give a star!
