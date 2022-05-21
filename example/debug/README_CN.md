# 使用调试功能

### 简介

本示例展示了 ioc-golang 框架提供的代码调试能力

调试能力对于程序性能有损耗，请您不要在追求性能的场景下开启调试能力。

ioc_golang.yaml:

```yaml
debug:
  enable: true # debug 开关，默认为 false
```

debug 模式下，本框架基于 AOP 的思路，为每个注册在框架的结构方法都封装了一组拦截器。基于这些拦截器，可以实现具有很好扩展能力的调试功能。

调试能力包括：

- 基于 ioc-debug 协议，暴露调试端口
- 查看所有接口、实现、方法列表
- 监听、修改任意方法的入参和返回值
- 性能瓶颈分析【开发中】
- 可观测性【开发中】

### 示例介绍

本示例实现了以下拓扑

![debug](https://raw.githubusercontent.com/ioc-golang/ioc-golang-website/main/resources/img/debug-topology.png)

在这个例子中，App 结构会依此调用所有依赖对象，进而调用一个单例模型注入的 gRPC 客户端，该客户端发起网络请求，并获得结果。

我们将开启 debug 模式，通过 ioc-go-cli 工具查看接口、实现、方法，并监听通过 gRPC Client 发送的所有请求和返回值。 

### 运行示例

1. 启动 grpc Server

   ```bash
   % cd example/denbug/grpc_server
   % go run .
   ```

2. 新开一个终端，启动客户端。

   **注意 GOARCH 环境变量和 -gcflags 编译参数, amd机器无需指定 GOARCH 环境变量。**

   正确在 ioc_golang.yaml 中开启debug模式后，会打印 

   `[Debug] Debug port is set to default :1999 `的日志。

   ```bash
   % cd example/debug/cmd
   % GOARCH=amd64 go run -gcflags="-N -l" .
     ___    ___     ____            ____           _                         
    |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
     | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
     | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
    |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                       |___/ 
   Welcome to use ioc-golang!
   [Boot] Start to load ioc-golang config
   [Config] Load config file from ../conf/ioc_golang.yaml
   [Boot] Start to load debug
   [Debug] Debug port is set to default :1999
   [Boot] Start to load autowire
   [Autowire Type] Found registered autowire type singleton
   [Autowire Struct Descriptor] Found type singleton registered SD Service1-Impl1
   [Autowire Struct Descriptor] Found type singleton registered SD Service2-Impl1
   [Autowire Struct Descriptor] Found type singleton registered SD Service2-Impl2
   [Autowire Struct Descriptor] Found type singleton registered SD Struct1-Struct1
   [Autowire Struct Descriptor] Found type singleton registered SD App-App
   [Autowire Type] Found registered autowire type grpc
   [Autowire Struct Descriptor] Found type grpc registered SD HelloServiceClient-HelloServiceClient
   [Debug] Debug server listening at :1999
   create conn target  localhost:8080
   App call grpc get: Hello laurence
   ExampleService1Impl1 call grpc get :Hello laurence_service1_impl1
   ExampleService2Impl1 call grpc get :Hello laurence_service2_impl1
   ExampleService2Impl2 call grpc get :Hello laurence_service2_impl2
   ExampleStruct1 call grpc get :Hello laurence_service1_impl1
   ```

   每隔 5s，所有的对象都会发起一次 gRPC 请求。

3. 新开一个终端，查看所有接口、实现和方法。

   ```bash
   % ioc-go-cli list
   App
   App
   [Run]
   
   HelloServiceClient
   HelloServiceClient
   [SayHello]
   
   Service1
   Impl1
   [Hello]
   
   Service2
   Impl1
   [Hello]
   
   Service2
   Impl2
   [Hello]
   
   Struct1
   Struct1
   [Hello]
   ```

4. 监听 gRPC Client 的所有流量，每隔 5s 会打印出相关的请求、返回值信息。

   ```bash
   % ioc-go-cli watch HelloServiceClient HelloServiceClient SayHello
   ========== On Call ==========
   HelloServiceClient.(HelloServiceClient).SayHello()
   Param 1: (*context.emptyCtx)(0xc0000280e0)(context.Background)
   
   Param 2: (*api.HelloRequest)(0xc000298640)(name:"laurence")
   
   Param 3: ([]grpc.CallOption) (len=2 cap=2) {
    (grpc.MaxRecvMsgSizeCallOption) {
     MaxRecvMsgSize: (int) 1024
    },
    (grpc.MaxRecvMsgSizeCallOption) {
     MaxRecvMsgSize: (int) 1024
    }
   }
   
   
   ========== On Response ==========
   HelloServiceClient.(HelloServiceClient).SayHello()
   Response 1: (*api.HelloResponse)(0xc000298740)(reply:"Hello laurence")
   
   Response 2: (interface {}) <nil>
   
   
   ========== On Call ==========
   HelloServiceClient.(HelloServiceClient).SayHello()
   Param 1: (*context.emptyCtx)(0xc0000280e0)(context.Background)
   
   Param 2: (*api.HelloRequest)(0xc000298900)(name:"laurence_service1_impl1")
   
   Param 3: ([]grpc.CallOption) <nil>
   
   
   ========== On Response ==========
   HelloServiceClient.(HelloServiceClient).SayHello()
   Response 1: (*api.HelloResponse)(0xc0002989c0)(reply:"Hello laurence_service1_impl1")
   
   Response 2: (interface {}) <nil>
   
   ...
   ```

### 小结

通过 Debug 能力，开发人员可以在测试环境内动态地监控流量，帮助排查问题。

也可以基于 ioc-golang 提供的拦截器层，注册任何自己期望的流量拦截器，扩展调试、可观测、运维能力。
