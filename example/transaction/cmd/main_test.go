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
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/alibaba/ioc-golang"
)

func (a *App) TestRun(t *testing.T) {
	err := a.TradeService.DoTradeWithTxAddMoneyFailed(1, 2, 10)
	assert.NotNil(t, err)
	assert.Equal(t, "add money num -1 is not positive", err.Error())
	assert.Equal(t, 100, a.BankService.GetMoney(1))
	assert.Equal(t, 100, a.BankService.GetMoney(2))

	err = a.TradeService.DoTradeWithTxFinallyFailed(1, 2, 10)
	assert.NotNil(t, err)
	assert.Equal(t, "finally failed", err.Error())
	assert.Equal(t, 100, a.BankService.GetMoney(1))
	assert.Equal(t, 100, a.BankService.GetMoney(2))

	err = a.TradeService.DoTradeWithTxSuccess(1, 2, 10)
	assert.Nil(t, err)
	assert.Equal(t, 0, a.BankService.GetMoney(1))
	assert.Equal(t, 200, a.BankService.GetMoney(2))
}

func TestTransactionRollback(t *testing.T) {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}
	app.TestRun(t)
}
