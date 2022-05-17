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

package mysql

import (
	"gorm.io/gorm"

	"github.com/alibaba/IOC-Golang/autowire/normal"
)

type Mysql interface {
	GetDB() *gorm.DB
	SelectWhere(queryStr string, result interface{}, args ...interface{}) error
	Insert(toInsertLines UserDefinedModel) error
	Delete(toDeleteTarget UserDefinedModel) error
	First(queryStr string, findTarget UserDefinedModel, args ...interface{}) error
	Update(queryStr, field string, target interface{}, args ...interface{}) error
}

const SDID = "Mysql-Impl"

// +ioc:autowire=true
// +ioc:autowire:interface=Mysql
// +ioc:autowire:type=normal
// +ioc:autowire:paramType=Config
// +ioc:autowire:paramLoader=paramLoader
// +ioc:autowire:constructFunc=New

type Impl struct {
	tableName string
	db        *gorm.DB
}

func (i *Impl) GetDB() *gorm.DB {
	return i.db
}

type UserDefinedModel interface {
	TableName() string
}

// SelectWhere 两个参数: queryStr result
// 调用次函数，相当于针对当前table执行 select * from table where `queryStr`，例如`userId = ?`, args = "1"
// 将结果写入result， result类型只能是注册好的model数组，类型为 &[]UserDefinedModel{}
func (mt *Impl) SelectWhere(queryStr string, result interface{}, args ...interface{}) error {
	if err := mt.db.Table(mt.tableName).Where(queryStr, args...).Find(result).Error; err != nil {
		return err
	}
	return nil
}

// Insert 一个参数 toInsertLines
// 调用次函数，相当于针对当前table，插入toInsertLines 对应的数据
// toInsertLines类型为 UserDefinedModel
func (mt *Impl) Insert(toInsertLines UserDefinedModel) error {
	if err := mt.db.Table(mt.tableName).Create(toInsertLines).Error; err != nil {
		return err
	}
	return nil
}

// Update 三个参数：queryStr、 field、target
// 调用此函数，相当于针对queryStr 筛选出的数据条目(例如queryStr = 'userId = ?', args = "1") ，将筛选出的数据条目的field字段替换为target内容
func (mt *Impl) Update(queryStr, field string, target interface{}, args ...interface{}) error {
	if err := mt.db.Table(mt.tableName).Where(queryStr, args...).Update(field, target).Error; err != nil {
		return err
	}
	return nil
}

// Delete 一个参数：toDeleteTarget
// 传入一个UserDefinedModel类型，如果此对象的userId = 1,则删除掉数据库中userId= 1的字段
func (mt *Impl) Delete(toDeleteTarget UserDefinedModel) error {
	if err := mt.db.Table(mt.tableName).Delete(toDeleteTarget).Error; err != nil {
		return err
	}
	return nil
}

// First 两个参数： queryStr、findTarget
// queryStr为筛选用的query，例如`userId = ?`, args = "1", findTarget为 UserDefinedModel 类型指针，为第一个找到的数据。
func (mt *Impl) First(queryStr string, findTarget UserDefinedModel, args ...interface{}) error {
	if err := mt.db.Table(mt.tableName).Where(queryStr, args...).Find(findTarget).Error; err != nil {
		return err
	}
	return nil
}

var _ Mysql = &Impl{}

func GetMysql(config *Config) (Mysql, error) {
	mysqlImpl, err := normal.GetImpl(SDID, config)
	if err != nil {
		return nil, err
	}
	return mysqlImpl.(Mysql), nil
}
