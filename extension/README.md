# IOC-golang 扩展组件

### [English](./README_EN.md) | 中文

### 基于 ‘自动装载模型’ 和 AOP 能力的扩展

- aop/
  
  提供了基于 AOP 的多种扩展能力，包括与之相关的覆盖 RPC、AOP 拦截器、Debug 服务、cli 命令的相关实现。
  
  - list/
  
    接口方法可视化能力
  
  - watch/
  
    接口、方法、实时参数观测能力
  
  - monitor/

    实时调用监控能力
  
  - trace/
  
    链路追踪、分布式链路追踪能力
  
  - transaction/
  
    事务、分布式事务能力
  
- autowire/
  
  提供了自动装载模型的扩展
  
  - config/
  
    配置字段自动装载模型
  
  - rpc/
  
    ioc-golang 原生支持的 RPC 模型，覆盖 client、server 端。
  
  - grpc/
  
    gRPC 客户端自动装载模型

### 第三方预置组件扩展

为方便开发者直接注入，IOC-golang 提供了覆盖多个领域的预置组件，可供开发者直接注入使用。

- config/

  提供了可以基于配置文件注入的基本数据类型

- config_center/

  提供了可以直接注入的配置中心客户端结构

  - nacos

- db/

  提供了可以操作数据库的客户端结构

  - gorm

- pubsub/

  提供了可以操作消息队列的客户端结构

  - rocketmq

- registry/

  提供了可以操作分布式场景下注册中心的客户端结构

  - nacos

- state/

  提供了可以保存状态、缓存的客户端结构

  - redis

    