# 事务能力

## 简介

本示例展示了 ioc-golag 框架提供的单体应用事务能力、分布式场景下的事务能力。

本事务基于 Saga 的概念，旨在通过一个长事务分解成多个子事务，需要使用事务能力的函数需要额外提供回滚函数，处于事务内的函数一旦返回错误（最后一个参数为 error 并且非 nil），则会触发事务回滚，框架会依次倒序调用之前成功执行函数的回滚方法，从而实现事务的最终一致性。

## 示例场景与回滚原理

本例子为一个交易场景，需要调用 TradeService 提供的方法，从而执行一次交易，从 id 为 1 的用户账户扣除 100 元钱，再给 id 为 2 的用户账户增加 100 元。

```
              -> BankService -> RemoveMoney(id = 1, num = 100)
            /
TradeServce
            \ 
              -> BankService -> AddMoney(id = 2, num = 100)
```



TradeService 提供了三个方法：

```go
type TradeServiceIOCInterface interface {
	DoTradeWithTxAddMoneyFailed(id1, id2, num int) error
	DoTradeWithTxFinallyFailed(id1, id2, num int) error
	DoTradeWithTxSuccess(id1, id2, num int) error
}
```

其中 DoTradeWithTxAddMoneyFailed 在给用户 2 增加余额的过程中出现错误，会触发 id 为 1 的用户扣钱动作的回滚。DoTradeWithTxFinallyFailed 为完成了给 1 扣钱和给 2 增加余额动作，最后抛出错误，会依次出发两次变动的回滚。DoTradeWithTxSuccess 为完成一次成功的转账。

我们提供了单体应用的事务回滚例子，位于 singleton/，也提供了分布式场景下，基于 RPC 的跨进程全链路事务回滚能力，位于 distributed/。

## 关键代码

在 ioc-golang 框架的事务使用中，需要为支持事物的结构体增加注解 ` // +ioc:tx:func=`

这一注解有两类格式，如果值为当前结构的一个函数名，则表示该方法的调用过程是一个事务。

` // +ioc:tx:func=$(functionName)`

如果值为使用 "-" 分隔的两个方法名，则表示该方法的调用过程是一个事务，并且在调用链路出现错误时，会执行回滚函数。

` // +ioc:tx:func=$(functionName)-$(rollbackFunctionName)`

在程序执行过程中，如果调用到的函数标注为事务，则其上层调用栈函数都将被当作子事务来处理，即如果事务内执行过程出现返回值错误，则其上层调用栈的，所有已完成调用，并且标注支持事务回滚的函数，都将被触发回滚。

对于例子中的 BankService ，其标注了两个事务函数 AddMoney，RemoveMoney，并提供了回滚方法 AddMoneyRollback 和 RemoveMoneyRollback。例如在一个事务内，如果 AddMoney 方法成功执行并返回了，但后续的方法中有返回错误，则 AddMoneyRollback 方法将被调用。

```go
// +ioc:autowire=true
// +ioc:autowire:type=rpc
// +ioc:autowire:constructFunc=InitBankService
// +ioc:tx:func=AddMoney-AddMoneyRollout
// +ioc:tx:func=RemoveMoney-RemoveMoneyRollout

type BankService struct {
	Money map[int]int
}

func (b *BankService) AddMoney(id, num int) error {
	if num <= 0 {
		// raise error, this would call all previous succeed branches: RemoveMoneyRollout function
		return fmt.Errorf("add money num %d is not positive", num)
	}
	b.Money[id] += num
	return nil
}

func (b *BankService) AddMoneyRollback(id, num int, errMsg string) {
	b.Money[id] -= num
	fmt.Printf("Transaction is failed, real cause is '%s'\n method BankService.AddMoney is rolling back, sub num %d\n", errMsg, num)
}

func (b *BankService) RemoveMoney(id, num int) error {
	if num <= 0 {
		return fmt.Errorf("remove money num %d is not positive", num)
	}
	b.Money[id] -= num
	return nil
}

func (b *BankService) RemoveMoneyRollback(id, num int, errMsg string) {
	b.Money[id] += num
	fmt.Printf("Transaction is failed, real cause is '%s'\nmethod BankService.RemoveMoney is rolling back, add num %d\n", errMsg, num)
}

```

- 回滚函数的函数签名

  我们以 RemoveMoneyRollback 这一回滚函数为例， 回滚函数签名是和其对应的事务函数签名相关的，RemoveMoneyRollback 的函数签名是 RemoveMoney 函数入参的基础之上，增加了一个 string 类型的参数在最后，其包含了触发回滚的错误信息，回滚函数无需返回值。

```go
 RemoveMoney(id, num int) error // 用户定义的事务函数
 
 RemoveMoneyRollout(id, num int, errMsg string) // RemoveMoney 的回滚函数
```



## 示例简介

在提供的两个示例中，都是通过调用 TradeService 触发交易。其使用注解标注了三个需要使用事务的函数，没有提供回滚函数。

