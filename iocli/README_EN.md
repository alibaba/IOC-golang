# iocli tool

## [中文](./README.md) | English

**iocli** is a command line tool, it has the following features: 

- program debug

  Developers can use **iocli** as debug client and debug the go program developed with ioc-golang.

- codes generation

  Developers can add annotations to structures that describe the struct to be   injected, and **iocli** will identify these annotations and generates structure-specific code that meets the requirements. Including structure description information, structure proxy layer, structure own interface, structure Get method and so on.

## Program Debug Feature

ioc-golang 框架拥有基于结构代理层的 Go 运行时程序调试能力，帮助故障排查，性能分析，提高应用可观测能力。在 [README](https://github.com/alibaba/ioc-golang#ioc-golang-a-golang-dependency-injection-framework)  Quickstart 中展示了接口信息的查看、参数监听能力。在 [基于 IOC-golang 的电商系统demo](https://github.com/ioc-golang/shopping-system)  中，可以展示基于 ioc-golang 的，业务无侵入的，方法粒度全链路追踪能力。

## Annotation and Code generation

注解是以特定字符串开头的注释，标注在期望注入的结构前。注解只具备静态意义，即在代码生成阶段，被iocli工具扫描识别到，从而获取结构相关信息。注解本身不具备程序运行时的意义。

iocli 可以识别以下注解：

```go
// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramLoader=paramLoader
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New
// +ioc:autowire:baseType=true
// +ioc:autowire:alias=MyAppAlias
```

- ioc:autowire 

  bool 类型，为 true 则在代码生成阶段被识别到。

- ioc:autowire:type

  string类型，表示依赖注入模型，目前支持以下五种，结构提供者可以选择五种中的一种或多种进行标注，从而生成相关的结构信息与 API，供结构使用者选用。

  ```go
  // +ioc:autowire:type=singleton
  // +ioc:autowire:type=normal
  
  type MyStruct struct{
  
  }
  ```

  - singleton

    单例模型，使用该注入模型获取到的结构体，全局只存在一个对象。

  - normal

    多例模型，使用该注入模型，每一个标签注入字段、每一次 API 获取，都会产生一个新的对象。

  - config:

    配置模型是基于多例模型的封装扩展，基于配置模型定义的结构体方便从 yaml 配置文件中注入信息。参考例子 [example/autowire_config](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_config)

  - grpc:

    grpc 模型是基于单例模型的封装扩展，基于 grpc 模型可以方便地从 yaml 配置文件中读取参数，生成 grpc 客户端。参考例子[example/autowire_grpc_client](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_grpc_client)

  - rpc:

    rpc 模型会在代码生成阶段产生 rpc 服务端注册代码，以及 rpc 客户端调用存根。参考例子 [example/autowire_rpc](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_rpc)

  

- ioc:autowire:paramLoader（非必填）

  string类型，表示需要定制的“参数加载器“类型名

  参数加载器由结构定义者可选定制。可参考：[ioc-go-extension/normal/redis](http://github.com/alibaba/ioc-golang/extension/blob/master/normal)

  参数加载器需要实现Load方法：

  ```go
  // ParamLoader is interface to load param
  type ParamLoader interface {
  	Load(sd *StructDescriptor, fi *FieldInfo) (interface{}, error)
  }
  ```

  定义结构的开发者可以通过实现参数加载器，来定义自己的结构初始化参数。例如，一个 redis 客户端结构 'Impl'，需要从Config 参数来加载，如下所示 New 方法。

  ```go
  type Config struct {
  	Address  string
  	Password string
  	DB       string
  }
  
  func (c *Config) New(impl *Impl) (*Impl, error) {
  	dbInt, err := strconv.Atoi(c.DB)
  	if err != nil {
  		return impl, err
  	}
  	client := redis.NewClient(&redis.Options{
  		Addr:     c.Address,
  		Password: c.Password,
  		DB:       dbInt,
  	})
  	_, err = client.Ping().Result()
  	if err != nil {
  		return impl, err
  	}
  	impl.client = client
  	return impl, nil
  }
  ```

  Config 包含的三个字段：Address Password DB，需要由使用者传入。

  从哪里传入？这就是参数加载器所做的事情。

  结构定义者可以定义如下加载器，从而将字段通过注入该结构的 tag 标签获取，如果tag信息标注了配置位置，则通过配置文件获取。

  ```go
  type paramLoader struct {
  }
  
  func (p *paramLoader) Load(sd *autowire.StructDescriptor, fi *autowire.FieldInfo) (interface{}, error) {
  	splitedTagValue := strings.Split(fi.TagValue, ",")
  	param := &Config{}
  	if len(splitedTagValue) == 1 {
  		return nil, fmt.Errorf("file info %s doesn't contain param infomration, create param from sd paramLoader failed", fi)
  	}
  	if err := config.LoadConfigByPrefix("extension.normal.redis."+splitedTagValue[1], param); err != nil {
  		return nil, err
  	}
  	return param, nil
  }
  ```

  例如 

  ```go
  type App struct {
  	NormalDB3Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/state/redis.Redis,address=127.0.0.1:6379&db=3"`
  }
  ```

  当然也可以从配置文件读入，tag中指定了key为 db1-redis

  ```go
  type App struct {
  	NormalDB3Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/state/redis.Redis,db1-redis"`
  }
  ```

  ioc-go.yaml： autowire.normal.Redis.Impl.db1-redis.param 读入参数

  ```yaml
  autowire:
    normal:
      github.com/alibaba/ioc-golang/extension/state/redis.Redis:
        db1-redis:
          param:
            address: localhost:6379
            db: 1
  ```
  
  
  
  **我们提供了预置的参数加载器**
  
  除非用户有强烈需求，我们更推荐用户直接使用我们预置的参数加载器：http://github.com/alibaba/ioc-golang/tree/master/autowire/param_loader。
  
  我们会先后尝试：标签重定向到配置、标签读入参数、配置文件的默认位置读入参数。每个注册到框架的结构都有唯一的ID，因此也会在配置文件中拥有配置参数的位置，这一默认位置在这里定义：http://github.com/alibaba/ioc-golang/blob/master/autowire/param_loader/default_config.go#L21，我们更希望和用户约定好这一点。
  
  当所有加载器都加载参数失败后，将会抛出错误。使用者应当查阅自己引入的结构加载器实现，并按照要求配置好参数。
  
- ioc:autowire:paramType（非必填）

  string类型，表示依赖参数的类型名，在上述例子，该类型名为 Config

- ioc:autowire:constructFunc（非必填）

  string类型，表示结构的构造方法名

  在给出 ioc:autowire:paramType 参数类型名的情况下，会使用参数的函数作为构造函数，例如在上述例子中，该构造方法为 Config 对象的 New 方法。

  如果没有给出 ioc:autowire:paramType 参数类型名，则会直接使用这一方法作为构造函数。

  我们要求该构造方法的函数签名是固定的，即：

  ```go
  func (*$(结构名)) (*$(结构名）, error)
  ```

- ioc:autowire:baseType=true （非必填）

  该类型是否为基础类型

  go 基础类型不可直接通过&在构造时取地址，因此我们针对基础类型单独设计了该注解。在 http://github.com/alibaba/ioc-golang/extension/tree/master/config 配置扩展中被使用较多。

- ioc:autowire:alias=MyAppAlias （非必填）

  该类型的别名，可在标签、API获取、配置中，通过该别名替代掉较长的类型全名来指定结构。

## iocli 操作命令文档

- `iocli init`

  生成初始化工程

- `iocli gen`

  递归遍历当前目录下的所有 go pkg ，根据注解生成结构体相关代码。

- `iocli list`

  查看应用所有接口和方法信息，默认端口 :1999

- `iocli watch [structID] [methodName]`

  监听一个方法的实时调用信息

- `iocli trace [structID] [methodName]`

  以当前方法为入口，开启调用链路追踪

具体操作参数可通过 -h 查看。



