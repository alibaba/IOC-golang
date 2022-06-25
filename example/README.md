# IOC-golang 示例

## 1. 简介

- autowire_config: 

  展示了如何从配置文件中注入值到结构体。

- autowire_gorm_db: 

  展示了注入 GORM 客户端的例子

- autowire_grpc_client

  展示了注入 gRPC 客户端的例子

- autowire_nacos_client

  展示了注入 Nacos 客户端的例子

- autowire_redis_client

  展示了注入 Redis 客户端的例子

- autowire_rpc

  展示了 IOC-golang 原生支持的 RPC 能力。

- debug

  展示了基于接口代理层的可观测能力。

- debug_with_monkey

  展示了基于 monkey 指针的，为结构体指针封装代理层的可观测能力。

- get_impl_by_api

  展示了基于 API 的对象获取方式

- helloworld

  展示了[README](https://github.com/alibaba/ioc-golang#ioc-golang-a-golang-dependency-injection-framework) 中给出的例子

## 2. 如何运行

### 2.1 通过命令行启动

1. git clone 本项目

2. 命令行进入示例目录下： cd example/helloworld 

3. 从命令行启动： `go run .`

   对于有依赖组件的例子，例如 autowire_redis_client ，可以通过命令行运行 ` go test `，基于 docker 启动组件。详情参阅测试文件代码。对于有依赖 server 的例子，例如 autowire_grpc_client ，需要先启动 server。

### 2.2 通过 Goland 启动

1. git clone 本项目
2. 修改需要启动的例子 main函数中的` config.WithSearchPath("../conf")` 为ioc_golang.yaml 所在文件夹的绝对路径。或者基于项目根目录的相对路径，例如 `config.WithSearchPath("./extension/autowire_redis_client/conf")`
3. 通过 Goland 启动 main 方法，运行或debug。

## 3. 更多

可以参考 [基于 IOC-golang 的电商系统demo](https://github.com/ioc-golang/shopping-system)  查看分布式场景下的应用系统示例





