# IOC-Golang: A golang dependency injection framework

```
  ___    ___     ____            ____           _                         
 |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                    |___/ 
```

[![IOC-Golang CI](https://github.com/alibaba/IOC-Golang/actions/workflows/github-actions.yml/badge.svg)](https://github.com/alibaba/IOC-Golang/actions/workflows/github-actions.yml)
[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

[IOC-Golang Docs](https://ioc-golang.github.io)

[中文 文档](https://ioc-golang.github.io/cn)

[中文 README](./README_CN.md)

IOC-Golang is a powerful golang dependency injection framework that provides a complete implementation of IoC containers. Its capabilities are as follows:

- [Dependency Injection](https://ioc-golang.github.io/docs/getting-started/turtoral/)

  Supports dependency injection of any structure and interface.

- Object life cycle management mechanism.

  Can take over object creation, parameter injection, factory methods. Customizable object parameter source.

- [Code debugging ability](https://ioc-golang.github.io/docs/examples/debug/)

  Based on the idea of AOP, it provides runtime monitoring and debugging capabilities for object methods taken over by the framework.

- [Automatic struct descriptor codes generation capability](https://ioc-golang.github.io/docs/reference/ioc-go-cli/)

  We provide a code generation tool, and developers can annotate the structure through annotations, so as to easily generate structure registration code.

- [Scalability](https://ioc-golang.github.io/cn/docs/contribution-guidelines/)

  Support the extension of struct to be injected, the extension of autowire model, and the extension of the debug AOP layer.

- [Complete prefabricated components](https://ioc-golang.github.io/cn/docs/examples/)

  Provides prefabricated objects covering mainstream middleware sdk for injection.

## Project Structure

- **autowire:** Provides two basic injection models: singleton model and multi-instance model
- **config:** Configuration loading module, responsible for parsing user yaml configuration files
- **debug:** Debug module: Provide debugging API, provide debugging injection layer implementation
- **extension:** Component extension directory: Provides preset implementation structures based on various injection models:

    - autowire: autoload model extensions

        - grpc: grpc client model definition

        - config: configure the model definition

    - config: configuration injection model extension structure

        - string,int,map,slice

    - normal: multi-instance model extension structure

        - redis

        - mysql

        - rocketmq

        - nacos

    - singleton: singleton model extension structure

        - http-server

- **example:** example repository

- **ioc-go-cli:** code generation/program debugging tool

## Quick start

### Install code generation tools

```shell
go install github.com/alibaba/ioc-golang/ioc-go-cli@latest
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
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceImpl1 Service `singleton:"ServiceImpl1"` // inject Service 's ServiceImpl1 implementation
	ServiceImpl2 Service `singleton:"ServiceImpl2"` // inject Service 's ServiceImpl2 implementation
	ServiceStruct *ServiceStruct `singleton:"ServiceStruct"` // inject ServiceStruct struct pointer
}

func (a*App) Run(){
	for {
		time.Sleep(time.Second*3)
		a.ServiceImpl1.Hello()
		a.ServiceImpl2.Hello()
		
		fmt.Println(a.ServiceStruct.GetString("laurence"))
    }
}


type Service interface{
	Hello()
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:interface=Service

type ServiceImpl1 struct {

}

func (s *ServiceImpl1) Hello(){
	fmt.Println("This is ServiceImpl1, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:interface=Service

type ServiceImpl2 struct {

}

func (s *ServiceImpl2) Hello(){
	fmt.Println("This is ServiceImpl2, hello world")
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type ServiceStruct struct {

}

func (s *ServiceStruct) GetString(name string)string{
	return fmt.Sprintf("Hello %s", name)
}

func main(){
	// start
	if err := ioc.Load(); err != nil{
		panic(err)
	}

	// App-App is the format of： '$(interfaceName)-$(implementationStructName)'
	// We can get instance by ths id
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil{
		panic(err)
	}
	app := appInterface.(*App)
	app.Run()
}
```
After writing, you can exec the following cli command to init go mod and generate codes.  (mac may require sudo due to permissions during code generation)

```bash
% go mod init ioc-golang-demo
% go mod tidy
% sudo ioc-go-cli gen
````

It will be generated in the current directory: zz_generated.ioc.go, developers **do not need to care about this file**, this file contains the description information of all interfaces,

```go
//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by ioc-go-cli

package main

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/autowire/singleton"
)

func init() {
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Interface: &App{},
		Factory: func() interface{} {
			return &App{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Interface: new(Service),
		Factory: func() interface{} {
			return &ServiceImpl1{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Interface: new(Service),
		Factory: func() interface{} {
			return &ServiceImpl2{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Interface: &ServiceStruct{},
		Factory: func() interface{} {
			return &ServiceStruct{}
		},
	})
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
````

#### Execute program

**Run with general mode**

`go run .`

Console printout:

```sh
  ___    ___     ____            ____           _                         
 |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                    |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golang config
[Config] Load config file from ../conf/ioc_golang.yaml
Load ioc-golang config file failed. open ../conf/ioc_golang.yaml: no such file or directory
The load procedure is continue
[Boot] Start to load debug
[Debug] Debug mod is not enabled
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD App-App
[Autowire Struct Descriptor] Found type singleton registered SD Service-ServiceImpl1
[Autowire Struct Descriptor] Found type singleton registered SD Service-ServiceImpl2
[Autowire Struct Descriptor] Found type singleton registered SD ServiceStruct-ServiceStruct
This is ServiceImpl1, hello world
This is ServiceImpl2, hello world
Hello laurence
...
```

It shows that the injection is successful and the program runs normally.

**Run with debug mode**

`GOARCH=amd64 go run -gcflags="-N -l" -tags iocdebug .`

Following logs can be found in console output:

```bash
[Debug] Debug server listening at :1999
```

List all interface, implementations and methods

```
% ioc-go-cli list
App
App
[Run]

Service
ServiceImpl1
[Hello]

Service
ServiceImpl2
[Hello]

ServiceStruct
ServiceStruct
[GetString]
```

Watch real-time param and return value. We take 'GetString' method as an example. The method would be called every  3s .

```
% ioc-go-cli watch ServiceStruct ServiceStruct GetString

========== On Call ==========
ServiceStruct.(ServiceStruct).GetString()
Param 1: (string) (len=8) "laurence"


========== On Response ==========
ServiceStruct.(ServiceStruct).GetString()
Response 1: (string) (len=14) "Hello laurence"
...
```



### Annotation Analysis

````go
// +ioc:autowire=true
The code generation tool recognizes objects marked with the +ioc:autowire=true annotation

// +ioc:autowire:type=singleton
The marker autowire model is the singleton, as well as the normal multi-instance model, the config configuration model, the grpc client model and other extensions.

// +ioc:autowire:interface=Service
Markers implement the interface Service and can be injected into objects of type Service .
````

###  More

More code generation annotations can be viewed at [ioc-golang-cli](https://github.com/alibaba/IOC-Golang/tree/master/ioc-go-cli).

You can go to [ioc-golang-example](https://github.com/alibaba/IOC-Golang/tree/master/example) for more examples and advanced usage.

### License

IOC-Golang developed by Alibaba and licensed under the Apache License (Version 2.0).
See the NOTICE file for more information.