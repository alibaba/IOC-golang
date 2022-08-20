# IOC-golang 示例

### [English](./README_EN.md) | 中文

## 1. 简介

### AOP

- obervability

  通过 iocli 工具展示接口、方法、实时参数、调用 RT 等可视化能力。

- transaction

  - singleton

    单体应用的事务回滚例子

  - distributed

    分布式 RPC 场景下的跨进程事务回滚例子。

### Autowire

- autowire_active_profile_implements

  展示了根据配置 profile 激活的情况注入对应实现的例子。

- autowire_allimpls:

  展示了注入一个接口的全部实现的模型。

- autowire_config: 

  展示了从配置文件中注入配置字段的自动装载模型

- autowire_rpc

  展示了 IOC-golang 原生支持的 RPC 能力。

- get_impl_by_api: 

  展示了基于 API 的对象获取方式

### Config File

- active_profile
- complex_example
- default_config_file
- mark_env_variable_in_config_file
- mark_nested_value_in_config_file
- set_config_file_search_path
- set_config_file_type
- set_config_name

### Helloworld

展示了[README](https://github.com/alibaba/ioc-golang#ioc-golang-a-golang-dependency-injection-framework) 中给出的例子

### Third Party

- autowire

  - grpc

    展示了注入 gRPC 客户端的例子

- db

  - autowire_gorm_db: 

    展示了注入 GORM 客户端的例子

- registry

  - Nacos

    展示了注入 Nacos 客户端的例子

- state

  - Redis

    展示了注入 Redis 客户端的例子

## 2. 如何运行

### 2.1 通过命令行启动

1. git clone 本项目

2. 命令行进入示例目录下： cd example/helloworld 

3. 从命令行启动： `go run .`

   对于有依赖组件的例子，例如 autowire_redis_client ，可以通过命令行运行 ` go test `，基于 docker 启动组件。详情参阅测试文件代码。对于有依赖 server 的例子，例如 autowire_grpc_client ，需要先启动 server。

### 2.2 通过 Goland 启动

1. git clone 本项目
2. 修改需要启动的例子 main函数中的` config.WithSearchPath("../conf")` 为ioc_golang.yaml 所在文件夹的绝对路径。或者基于项目根目录的相对路径，例如 `config.WithSearchPath("./extension/third_party/state/redis/conf")`
3. 通过 Goland 启动 main 方法，运行或debug。

## 3. 更多

可以参考 [基于 IOC-golang 的电商系统demo](https://github.com/ioc-golang/shopping-system)  查看分布式场景下的应用系统示例





