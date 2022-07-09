# IOC-golang Extension Components

### English | [中文](./README.md)

### Extensions based on autowire and aop feature

- aop/
  
  Support many extensions feature based on AOP, including related implementations of RPC Interceptor, AOP Interceptor, debug service, cli command.  
  
  - list/
  
    interaface and methods visibility feature
  
  - watch/
  
    realtime param & return values watch feature
  
  - trace/
  
    tracing & distributed system tracing feature
  
  - transaction/
  
    transaction & distributed system transaction feature
  
- autowire/
  
  Support extension of autowire model.
  
  - config/
  
    config field autowire model implementation
  
  - rpc/
  
    ioc-golang native support RPC model, including client and server side implementation.
  
  - grpc/
  
    gRPC client autowire model.

### Third-party components extensions

To support components direct injection by developers, IOC-golang provides pre-built third-party components covering multiple domains that can be directly injected.

- config/

  support basic data types to inject value from config file.

- config_center/

  support config center client implementation.

  - nacos

- db/

  support database client implementation.

  - gorm

- pubsub/

  support message queue client implemetation.

  - rocketmq

- registry/

  support reigstry center client implementation in distributed system.

  - nacos

- state/

  support state storage client implementation.

  - redis

    