```go
// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:tx:func=DoTradeWithTxFinallyFailed
// +ioc:tx:func=DoTradeWithTxAddMoneyFailed
// +ioc:tx:func=DoTradeWithTxSuccess

type TradeService struct {
	BankRPCService api.BankServiceIOCRPCClient `rpc-client:",address=localhost:2022"`
}

```

在其中的一个函数内，会调用 BankService 提供的 RemoveMoney 和 AddMoney 方法，以DoTradeWithTxAddMoneyFailed 函数为例

```go
func (b *TradeService) DoTradeWithTxAddMoneyFailed(id1, id2, num int) error {
	if err := b.BankService.RemoveMoney(id1, num); err != nil {
		return err
	}

	if err := b.BankService.AddMoney(id2, -1); err != nil {
		// -1 num cause error, previous succeeded branch b.BankService.RemoveMoneyRollout() would be called
		return err
	}
	return nil
}
```

可看到其在调用 AddMoney 的时候传递了错误的数值 -1，这回导致调用失败，返回错误，触发已经完成的函数 RemoveMoney 的回滚逻辑，为用户1 增加回 100 元，保证最终一致性。

对于 distributed/ 给出的例子，其 BankService 结构标注了 rpc 注解，另一个进程作为 RPC server ，它也同样支持事务回滚。当基于调用  AddMoney  失败后，客户端也会以 RPC 的形式调用已完成的回滚函数 RemoveMoneyRollback，从而保证整个分布式系统的最终一致性。

## 运行例子

- 单体应用 

```bash
% cd singleton/cmd
% go run .
...
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/aop/transaction/singleton/service.BankService
[Autowire Struct Descriptor] Found type singleton registered SD github.com/alibaba/ioc-golang/example/aop/transaction/singleton/service.TradeService
[Autowire Struct Descriptor] Found type singleton registered SD main.App
user 1 have  100
user 2 have  100
---
Transaction is failed, real cause is 'add money num -1 is not positive'
method BankService.RemoveMoney is rolling back, add num 100
ops! DoTradeWithTxAddMoneyFailed failed with error = add money num -1 is not positive
user 1 have  100
user 2 have  100
---
Transaction is failed, real cause is 'finally failed'
 method BankService.AddMoney is rolling back, sub num 100
Transaction is failed, real cause is 'finally failed'
method BankService.RemoveMoney is rolling back, add num 100
ops! DoTradeWithTxFinallyFailed failed with error = finally failed
user 1 have  100
user 2 have  100
---
user 1 have  0
user 2 have  200

```

可看到前两次触发交易都失败了，触发回滚的错误原因被打印了出来，账户金额在回滚后没有变化。第三次调用没有报错，成功完成了交易。

- 基于 RPC 的分布式事务

启动 server 

```
% cd distributed/server
% go run .
....
[GIN-debug] POST   /github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service/api.BankServiceIOCRPCClient/RemoveMoneyRollback --> github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl.(*IOCProtocol).Export.func1 (3 handlers)
[GIN-debug] Listening and serving HTTP on :2022
```

客户端发起调用

```bash
% go run .
...
[Autowire Struct Descriptor] Found type singleton registered SD main.App
[Debug] Debug server listening at :2000
trade failed, error =  finally failed
user 1 have money 100
user 2 have money 100
trade failed error =  add money num -1 is not positive
user 1 have money 100
user 2 have money 100
trade success
user 1 have money 0
user 2 have money 200
```

可看到符合预期，server 端返回错误后，金额按照预期被回滚，错误信息被正确打印。

- 观测 RPC 场景下回滚函数的调用情况

我们也可以使用 ioc-golang 的可观测性，监控一下 server 端的接口被调用情况。使用 iocli 的 list 和 monitor 方法。

```bash
% iocli list                                                                                                     
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService
[RemoveMoneyRollback GetMoney AddMoney AddMoneyRollback RemoveMoney]

github.com/alibaba/ioc-golang/extension/autowire/rpc/protocol/protocol_impl.IOCProtocol
[Invoke Export]

% iocli monitor -i 3 github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService 
====================
2022/07/09 21:18:41
====================
2022/07/09 21:18:44

```

目前没有调用，三秒钟会打印一次监控信息，我们开启 client 端尝试调用一下， 可在监控侧查看到日志：

```bash
====================
2022/07/09 21:20:11
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService.AddMoney()
Total: 3, Success: 2, Fail: 1, AvgRT: 0.00ms, FailRate: 33.33%
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService.AddMoneyRollback()
Total: 1, Success: 1, Fail: 0, AvgRT: 0.00ms, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService.GetMoney()
Total: 6, Success: 6, Fail: 0, AvgRT: 0.00ms, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService.RemoveMoney()
Total: 3, Success: 3, Fail: 0, AvgRT: 0.00ms, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service.BankService.RemoveMoneyRollback()
Total: 2, Success: 2, Fail: 0, AvgRT: 0.00ms, FailRate: 0.00%

```

可观测到，在这三秒中，AddMoney 失败了一次，回滚函数 RemoveMoneyRollback 被执行了两次 回滚函数 AddMoneyRollback 被执行了一次，符合我们的预期。

