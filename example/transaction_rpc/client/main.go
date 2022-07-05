/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"

	"github.com/alibaba/ioc-golang"
	"github.com/alibaba/ioc-golang/example/transaction_rpc/server/pkg/service/api"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:tx:func=DoTradeWithTxFinallyFailed
// +ioc:tx:func=DoTradeWithTxAddMoneyFailed
// +ioc:tx:func=DoTradeWithTxSuccess

type TradeService struct {
	BankRPCService api.BankServiceIOCRPCClient `rpc-client:",address=localhost:2022""`
}

func (b *TradeService) DoTradeWithTxFinallyFailed(id1, id2, num int) error {
	if err := b.BankRPCService.RemoveMoney(id1, num); err != nil {
		return err
	}

	if err := b.BankRPCService.AddMoney(id2, num); err != nil {
		return err
	}
	// previous succeeded branch b.BankService.AddMoneyRollout() and b.BankService.RemoveMoneyRollout() would be called in order
	return fmt.Errorf("finally failed")
}

func (b *TradeService) DoTradeWithTxAddMoneyFailed(id1, id2, num int) error {
	if err := b.BankRPCService.RemoveMoney(id1, num); err != nil {
		return err
	}

	if err := b.BankRPCService.AddMoney(id2, -1); err != nil {
		// -1 num cause error, previous succeeded branch b.BankService.RemoveMoneyRollout() would be called
		return err
	}
	return nil
}

func (b *TradeService) DoTradeWithTxSuccess(id1, id2, num int) error {
	if err := b.BankRPCService.RemoveMoney(id1, num); err != nil {
		return err
	}

	if err := b.BankRPCService.AddMoney(id2, num); err != nil {
		return err
	}
	return nil
}

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	BankRPCService api.BankServiceIOCRPCClient `rpc-client:",address=localhost:2022"`
	TradeService   TradeServiceIOCInterface    `singleton:""`
}

func (a *App) Run() {
	if err := a.TradeService.DoTradeWithTxFinallyFailed(1, 2, 100); err != nil {
		// finally failed
		fmt.Println("trade failed, error = ", err)
		fmt.Println("user 1 have money", a.BankRPCService.GetMoney(1))
		fmt.Println("user 2 have money", a.BankRPCService.GetMoney(2))
	}

	if err := a.TradeService.DoTradeWithTxAddMoneyFailed(1, 2, 100); err != nil {
		// finally failed
		fmt.Println("trade failed error = ", err)
		fmt.Println("user 1 have money", a.BankRPCService.GetMoney(1))
		fmt.Println("user 2 have money", a.BankRPCService.GetMoney(2))
	}

	if err := a.TradeService.DoTradeWithTxSuccess(1, 2, 100); err != nil {
		panic(err)
	}
	// finally failed
	fmt.Println("trade success")
	fmt.Println("user 1 have money", a.BankRPCService.GetMoney(1))
	fmt.Println("user 2 have money", a.BankRPCService.GetMoney(2))
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
