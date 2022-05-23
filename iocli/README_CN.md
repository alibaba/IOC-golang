# iocli 工具

iocli 是一款命令行工具，提供了以下能力：

- 代码调试

  开发者可以配置启动 ioc-go 内核提供的调试能力，iocli 作为调试客户端。

- 结构描述注册信息生成

  开发者可以为需要依赖注入的结构体增加注解，iocli 会识别这些注解，并产生结构描述符，注册在 ioc-go 框架内。

## 代码调试能力

ioc-golang 框架拥有首创的基于 AOP 思路的 Go 运行时程序调试能力。

### 开启代码调试能力

我们以的 demo 为基础，开启代码调试能力，并使用 iocli 进行调试。

修改 main 函数 Run 方法调用方式为隔 5s 调用一次

```go
func main(){
	// 框架启动
	if err := ioc.Load(); err != nil{
		panic(err)
	}

	// App-App 即结构ID： '$(接口名)-$(方法名)'， 对于结构指针，接口名默认为方法名
	// 可通过这一 ID 获取实例
	appInterface, err := singleton.GetImpl("App-App")
	if err != nil{
		panic(err)
	}
	app := appInterface.(*App)
	for{
		time.Sleep(time.Second*5) // 5s 调用一次
		app.Run()
	}
}
```

指定配置文件环境变量

```bash
export IOC_GOLANG_CONFIG_PATH=./ioc_golang.yaml
```

ioc_golang.yaml:

```yaml
debug:
  enable: true # 开启调试模式
```

```bash
% tree .
.
├── ioc_golang.yaml
├── go.mod
├── go.sum
├── main.go
└── zz_generated.ioc.go
```

启动程序：非amd64 机器需要指定架构，GOARCH=amd64。-gcflags="-N -l" 关闭编译器优化内联。

```bash
 %  GOARCH=amd64 go run -gcflags "-N -l" .
  ___    ___     ____            ____           _                         
 |_ _|  / _ \   / ___|          / ___|   ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____  | |  _   / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | |_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \____|  \___/  |_|  \__,_| |_| |_|  \__, |
                                                                    |___/ 
Welcome to use ioc-golang!
[Boot] Start to load ioc-golan config
[Config] Environment IOC_GOLANG_CONFIG_PATH is set to ./ioc_golang.yaml
[Config] Load config file from ./ioc-golang.yaml
[Boot] Start to load debug
[Debug] Debug port is set to default :1999
[Boot] Start to load autowire
[Autowire Type] Found registered autowire type singleton
[Autowire Struct Descriptor] Found type singleton registered SD App-App
[Autowire Struct Descriptor] Found type singleton registered SD Service-ServiceImpl1
[Autowire Struct Descriptor] Found type singleton registered SD Service-ServiceImpl2
[Autowire Struct Descriptor] Found type singleton registered SD ServiceStruct-ServiceStruct
[Debug] Debug server listening at :1999
This is ServiceImpl1, hello world
This is ServiceImpl2, hello world
This is ServiceStruct, hello world
```

可看到 debug 端口监听1999

新启动一个终端，查看所有接口实现和方法：

```bash
$ iocli list
App
App
[Run]

Service
ServiceImpl1
[Hello]

Service
ServiceImpl2
[Hello]

ServiceStruct
ServiceStruct
[Hello]
```

监听一个实现类的方法：

```bash
% iocli watch Service ServiceImpl1 Hello

========== On Call ==========
Service.(ServiceImpl1).Hello()
```

对于有入参和返回值的参数，可以监听到具体参数类型和值。一个监听 grpc 客户端的例子如下：

```bash
========== On Call ==========
HelloServiceClient.(HelloServiceClient).SayHello()
Param 1: (*context.emptyCtx)(0xc0000a0000)(context.Background)

Param 2: (*api.HelloRequest)(0xc00050c640)(name:"laurence")

Param 3: ([]grpc.CallOption) (len=2 cap=2) {
 (grpc.MaxRecvMsgSizeCallOption) {
  MaxRecvMsgSize: (int) 1024
 },
 (grpc.MaxRecvMsgSizeCallOption) {
  MaxRecvMsgSize: (int) 1024
 }
}


========== On Response ==========
HelloServiceClient.(HelloServiceClient).SayHello()
Response 1: (*api.HelloResponse)(0xc00050c700)(reply:"Hello laurence")

Response 2: (interface {}) <nil>


========== On Call ==========
HelloServiceClient.(HelloServiceClient).SayHello()
Param 1: (*context.emptyCtx)(0xc0000a0000)(context.Background)

Param 2: (*api.HelloRequest)(0xc00050c8c0)(name:"laurence_service1_impl1")

Param 3: ([]grpc.CallOption) <nil>


========== On Response ==========
HelloServiceClient.(HelloServiceClient).SayHello()
Response 1: (*api.HelloResponse)(0xc00050c980)(reply:"Hello laurence_service1_impl1")

Response 2: (interface {}) <nil>


```



基于该思路，我们可以扩展更丰富的调试能力，例如：

- 流量过滤、监控
- 参数编辑
- 故障注入
- 耗时瓶颈分析
- ...



## 结构注解

iocli 可以识别以下注解：

```go
// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:interface=Redis
// +ioc:autowire:paramLoader=paramLoader
// +ioc:autowire:paramType=Config
// +ioc:autowire:constructFunc=New
// +ioc:autowire:baseType=true
```

- ioc:autowire 

  bool 类型，为 true 则在代码生成阶段被识别到。

- ioc:autowire:type

  string类型，表示依赖注入模型，目前支持四种：

  - singleton

    单例模型，该结构体全局只能产生一个对象，无论是 API 获取还是字段注入。

  - normal

    多例模型，每一个标签注入字段、每一次 API 获取，都会产生一个新的对象。

  - config:

    配置模型是基于多例模型的封装扩展，基于配置模型定义的结构体方便从 yaml 配置文件中注入信息。

  - grpc:

    grpc 模型是基于单例模型的封装扩展，基于 grpc 模型可以方便地从 yaml 配置文件中读取参数，生成 grpc 客户端。

- ioc:autowire:interface（非必填）

  string类型，表示实现的接口名，如果不存在这个标注，将作为结构体指针注入给使用方。

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
  	NormalDB3Redis normalRedis.Redis `normal:"Impl,address=127.0.0.1:6379&db=3"`
  }
  ```

  当然也可以从配置文件读入，tag中指定了key为 db1-redis

  ```go
  type App struct {
  	NormalDB3Redis normalRedis.Redis `normal:"Impl,db1-redis"`
  }
  ```

  ioc-go.yaml： autowire.normal.Redis.Impl.db1-redis.param 读入参数

  ```yaml
  autowire:
    normal:
      Redis:
        Impl:
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

  string类型，表示结构的构造方法名，作为依赖参数的一个函数。

  在上述例子中，该方法名为 New。

  思路为：依赖参数是构造对象的前提，因此该方法放置在参数下，会被框架调用。对于不依赖外部传递参数，但拥有构造函数的对象，可以使用一个空Struct 作为参数，在这一Struct 内实现构造方法。

  我们要求该构造方法的函数签名是固定的，即：

  ```go
  func (*$(结构名)) (*$(结构名）, error)
  ```

- ioc:autowire:baseType=true （非必填）

  该类型是否为基础类型

  go 基础类型不可直接通过&在构造时取地址，因此我们针对基础类型单独设计了该注解。在 http://github.com/alibaba/ioc-golang/extension/tree/master/config 配置扩展中被使用较多。





