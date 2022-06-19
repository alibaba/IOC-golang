# 使用基于 Monkey 的调试功能

### 简介

基于 monkey 的调试能力，可以为所有结构体指针也增加代理层，从而赋予运维能力。

但由于 monkey 指针的限制，这种方式对性能有损耗，**不推荐在生产环境中使用**。

### 例子解析

在本例子中，我们将 ServiceImpl1 以结构指针的方式注入，对于另一个结构体 ServiceImpl2，我们以专属接口的形式注入。按照正常的启动方式，注入至接口的结构都封装了代理层，我们可以调试 ServiceImpl2 ，而不能调试 ServiceImpl1。

```go
// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceImpl1Ptr       *ServiceImpl1            `singleton:""`
	ServiceImpl2Interface ServiceImpl2IOCInterface `singleton:""`
}
```

使用了 Monkey 的调试功能，我们可以同时满足调试 ServiceImpl1 和 ServiceImpl2。

### 运行例子

1. 按照正常方式启动

   ```bash
    % go run .
     ___    ___     ____                           _                         
    |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
     | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
     | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
    |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                   |___/                               |___/ 
   Welcome to use ioc-golang!
   [Boot] Start to load ioc-golang config
   [Config] Load default config file from ../conf/ioc_golang.yaml
   [Config] Load ioc-golang config file failed. open /Users/laurence/Desktop/workplace/alibaba/IOC-Golang/example/conf/ioc_golang.yaml: no such file or directory
    The load procedure is continue
   [Boot] Start to load debug
   [Debug] Debug port is set to default :1999
   [Boot] Start to load autowire
   [Autowire Type] Found registered autowire type normal
   [Autowire Struct Descriptor] Found type normal registered SD main.app_
   [Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl1_
   [Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl2_
   [Autowire Type] Found registered autowire type singleton
   [Autowire Struct Descriptor] Found type singleton registered SD main.App
   [Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl1
   [Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl2
   [Debug] Debug server listening at :1999
   This is ServiceImpl1, hello laurence
   This is ServiceImpl2, hello laurence
   This is ServiceImpl1, hello laurence
   ...
   ```

   另一个终端查看所有接口方法：

   ```bash
   % iocli list
   main.App
   [Run]
   
   main.ServiceImpl2
   [GetHelloString]
   
   ```

   可看到只能查看到具有代理层的 App 结构和 main.ServiceImpl2 结构。

2. 使用 monkey 调试模式

   关闭上述终端，重新以 monkey 调试模式启动程序，构建参数和标签 `-gcflags="-N -l" -tags iocMonkeydebug`都是必要的。CPU架构如不是 amd64，需要通过环境变量 GOARCH 来修改成 amd64 跨平台编译。

   ```bash
   % GOARCH=amd64 go run -gcflags="-N -l" -tags iocMonkeydebug .
     ___    ___     ____                           _                         
    |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
     | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
     | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
    |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                   |___/                               |___/ 
   Welcome to use ioc-golang!
   [Boot] Start to load ioc-golang config
   [Config] Load default config file from ../conf/ioc_golang.yaml
   [Config] Load ioc-golang config file failed. open /Users/laurence/Desktop/workplace/alibaba/IOC-Golang/example/conf/ioc_golang.yaml: no such file or directory
    The load procedure is continue
   [Boot] Start to load debug
   [Debug] Debug port is set to default :1999
   [Boot] Start to load autowire
   [Autowire Type] Found registered autowire type singleton
   [Autowire Struct Descriptor] Found type singleton registered SD main.App
   [Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl1
   [Autowire Struct Descriptor] Found type singleton registered SD main.ServiceImpl2
   [Autowire Type] Found registered autowire type normal
   [Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl1_
   [Autowire Struct Descriptor] Found type normal registered SD main.serviceImpl2_
   [Autowire Struct Descriptor] Found type normal registered SD main.app_
   [Debug] Debug server listening at :1999
   This is ServiceImpl1, hello laurence
   This is ServiceImpl2, hello laurence
   ...
   ```

   查看接口和方法：

   ```bash
   % iocli list
   main.App
   [Run]
   
   main.ServiceImpl1
   [GetHelloString]
   
   main.ServiceImpl2
   [GetHelloString]
   ```

   可以看到，main.ServiceImpl1 方法也具备了运维能力，我们监听该接口的 GetHelloString 方法，得到如下信息：

   ```bash
    iocli watch main.ServiceImpl1 GetHelloString
   ========== On Call ==========
   main.ServiceImpl1.GetHelloString()
   Param 1: (string) (len=8) "laurence"
   
   ========== On Response ==========
   main.ServiceImpl1.GetHelloString()
   Response 1: (string) (len=36) "This is ServiceImpl1, hello laurence"
   
   ```

   可看到，基于 monkey 指针的调试模式，可以对结构体指针进行代理层注入。
