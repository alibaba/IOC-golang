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

package service

import "fmt"

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:tx:func=DoTradeWithTxAddMoneyFailed
// +ioc:tx:func=DoTradeWithTxFinallyFailed
// +ioc:tx:func=DoTradeWithTxSuccess

type TradeService struct {
	BankService BankServiceIOCInterface `singleton:""`
}

func (b *TradeService) DoTradeWithTxAddMoneyFailed(id1, id2, num int) error {
	if err := b.BankService.RemoveMoney(id1, 100); err != nil {
		return err
	}

	if err := b.BankService.AddMoney(id2, -1); err != nil {
		// -1 num cause error, previous succeeded branch b.BankService.RemoveMoneyRollout() would be called
		return err
	}
	return nil
}
func (b *TradeService) DoTradeWithTxFinallyFailed(id1, id2, num int) error {
	if err := b.BankService.RemoveMoney(id1, 100); err != nil {
		return err
	}

	if err := b.BankService.AddMoney(id2, 100); err != nil {
		return err
	}
	// previous succeeded branch b.BankService.AddMoneyRollout() and b.BankService.RemoveMoneyRollout() would be called in order
	return fmt.Errorf("finally failed")
}

func (b *TradeService) DoTradeWithTxSuccess(id1, id2, num int) error {
	if err := b.BankService.RemoveMoney(id1, 100); err != nil {
		return err
	}

	if err := b.BankService.AddMoney(id2, 100); err != nil {
		return err
	}
	return nil
}
