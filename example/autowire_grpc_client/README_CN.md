# 注入 gRPC 客户端示例

### 简介

本示例展示了基于 ioc-golang 框架的 gRPC 客户端注入能力。

在进行微服务开发过程中，服务间通信尤为重要，gRPC 是被应用最为广泛的 RPC 框架之一。

在常规开发中，开发者需要从配置中手动读取下游主机名，启动 grpc 客户端。针对一个接口的网络客户端往往是单例模型，如果多个服务都需要使用同一客户端，则还需要开发者维护这个单例模型。

基于 ioc-golang 框架的 gRPC 客户端注入能力，我们可以将客户端的生命周期交付给框架管理，并赋予客户端调试能力，开发者只需要关注注册和使用。

### 示例介绍

本示例实现了以下拓扑

![debug](https://raw.githubusercontent.com/ioc-golang/ioc-golang-website/main/resources/img/debug-topology.png)

在这个例子中，App 结构会依此调用所有依赖对象，进而调用一个单例模型注入的 gRPC 客户端，该客户端发起网络请求，并获得结果。

### 依赖注入模型

[grpc 依赖注入模型](https://github.com/alibaba/IOC-Golang/tree/master/extension/grpc)

### 关键代码

```go
import(
	"github.com/alibaba/ioc-golang/extension/autowire/grpc"
	googleGRPC "google.golang.org/grpc"
)
func init() {
	// register grpc client
	grpc.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return new(api.HelloServiceClient)
		},
		ParamFactory: func() interface{} {
			return &googleGRPC.ClientConn{}
		},
		ConstructFunc: func(impl interface{}, param interface{}) (interface{}, error) {
			conn := param.(*googleGRPC.ClientConn)
			fmt.Println("create conn target ", conn.Target())
			return api.NewHelloServiceClient(conn), nil
		},
	})
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	HelloServiceClient api.HelloServiceClient `grpc:"hello-service"`
}
```

需要在代码中手动注册 gRPC 客户端。在需要使用的地方，增加  `grpc:"xxx"` 标签

框架会默认从 autowire.grpc.xxx 读取参数, 在例子中，为`autowire#grpc#hello-service `

```yaml
autowire:
  grpc:
    hello-service:
      address: localhost:8080
```

### 运行示例

1. 启动 grpc Server

   ```bash
   % cd example/denbug/grpc_server
   % go run .
   ```

2. 新开一个终端，启动客户端。

```bash
% cd example/autowire_grpc_client/cmd
% go run .
  ___    ___     ____            ____           _                         
 |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                    |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golang config
[Config] Load default config file from ../conf/ioc_golang.yaml
[Config] merge config map, depth: [0]
[Boot] Start to load debug
[Debug] Debug mod is not enabled
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type grpc
[Autowire Struct Descriptor] Found type grpc registered SD github.com/alibaba/ioc-golang/example/autowire_grpc_client/api.HelloServiceClient
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/autowire_grpc_client/cmd/struct1.Struct1
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/autowire_grpc_client/cmd/service1.Impl1
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/autowire_grpc_client/cmd/service2.Impl1
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/autowire_grpc_client/cmd/service2.Impl2
create conn target  localhost:8080
App call grpc get: Hello laurence
ExampleService1Impl1 call grpc get :Hello laurence_service1_impl1
ExampleService2Impl1 call grpc get :Hello laurence_service2_impl1
ExampleService2Impl2 call grpc get :Hello laurence_service2_impl2
ExampleStruct1 call grpc get :Hello laurence_struct
```
   

### 小节

gRPC 客户端注入能力，所代表的是 ioc-golang 框架具备网络模型注入的能力。

针对特定网络框架，可以在其服务提供者接口处，提供客户端注入代码，从而便于客户端引入后直接注入。

