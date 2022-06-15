# 使用 RPC 能力

### 简介

本示例展示了基于 IOC-golang 框架的 RPC 能力

在微服务开发过程中，暴露一些对象的方法给外部程序调用是一种常见的场景。

在 IOC-golang 的 RPC 能力的使用场景下，客户端可以直接注入下游接口的客户端存根，以方法调用的形式发起RPC请求。

### 关键代码

**服务端**

需要暴露的RPC服务结构使用  `// +ioc:autowire:type=rpc` 标注，例如 service.go：

```go
import (
	"github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/dto"
)

// +ioc:autowire=true
// +ioc:autowire:type=rpc

type ServiceStruct struct {
}

func (s *ServiceStruct) GetUser(name string, age int) (*dto.User, error) {
	return &dto.User{
		Id:   1,
		Name: name,
		Age:  age,
	}, nil
}

```

使用 iocli 工具生成相关代码。

`sudo iocli gen`

```bash
% tree
.
├── api
│   └── zz_generated.ioc_rpc_client_servicestruct.go
├── service.go
└── zz_generated.ioc.go
```

会在当前文件目录下生成 `zz_generated.ioc.go` 包含了服务提供者的结构描述信息。也会在当前目录下创建 api/ 文件夹，并创建当前结构的客户端存根文件 `zz_generated.ioc_rpc_client_servicestruct.go` 

**客户端**

可以通过标签注入的方法，注入客户端存根，存根中给出下游地址。默认服务暴露端口为`2022`

```go
import(
  "github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/service/api"
)
// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	ServiceStruct *api.ServiceStructIOCRPCClient `rpc-client:",address=127.0.0.1:2022"`
}

func (a *App) Run() {
	for {
		time.Sleep(time.Second * 3)
		usr, err := a.ServiceStruct.GetUser("laurence", 23) // RPC调用
		if err != nil {
			panic(err)
		}
		fmt.Printf("get user = %+v\n", usr)
	}
}
```

### 运行示例

**启动服务端**

服务端需要确保引入对应服务提供者结构包： `_ "github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/service"`

```go
package main

import (
	"github.com/alibaba/ioc-golang"
	_ "github.com/alibaba/ioc-golang/example/autowire_rpc/server/pkg/service"
)

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	select {}
}

```

启动服务端进程，默认 rpc 服务监听 2022 端口

```bash
% cd server
% go run .
...
[negroni] listening on :2022
```

**启动客户端**

开启另一个终端，进入client目录，启动进程，不断发起调用，获得返回值。

```bash
% cd client
% go run .
...
get user = &{Id:1 Name:laurence Age:23}
get user = &{Id:1 Name:laurence Age:23}
get user = &{Id:1 Name:laurence Age:23}
```

