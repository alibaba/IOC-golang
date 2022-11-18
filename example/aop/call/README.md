# 触发接口方法调用

### 简介

基于 IOC-golang 框架注入的接口默认具备 AOP 能力，我们可以使用 iocli call 命令触发这些接口的方法调用，方便调试和问题排查。

对于一个正在运行的基于 IOC-golang 的程序，我们可以通过 iocli list 命令获取到其全部接口方法，从而按需触发方法调用。


### 运行例子：

1. 在当前命令下执行：`go run .` 启动程序

```go
% go run .
  ___    ___     ____                           _                         
 |_ _|  / _ \   / ___|           __ _    ___   | |   __ _   _ __     __ _ 
  | |  | | | | | |      _____   / _` |  / _ \  | |  / _` | | '_ \   / _` |
  | |  | |_| | | |___  |_____| | (_| | | (_) | | | | (_| | | | | | | (_| |
 |___|  \___/   \____|          \__, |  \___/  |_|  \__,_| |_| |_|  \__, |
                                |___/                               |___/ 
Welcome to use ioc-golang 1.0.2!
[Boot] Start to load ioc-golang config
[Config] Config files load options is {AbsPath:[] ConfigName:config ConfigType:yaml SearchPath:[. ./config ./configs] ProfilesActive:[] MergeDepth:8 Properties:map[]}
[Boot] Start to load debug
[AOP Log] log config level is using default 'info'
[AOP Log] log config print-params-max-depth is set to default 3
[AOP Log] log config invocation ctx logs level is using default 'info'
[AOP Log] log config invocation ctx logs print level is using default 'debug'
[Wrapper Autowire] Parse param from config file with sdid github.com/alibaba/ioc-golang/extension/aop/log.logInterceptor failed, error: property [autowire singleton github.com/alibaba/ioc-golang/extension/aop/log.logInterceptor param]'s key autowire not found, continue with nil param.
[AOP] Debug server port is set to default :1999
[Wrapper Autowire] Parse param from config file with sdid github.com/alibaba/ioc-golang/extension/aop/log.logInterceptor failed, error: property [autowire singleton github.com/alibaba/ioc-golang/extension/aop/log.logInterceptor param]'s key autowire not found, continue with nil param.
[Boot] Start to load autowire
[Debug] Debug server listening at :1999
[Autowire Type] Found registered autowire type normal
...
[AOP] [Trace] Set trace logger to logrus success
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 }
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 }
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 }

```


2. 新开一个终端，使用 iocli list 命令查询接口方法
```bash
% iocli list
main.App
[Run]

main.ServiceImpl1
[GetHelloString]

main.ServiceImpl2
[GetHelloString]

main.UserService
[CreateUser ParseUserInfo SetMark]

```

可看到当前程序包含的业务接口和方法，我们尝试调用 main.UserService 结构的 SetMark 方法，从而修改其 mark 字段

```go
// +ioc:autowire=true
// +ioc:autowire:type=singleton

type UserService struct {
	mark string
}

func (s *UserService) SetMark(mark string) {
	s.mark = mark // 修改mark字段
}
```

3. iocli call 命令用法

```bash
% iocli call 
iocli call command needs 3 arguments: 
${autowireType} ${sdid} ${methodName} #iocli call 命令用法
% iocli call -h
 iocli call -h
Usage:
  iocli call [flags]

Flags:
  -h, --help                help for call
      --host string         debug host (default "127.0.0.1")
      --params string       request call params json string # 以 json 字符串形式通过终端传入参数
      --paramsFile string   request call params json file path # 以 json 文件形式传入函数参数
  -p, --port int            debug port (default 1999)
```

调用上述 SetMark 方法

```bash
% iocli call singleton main.UserService SetMark --params "[\"mymarks\"]"
Call singleton: main.UserService.SetMark() success!
Param = ["mymarks"]
Return values = []
```

可看到调用成功，打印AOP调试触发日志 `Receive call request`，程序原本日志出现变化，验证调用成功，mark 被设置为 `mymarks`

```bash
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 }
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 }
[Debug Server] Receive call request sdid:"main.UserService" methodName:"SetMark" autowireType:"singleton" params:"[\"mymarks\"]"
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 mymarks}
This is ServiceImpl2, hello laurence
create user &{1 laurence 23 mymarks}
```

