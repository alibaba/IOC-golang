# Mysql

## 基本信息

- Autowire 类型：normal

- SDID：Mysql-Impl

- 参数来源：定制化参数加载方式

  ​	标签 + 配置文件。需要同时设定好注入标签与配置文件。

## 参数说明：

```go
type Config struct {
	Host         string `yaml:"host"` // 数据库地址：从配置文件读入
	Port         string `yaml:"port"` // 端口号：从配置文件读入
	Username     string `yaml:"username"` // 用户名：从配置文件读入
	Password     string `yaml:"password"` // 密码：从配置文件读入
	DBName       string `yaml:"dbname"` // db 名：从配置文件读入
	TableName    string // 表名：从标签读入
}
```

## 注入示例

### 通过依赖注入

标签

```go
// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	MyDataTable normalMysql.Mysql `normal:"Impl,my-mysql,mydata"` // Impl,配置key，表明
}
```

配置

```yaml
autowire:
  normal:
    github.com/alibaba/ioc-golang/extension/normal/mysql.Impl:
        param:
          my-mysql:
            host: "127.0.0.1"
            port: 3306
            username: "root"
            password: "root"
            dbname: "test"
```

标签中 Impl 为固定值。mysql 为

可获取到参数：

```go
type Config struct {
	Host         string // 127.0.0.1
	Port         string // 3306
	Username     string // root
	Password     string // root
	DBName       string // test
	TableName    string // mydata
}
```

### 通过 API 获取

```go
import (
  normalMysql "github.com/alibaba/ioc-golang/extension/normal/mysql"
)


mysqlImpl, err := normalMysql.GetMysql(&normalMysql.Config{
		Host: "127.0.0.1",
		Port: "3306",
		...
})
if err != nil{
  panic(err)
}
```

## 方法说明

```go
type Mysql interface {
	GetDB() *gorm.DB // 获取 gorm 数据库连接
  
  // 包装 API 
	SelectWhere(queryStr string, result interface{}, args ...interface{}) error
	Insert(toInsertLines UserDefinedModel) error
	Delete(toDeleteTarget UserDefinedModel) error
	First(queryStr string, findTarget UserDefinedModel, args ...interface{}) error
	Update(queryStr, field string, target interface{}, args ...interface{}) error
}
```



