# 通过 API 获取对象

### 简介

本示例展示了通过 API 调用的方式获取对象的能力。

在应用开发过程中，部分依赖通过注入的方式预置在字段中，也有一部分依赖是在程序运行过程中动态生成的。常规做法是通过手动拼装结构的方式，或者通过调用构造函数的方式获取对象。

ioc-golang 框架提供了中心化的获取对象 API。

- 该 API 推荐被结构提供者封装，从而提供一个具象的 API 供用户调用。该 API 传入所需配置结构，返回具体的接口。

extension/normal/redis/redis.go

```go
func GetRedis(config *Config) (Redis, error) {
	mysqlImpl, err := normal.GetImpl(SDID, config)
	if err != nil {
		return nil, err
	}
	return mysqlImpl.(Redis), nil
}
```

- 如果结构提供者并没有提供上述 API，用户同样也可以直接调用，传入参数并获取对象。

### 对象获取 API

- 多例（normal）

  autowire/normal/normal.go

  ```go
  func GetImpl(sdID string, param interface{}) (interface{}, error) {}
  ```

  每次调用多例获取 API，将创建一个新对象。

- 单例（singleton）

  autowire/singleton/singleton.go

  ```go
  func GetImpl(sdID string) (interface{}, error) {}
  ```

  单例模型全局只拥有一个对象，通过 API 只能获取而不能创建，在框架启动时所有单例模型指针将基于配置/标签创建好。

### 关键代码

```go
import(
	"github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/singleton"
	"github.com/alibaba/ioc-golang/extension/normal/redis"
)


func (a *App) Run() {
  // 通过 normal 提供的全局 API，传递结构描述 ID 和配置结构，创建多例对象
	normalRedis, err := normal.GetImpl("Redis-Impl", &redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
  // 通过 redis 结构提供者定义好的 GetRedis 方法，传递配置结构，创建多例对象
	normalRedis2, err := redis.GetRedis(&redis.Config{
		Address: "localhost:6379",
		DB:      "0",
	})
  ...
}

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
  // 通过 singleton 提供的全局 API，传递结构描述 ID 获取单例对象
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil {
		panic(err)
	}
	app := appInterface.(*App)

	app.Run()
}

```

### 运行示例

需要确保您本地运行了redis，可通过 `docker run -p6379:6379 redis:latest` 快速启动一个。

例子会通过 API 获取 App 对象和 redis 对象，并调用 redis 对象提供的方法。

```bash
 cd example/get_impl_by_api/cmd
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
Load ioc-golang config file failed. open ../conf/ioc_golang.yaml: no such file or directory
The load procedure is continue
[Boot] Start to load debug
[Debug] Debug mod is not enabled
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type normal
[Autowire Struct Descriptor] Found type normal registered SD Redis-Impl
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD App-App
get val =  db0
```

可看到打印出了 redis 中的数据。





