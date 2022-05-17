# Redis

## 基本信息

- Autowire 类型：normal
- SDID：Redis-Impl
- 参数来源：默认参数加载方式

## 参数说明：

```go
type Config struct {
	Address  string // redis 地址
	Password string // 密码
	DB       string // 数据库序号
}
```

## 通过依赖注入

略

## 通过 API 获取

```go
import (
  normalRedis "github.com/alibaba/IOC-Golang/extension/normal/redis"
)

redisImpl, err := normalRedis.GetRedis(&normalRedis.Config{
  Address: "localhost:6379",
})
if err != nil{
  panic(err)
}

```


## 方法说明

```go
type Redis interface {
	GetRawClient() *redis.Client // 获取 redis 数据库连接
  
	Set(key string, value interface{}, expiration time.Duration) (string, error)
	Get(key string) (string, error)
	HGetAll(key string) (map[string]string, error)
}
```
