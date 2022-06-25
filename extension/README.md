# IOC-golang 扩展组件

为方便开发者直接注入，IOC-golang 提供了覆盖多个领域的预置组件

- autowire/
  
  提供了自动装载模型的扩展
  
  - config/
  
    配置字段自动装载模型
  
  - rpc/
  
    ioc-golang 原生支持的 RPC 模型，覆盖 client、server 端。
  
  - grpc/
  
    gRPC 客户端自动装载模型
  
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
  
    