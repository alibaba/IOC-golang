# 数据库连接注入示例

### 简介

本示例展示了注入基于 [GORM](https://gorm.io/index.html) 的数据库客户端能力。

在应用开发过程中，通过 SDK 操作数据库是一个常见的诉求，GORM 是应用较广泛的 Go 数据库 sdk。

ioc-golang 框架提供了注入数据库连接的能力，开发者可以在配置文件中指定好数据库地址、密码信息，通过标签注入连接，无需手动创建、组装 GORM 客户端。

### 注入模型与结构

[多例（normal）依赖注入模型](https://github.com/alibaba/IOC-Golang/tree/master/extension/normal)

[预定义的 mysql 结构](https://github.com/alibaba/IOC-Golang/tree/master/extension/normal/mysql)

### 关键代码

```go
import(
	normalMysql "github.com/alibaba/ioc-golang/extension/normal/mysql"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	MyDataTable normalMysql.Mysql `normal:"github.com/alibaba/ioc-golang/extension/normal.Impl,my-mysql,mydata"`
}
```

- 被注入字段类型

  normalMysql.Mysql 接口

- 标签

  开发人员可以为 normalMysql.Mysql 类型的字段增加 `normal:"Impl,$(configKey),$(tableName)" `标签。从而注入指定数据库的指定表 sdk。

  例子中的 `normal:"github.com/alibaba/ioc-golang/extension/normal/mysql.Impl,my-mysql,mydata"` 的意义为，将配置文件内`autowire.normal.<github.com/alibaba/ioc-golang/extension/normal/mysql.Impl>.my-mysql.param`定义的值作为参数。

  ```yaml
  autowire:
    normal:
        github.com/alibaba/ioc-golang/extension/normal/mysql.Impl:
          my-mysql:
            param:
              host: "127.0.0.1"
              port: 3306
              username: "root"
              password: "root"
              dbname: "test"
  ```

  例子会建立一个位于127.0.0.1:3306 的数据库连接，用户名为root、密码为 root、数据库名为test、表名为mydata。

  注入后，方可调用该接口提供的方法，获取裸 gorm 连接，或者直接使用封装好的 API 操作数据表。

  ```go
  type Mysql interface {
      GetDB() *gorm.DB
      SelectWhere(queryStr string, result interface{}, args ...interface{}) error
      Insert(toInsertLines UserDefinedModel) error
      Delete(toDeleteTarget UserDefinedModel) error
      First(queryStr string, findTarget UserDefinedModel, args ...interface{}) error
      Update(queryStr, field string, target interface{}, args ...interface{}) error
  }
  ```

  

