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

import (
	"fmt"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:constructFunc=InitBankService
// +ioc:tx:func=AddMoney-AddMoneyRollback
// +ioc:tx:func=RemoveMoney-RemoveMoneyRollback

type BankService struct {
	Money map[int]int
}

func InitBankService(b *BankService) (*BankService, error) {
	b.Money = make(map[int]int)
	b.Money[1] = 100
	b.Money[2] = 100
	return b, nil
}

func (b *BankService) GetMoney(id int) int {
	return b.Money[id]
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
