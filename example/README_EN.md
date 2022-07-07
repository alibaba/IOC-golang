# IOC-golang Example

## 1. Overview

### AOP

- list_watch

  A debug demo to show feature to list and watch interface, methods and real-time param values.

- transaction

  - singleton

    A demo to show transaction rollback feature with singleton application.

  - distributed

    A demo to show transaction rollback feature in distributed RPC scene.

### Autowire

- autowire_config: 

  Shows how to Inject config field from config file

- autowire_rpc

  An ioc-native RPC demo with client and server. The recommended RPC usage of this framework.

- get_impl_by_api: 

  An example shows how to get object by API.

### Config File

- activate_profile
- complex_example
- default_config_file
- mark_env_variable_in_config_file
- mark_nested_value_in_config_file
- set_config_file_search_path
- set_config_file_type
- set_config_name

### Helloworld

An example to show demo of [README](https://github.com/alibaba/ioc-golang#ioc-golang-a-golang-dependency-injection-framework)

### Third Party

- autowire

  - grpc

    Shows how to inject gRPC client

- db

  - autowire_gorm_db: 

    Shows how to inject GORM client

- registry

  - Nacos

    Shows how to inject Nacos client

- state

  - Redis

    Shows how to inject redis client

## 2. How to run

### 2.1 Run with command line

1. git clone the project

2. cd to demo dir: cd example/helloworld 

3. run with comand line:  `go run .`

   Some example has third-part component dependency, such as autowire_redis_client, we can run command `go test` to start demo with component based on docker, and the detail refers to test files. For examples that have sever dependency, such as autowire_grpc_client, we should start server first.

### 2.2 Run with Goland

1. git clone they project
2. Modify the search path in examples's main method: ` config.WithSearchPath("../conf")`  to ioc_goland.yaml 's parenet dir, or modify to relative path from root dir of the project, such as: `config.WithSearchPath("./extension/third_party/state/redis/conf")`
3. Run or debug with goland.

## 3. More

You can go to [E-commercial system demo based on ioc-golang](https://github.com/ioc-golang/shopping-system) to refer to applications system on distributed scene.

