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
	"github.com/alibaba/ioc-golang/example/transaction/service"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton

type App struct {
	TradeService service.TradeServiceIOCInterface `singleton:""`

	BankService service.BankServiceIOCInterface `singleton:""`
}

func (a *App) printMoney() {
	fmt.Println("user 1 have ", a.BankService.GetMoney(1))
	fmt.Println("user 2 have ", a.BankService.GetMoney(2))
}

func (a *App) Run() {
	a.printMoney()
	fmt.Println("---")
	if err := a.TradeService.DoTradeWithTxAddMoneyFailed(1, 2, 10); err != nil {
		fmt.Printf("ops! DoTradeWithTxAddMoneyFailed failed with error = %s\n", err)
		a.printMoney()
		fmt.Println("---")
	}

	if err := a.TradeService.DoTradeWithTxFinallyFailed(1, 2, 10); err != nil {
		fmt.Printf("ops! DoTradeWithTxFinallyFailed failed with error = %s\n", err)
		a.printMoney()
		fmt.Println("---")

	}

	if err := a.TradeService.DoTradeWithTxSuccess(1, 2, 10); err != nil {
		fmt.Printf("ops! DoTradeWithTxSuccess failed with error = %s\n", err)
		a.printMoney()
		return
	}
	a.printMoney()
}

func main() {
	// start
	if err := ioc.Load(); err != nil {
		panic(err)
	}

	// app, err := GetAppIOCInterfaceSingleton is ok too
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.Run()
}
