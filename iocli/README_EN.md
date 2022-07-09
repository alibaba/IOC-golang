# iocli tool

###  English | [中文](./README.md)

**iocli** is a command line tool, it has the following features: 

- program debug

  Developers can use **iocli** as debug client and debug the go program developed with ioc-golang.

- codes generation

  Developers can add annotations to structures that describe the struct to be   injected, and **iocli** will identify these annotations and generates structure-specific code that meets the requirements. Including structure description information, structure proxy layer, structure own interface, structure Get method and so on.

## Program Debug Feature

IOC golang framework has the ability to debug go programs based on the struct AOP layer, help with troubleshooting, performance analysis, and improve the observability of applications. In [README]( https://github.com/alibaba/ioc-golang#ioc -Golang-a-golang-dependency-injection-framework) QuickStart page shows the ability to view interface information and monitor parameters. In [IOC golang based e-commerce system demo]( https://github.com/ioc-golang/shopping-system ). It shows the IOC-golang based, non intrusive, method granularity whole invoking link tracking capability.

## Annotation and Code generation

An annotation is an annotation that begins with a specific string and is marked in front of the structure that is expected to be injected. Annotations only have static meaning, that is, in the code generation stage, they are scanned and recognized by iocli tools to obtain structure related information. Annotations themselves do not have the meaning of program runtime.

iocli can identify the following annotation keys, and the values after '=' are just for example.

```go
// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:paramLoader=paramLoader
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New
// +ioc:autowire:baseType=true
// +ioc:autowire:alias=MyAppAlias
// +ioc:tx:func=MyTransactionFunction
```

- ioc:autowire (required)

  bool type, The identify flag for code generation to pick the struct.

- ioc:autowire:type (required)

  string type. It represents the autowire  model. Currently, it supports the following five types. The structure provider can select one or more of the five types to annotate, so as to generate relevant structure information and APIs for structure users to choose.

  ```go
  // +ioc:autowire:type=singleton
  // +ioc:autowire:type=normal
  
  type MyStruct struct{
  
  }
  ```

  - singleton

    For the singleton model, there is only one object globally in the structure obtained by using the injection model.

  - normal

    Multiple case model. With this injection model, each tag injection field and API acquisition will generate a new object.

  - config:

     [example/autowire_config](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_config)

  - grpc:

    The configuration model is a encapsulation extension based on the multi instance model. The structure defined based on the configuration model is convenient to inject information from the yaml configuration file. Reference to [example/autowire_grpc_client](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_grpc_client)

  - rpc:

    RPC model will generate RPC server registration code and RPC client API stub in the code generation phase. Reference to  [example/autowire_rpc](https://github.com/alibaba/IOC-golang/tree/master/example/autowire/autowire_rpc)

  

- ioc:autowire:paramLoader (not required)

  String type, indicating the type name of "parameter loader" that needs to be customized

  

  The parameter loader is optionally customized by the structure definer. Refer to [ioc-go-extension/normal/redis](http://github.com/alibaba/ioc-golang/extension/blob/master/normal)

  param loader should import Load method：
  
  ```go
  // ParamLoader is interface to load param
  type ParamLoader interface {
  	Load(sd *StructDescriptor, fi *FieldInfo) (interface{}, error)
  }
  ```

  Developers who define structures can define their own structure initialization parameters by implementing parameter loaders. For example, a redis client structure'impl'needs to be loaded from the config parameter, as shown in the new method below.
  
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

  Config contains three fields: address password dB, which need to be passed in by the user.

  From where? This is what the parameter loader does.

  The structure definer can define the following loaders, so that the fields can be obtained through the tag tag injected into the structure. If the tag information marks the configuration location, it can be obtained through the configuration file.
  
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

  Such as
  
  ```go
  type App struct {
  	NormalDB3Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/state/redis.Redis,address=127.0.0.1:6379&db=3"`
  }
  ```

  Of course, it can also be read from the configuration file. The key specified in the tag is db1 redis
  
  ```go
  type App struct {
  	NormalDB3Redis normalRedis.Redis `normal:"github.com/alibaba/ioc-golang/extension/state/redis.Redis,db1-redis"`
  }
  ```

  ioc-go.yaml： autowire.normal.Redis.Impl.db1-redis.param  Read from the config file.
  
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

## iocli 操作命令

- `iocli init`

  生成初始化工程

- `iocli gen`

  递归遍历当前目录下的所有 go pkg ，根据注解生成结构体相关代码。

- `iocli list`

  查看应用所有接口和方法信息，默认端口 :1999

- `iocli watch [structID] [methodName]`

  监听一个方法的实时调用信息


- `iocli monitor [structID] [methodName]`

  开启调用监控，structID 和 methodName 可不指定，则监控所有接口方法。

- `iocli trace [structID] [methodName]`

  以当前方法为入口，开启调用链路追踪

具体操作参数可通过 -h 查看。



