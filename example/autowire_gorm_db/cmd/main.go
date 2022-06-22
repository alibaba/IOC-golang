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
	normalMysql "github.com/alibaba/ioc-golang/extension/db/gorm"
)

// +ioc:autowire=true
// +ioc:autowire:type=singleton
// +ioc:autowire:alias=AppAlias
type App struct {
	MyDB normalMysql.GORMDBIOCInterface `normal:",my-mysql"`
}

type MyDataDO struct {
	Id    int32
	Value string
}

func (a *MyDataDO) TableName() string {
	return "mydata"
}

func (a *App) Run() {
	// create table
	if err := a.MyDB.Model(&MyDataDO{}).AutoMigrate(&MyDataDO{}); err != nil {
		panic(err)
	}
	toInsertMyData := &MyDataDO{
		Value: "first value",
	}
	if err := a.MyDB.Model(&MyDataDO{}).Create(toInsertMyData).Error(); err != nil {
		panic(err)
	}
	myDataDOs := make([]MyDataDO, 0)
	if err := a.MyDB.Model(&MyDataDO{}).Where("id = ?", 1).Find(&myDataDOs).Error(); err != nil {
		panic(err)
	}
	fmt.Println(myDataDOs)
}

func main() {
	if err := ioc.Load(); err != nil {
		panic(err)
	}
	app, err := GetAppSingleton()
	if err != nil {
		panic(err)
	}

	app.Run()
}
