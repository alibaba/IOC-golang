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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang"

	_ "github.com/alibaba/ioc-golang/example/aop/transaction/distributed/server/pkg/service"
)

func (a *App) TestRun(t *testing.T) {
	assert.NotNil(t, a.TradeService.DoTradeWithTxFinallyFailed(1, 2, 100))
	assert.Equal(t, 100, a.BankRPCService.GetMoney(1))
	assert.Equal(t, 100, a.BankRPCService.GetMoney(2))

	assert.NotNil(t, a.TradeService.DoTradeWithTxAddMoneyFailed(1, 2, 100))
	assert.Equal(t, 100, a.BankRPCService.GetMoney(1))
	assert.Equal(t, 100, a.BankRPCService.GetMoney(2))

	assert.Nil(t, a.TradeService.DoTradeWithTxSuccess(1, 2, 100))
	assert.Equal(t, 0, a.BankRPCService.GetMoney(1))
	assert.Equal(t, 200, a.BankRPCService.GetMoney(2))
}

func TestRPCClient(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	app.TestRun(t)
}
