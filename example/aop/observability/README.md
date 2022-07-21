# 可视化能力

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
	Service1OwnInterface ServiceImpl1IOCInterface `singleton:""`
}
```

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
This is ServiceImpl2, hello laurence
This is ServiceImpl2, hello laurence
...
```

可看到每隔三秒钟，ServiceImpl1 和 ServiceImpl2 的方法就会被调用。下面我们新启动一个终端，使用 iocli 工具调试这个程序: 

- 查看所有接口和方法：

```bash
% iocli list
main.App
[Run]

main.ServiceImpl1
[GetHelloString]

main.ServiceImpl2
[GetHelloString]

```

- 监听接口参数和返回值：

```bash
% iocli watch main.ServiceImpl1 GetHelloString
iocli watch started, try to connect to debug server at 127.0.0.1:1999
debug server connected, watch info would be printed when invocation occurs, param info max depth = 5
========== On Call ==========
main.ServiceImpl1.GetHelloString()
Param 1: (string) (len=8) "laurence"

========== On Response ==========
main.ServiceImpl1.GetHelloString()
Response 1: (string) (len=36) "This is ServiceImpl2, hello laurence"

========== On Call ==========
main.ServiceImpl1.GetHelloString()
Param 1: (string) (len=8) "laurence"

========== On Response ==========
main.ServiceImpl1.GetHelloString()
Response 1: (string) (len=36) "This is ServiceImpl2, hello laurence"

...
```

可看到每隔三秒钟，就会监听到方法调用的参数和返回值。

- 监控应用

```
% iocli monitor
iocli monitor started, try to connect to debug server at 127.0.0.1:1999
debug server connected, monitor info would be printed every 5s
====================
2022/07/10 19:39:26
main.ServiceImpl1.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: 39.00us, FailRate: 0.00%
main.ServiceImpl2.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: 22.00us, FailRate: 0.00%
====================
2022/07/10 19:39:31
main.ServiceImpl1.GetHelloString()
Total: 2, Success: 2, Fail: 0, AvgRT: 57.00us, FailRate: 0.00%
main.ServiceImpl2.GetHelloString()
Total: 2, Success: 2, Fail: 0, AvgRT: 27.50us, FailRate: 0.00%

...
^C
Got interrupt signal, collecting data during 15663ms
====================Collection====================
2022/07/10 19:39:36
main.ServiceImpl1.GetHelloString()
Total: 5, Success: 5, Fail: 0, AvgRT: 46.17us, FailRate: 0.00%
main.ServiceImpl2.GetHelloString()
Total: 5, Success: 5, Fail: 0, AvgRT: 20.50us, FailRate: 0.00%

```

可看到在一段时间内，所有接口方法的调用情况。默认每隔五秒钟刷新一次这五秒内的调用情况，Control+C 终止进程时，会打印这段时间内的全部调用信息统计。包括请求次数、RT、失败率等信息。

- 调用链路追踪

```go
% iocli trace main.ServiceImpl1 GetHelloString
iocli trace started, try to connect to debug server at 127.0.0.1:1999
debug server connected, tracing info would be printed when invocation occurs
==================== Trace ==================== 
Duration 9us, OperationName: main.(*serviceImpl2_).GetHelloString, StartTime: 2022/07/10 11:41:32, ReferenceSpans: [{TraceID:75698db3dfec8990 SpanID:01dc1b3c44bf9bb8 RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
Duration 957us, OperationName: main.(*serviceImpl1_).GetHelloString, StartTime: 2022/07/10 11:41:32, ReferenceSpans: [{TraceID:75698db3dfec8990 SpanID:75698db3dfec8990 RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
==================== Trace ==================== 
Duration 10us, OperationName: main.(*serviceImpl2_).GetHelloString, StartTime: 2022/07/10 11:41:35, ReferenceSpans: [{TraceID:301fabdc183b603d SpanID:5bfc284b130c0d18 RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
Duration 56us, OperationName: main.(*serviceImpl1_).GetHelloString, StartTime: 2022/07/10 11:41:35, ReferenceSpans: [{TraceID:301fabdc183b603d SpanID:301fabdc183b603d RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
...
```

可看到每一次调用到监控方法对应一个Trace，其全部调用 Span 都被打印了出来，包含 RT、时间等信息

- 在本地可视化调用链路

需要在本地启动 jaeger-collector、jaeger-query、elsaticsearch, 可参考 shopping-system 例子的 docker 启动方式 [docker-compose](https://github.com/ioc-golang/shopping-system/blob/main/deploy/docker-compose/docker-compose.yaml)

```go
 iocli trace  main.ServiceImpl1 GetHelloString --pushAddr localhost:14268
iocli trace started, try to connect to debug server at 127.0.0.1:1999
debug server connected, tracing info would be printed when invocation occurs
try to push span batch data to localhost:14268
==================== Trace ====================
Duration 8us, OperationName: main.(*serviceImpl2_).GetHelloString, StartTime: 2022/07/10 11:47:23, ReferenceSpans: [{TraceID:79425e930368c5dd SpanID:6c75b07c1f19a82b RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
Duration 75us, OperationName: main.(*serviceImpl1_).GetHelloString, StartTime: 2022/07/10 11:47:23, ReferenceSpans: [{TraceID:79425e930368c5dd SpanID:79425e930368c5dd RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
==================== Trace ====================
Duration 1312us, OperationName: main.(*serviceImpl2_).GetHelloString, StartTime: 2022/07/10 11:47:26, ReferenceSpans: [{TraceID:6a36fcdd7909db2c SpanID:5abedfc7a9c2f3ca RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================
Duration 1341us, OperationName: main.(*serviceImpl1_).GetHelloString, StartTime: 2022/07/10 11:47:26, ReferenceSpans: [{TraceID:6a36fcdd7909db2c SpanID:6a36fcdd7909db2c RefType:CHILD_OF XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
====================

```

`--pushAddr` 参数指定需要推送至本机的 jaeger-collector 地址，该地址只需 iocli 进程可达，调用链路数据将由 iocli 工具拉取到本机，并推送至 jaeger-collector，不要求应用进程的部署环境有可视化组件。

浏览器访问 localhost:16686 查看调用链路信息。

![img.png](https://raw.githubusercontent.com/ioc-golang/ioc-golang-website/main/resources/img/example-aop-observability-tracing.png)

可通过`iocli -h` 查看更多命令和参数

