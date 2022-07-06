# 使用接口方法、参数的可视化功能

### 简介

本示例展示了 ioc-golang 框架提供的接口方法、参数的可视化功能。

所有在调用中流量经过代理层的结构，都具备调试能力，我们可以通过 iocli 工具，动态监听所有包含代理层的接口、方法、以及实时的参数/返回值内容。

### 专属接口

iocli 会为任何期望注册在框架的结构，生成专属接口，在本例子中，在 zz_generated.ioc.go 为 ServiceImpl1 生成了ServiceImpl1IOCInterface 接口，该接口会包含 ServiceImpl1 结构的全部方法。为 ServiceImpl2 结构也生成了 ServiceImpl2IOCInterface 专属接口。

专属接口的命名规则为 $(结构名)IOCInterface

### 为注入结构封装代理层

任何被注入到接口的字段，都会被框架自动封装代理 AOP 层，即注入到接口的结构体指针，并非真实结构体指针，而是封装了结构体的代理指针。例如：

```go
// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceImpl2 Service `singleton:"main.ServiceImpl2"`
  
	Service1OwnInterface ServiceImpl1IOCInterface `singleton:""`
}
```

例子中的 App.ServiceImpl2 字段，标签中指定的注入结构是 main.ServiceImpl2 是期望将 main.ServiceImpl2 结构体注入至 Service 接口，这个过程被框架注入的接口即包含代理层。

例子中的 ServiceImpl1IOCInterface 字段，是期望注入 ServiceImpl1 结构至它的专属接口，专属接口的注入就不需要在标签中指定结构体ID了。只需要填写空 `singleton:""` 即可。

### 通过 API 获取代理接口

就像通过 API 的方式获取结构体指针一样，也可以通过 API 的形式获得封装了代理层的接口。如例子中的：

```go
// app, err := GetAppSingleton() 获取真实结构体指针
app, err := GetAppIOCInterfaceSingleton()
if err != nil {
  panic(err)
}
```

我们可以调用 iocli 为结构生成的 Get 方法：GetAppIOCInterface，来获取封装了代理层的对象。

### 运行例子：

```bash
% sudo iocli gen
% go run .
  ___    ___     ____                           _                         
 |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                |___/                               |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golang config
[Config] Load default config file from ../conf/ioc_golang.yaml
[Config] Load ioc-golang config file failed. open /Users/laurence/Desktop/workplace/alibaba/IOC-Golang/example/conf/ioc_golang.yaml: no such file or directory
 The load procedure is continue
[Boot] Start to load debug
[Debug] Debug port is set to default :1999
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type normal
[Autowire Struct Descriptor] Found type normal registered SD main.app_
[Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl1_
[Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl2_
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl1
[Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl2
[Debug] Debug server listening at :1999
This is ServiceImpl1, hello laurence
This is ServiceImpl2, hello laurence
This is ServiceImpl1, hello laurence
This is ServiceImpl2, hello laurence
...
```

可看到每隔三秒钟，ServiceImpl1 和 ServiceImpl2 的方法就会被调用。下面我们新启动一个终端，使用 iocli 工具调试这个程序: 

查看所有接口和方法：

```bash
% iocli list
main.App
[Run]

main.ServiceImpl1
[GetHelloString]

main.ServiceImpl2
[GetHelloString]

```

监听接口参数和返回值：

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

可看到每隔三秒钟，就会监听到方法调用的参数和返回值。
