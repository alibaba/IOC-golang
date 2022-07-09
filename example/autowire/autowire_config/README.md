# 配置注入示例

### 简介

本示例展示了从配置文件注入字段的能力。

在应用开发过程中，从配置文件中读入配置是常见的诉求。例如读取数据库的账号密码、下游服务的主机名，以及一些业务配置等。

ioc-golang 框架提供了便捷的基于文件注入配置的能力，使开发者无需手动解析配置文件，无需手动组装对象。

### 依赖注入模型

[config 依赖注入模型](https://github.com/alibaba/IOC-golang/tree/master/extension/config)

### 关键代码：

```go
import (
  github.com/alibaba/ioc-golang/extension/config
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	DemoConfigString *config.ConfigString `config:",autowire.config.demo-config.string-value"`
	DemoConfigInt    *config.ConfigInt    `config:",autowire.config.demo-config.int-value"`
	DemoConfigMap    *config.ConfigMap    `config:",autowire.config.demo-config.map-value"`
	DemoConfigSlice  *config.ConfigSlice  `config:",autowire.config.demo-config.slice-value"`
}
```

- 被注入字段类型

  目前支持 ConfigString，ConfigInt，ConfigMap，ConfigSlice 四种类型。

  需要以 **指针** 的形式声明字段类型

- 标签与注入位置

  开发人员可以给结构增加 ``config:`"xxx" ` 标签, 标注需要注入的值类型，以及该字段位于配置文件的位置。

  例子中的

  `config:",autowire.config.demo-config.string-value"`

  的意义为，将配置文件内 `autowire.config.demo-config.string-value` 的值注入到该字段。

  对应配置文件：ioc_golang.yaml 中的字符串 "stringValue"

```yaml
autowire:
  config:
    demo-config:
      int-value: 123
      int64-value: 130117537261158665
      float64-value: 0.001
      string-value: stringValue
      map-value:
        key1: value1
        key2: value2
        key3: value3
        obj:
          objkey1: objvalue1
          objkey2: objvalue2
          objkeyslice: objslicevalue
      slice-value:
        - sliceValue1
        - sliceValue2
        - sliceValue3
        - sliceValue4
 ```

### 运行示例

```bash
cd example/autowire_config/cmd
go run .
  ___    ___     ____            ____           _                         
 |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                    |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golang config
[Config] merge config map, depth: [0]
[Boot] Start to load debug
[Debug] Debug mod is not enabled
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Autowire Type] Found registered autowire type normal
[Autowire Type] Found registered autowire type config
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigInt64
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigInt
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigMap
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigSlice
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigString
[Autowire Struct Descriptor] Found type config registered SD github.com/alibaba/ioc-golang/extension/config.ConfigFloat64
2022/06/06 18:01:22 load config path autowire.config#demo-config#float64-value error =  property [autowire config#demo-config#float64-value]'s key config#demo-config#float64-value not found
stringValue
123
map[key1:value1 key2:value2 key3:value3 obj:map[objkey1:objvalue1 objkey2:objvalue2 objkeyslice:objslicevalue]]
[sliceValue1 sliceValue2 sliceValue3 sliceValue4]
130117537261158665
0
stringValue
123
map[key1:value1 key2:value2 key3:value3 obj:map[objkey1:objvalue1 objkey2:objvalue2 objkeyslice:objslicevalue]]
[sliceValue1 sliceValue2 sliceValue3 sliceValue4]
130117537261158665
0

```

可看到依次打印出了不同结构的注入配置。

