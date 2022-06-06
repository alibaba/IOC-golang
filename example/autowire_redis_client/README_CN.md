# Redis 客户端注入示例

### 简介

本示例展示了注入 Redis 客户端能力。

在应用开发过程中，通过 SDK 操作 Redis 是一个常见的诉求。

ioc-golang 框架提供了注入 Redis 连接的能力，开发者可以在配置文件中指定好 Redis 地址、密码、db名等信息，通过标签注入连接，无需手动创建、组装。

### 注入模型与结构

[多例（normal）依赖注入模型](https://github.com/alibaba/IOC-Golang/tree/master/extension/normal)

[预定义的 redis 结构](https://github.com/alibaba/IOC-Golang/tree/master/extension/normal/redis)

### 关键代码

```go
import(
	normalMysql "github.com/alibaba/ioc-golang/extension/normal/mysql"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	NormalRedis    normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl"`
	NormalDB1Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,db1-redis"`
	NormalDB2Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,db2-redis"`
	NormalDB3Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,address=127.0.0.1:6379&db=3"`
}
```

- 被注入字段类型

  normalRedis.Redis 接口

- 标签

  开发人员可以为 normalRedis.Redis 类型的字段增加 `normal:"Impl,$(configKey),$(tableName)" `标签。从而注入Redis  sdk。

  例子中的 `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl"` 的意义为，将配置文件内 `autowire.normal.Redis.Impl.param`定义的值作为参数。

  例子中的 `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,db1-redis"` 的意义为，将配置文件内 `autowire.normal.Redis.Impl.db1-redis.param`定义的值作为参数。
  
  例子中的 `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,db2-redis"`的意义为，将配置文件内 `autowire.normal.Redis.Impl.db2-redis.param`定义的值作为参数。
  
  例子中的 `normal:"github.com/alibaba/ioc-golang/extension/normal/redis.Impl,address=127.0.0.1:6379&db=3"` 的意义为，使用标签内定义的 key-value 作为参数配置。
  
  默认参数加载策略详情请参阅[参数加载器](/docs/concept/param_loader/)
  
  ```yaml
  autowire:
    normal:
      github.com/alibaba/ioc-golang/extension/normal/redis.Impl:
        db1-redis:
          param:
            address: localhost:6379
            db: 1
        db2-redis:
          param:
            address: localhost:6379
            db: 2
        param:
          address: localhost:6379
          db: 0
  ```
  

### 运行示例

例子会注入多个位于127.0.0.1:6379 的Redis 客户端，数据库id分别为 0、1、2、3. 

需要确保您本地运行了redis，可通过 `docker run -p6379:6379 redis:latest` 快速启动一个。

注入后，方可调用该接口提供的方法。可获取裸 Redis 连接，也可以直接使用封装好的 API 操作Redis。

```bash
cd example/autowire_redis_client/cmd
go run .
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
[Debug] Debug mod is not enabled
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Autowire Type] Found registered autowire type normal
[Autowire Struct Descriptor] Found type normal registered SD github.com/alibaba/ioc-golang/extension/normal/redis.Impl
client0 get  db0
client1 get  db1
client2 get  db2
client3 get  db3
```

可看到打印出了已经写入四个 redis db 的值。





