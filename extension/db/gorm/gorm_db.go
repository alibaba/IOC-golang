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

package gorm

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/alibaba/ioc-golang/autowire"
)

// +ioc:autowire=true
// +ioc:autowire:type=normal
// +ioc:autowire:type=singleton
// +ioc:autowire:paramType=Param
// +ioc:autowire:constructFunc=New

type GORMDB struct {
	db *gorm.DB
}

func fromDB(db *gorm.DB) GORMDBIOCInterface {
	return autowire.GetProxyFunction()(&GORMDB{
		db: db,
	}).(GORMDBIOCInterface)
}

// Session create new db session
func (db *GORMDB) Session(config *gorm.Session) GORMDBIOCInterface {
	return fromDB(db.db.Session(config))
}

// WithContext change current instance db's context to ctx
func (db *GORMDB) WithContext(ctx context.Context) GORMDBIOCInterface {
	return fromDB(db.db.WithContext(ctx))
}

// Debug start debug mode
func (db *GORMDB) Debug() GORMDBIOCInterface {
	return fromDB(db.db.Debug())
}

// Set store value with key into current db instance's context
func (db *GORMDB) Set(key string, value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Set(key, value))
}

// Get get value with key from current db instance's context
func (db *GORMDB) Get(key string) (interface{}, bool) {
	return db.db.Get(key)
}

// InstanceSet store value with key into current db instance's context
func (db *GORMDB) InstanceSet(key string, value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.InstanceSet(key, value))
}

// InstanceGet get value with key from current db instance's context
func (db *GORMDB) InstanceGet(key string) (interface{}, bool) {
	return db.db.InstanceGet(key)
}

// AddError add error to db
func (db *GORMDB) AddError(err error) error {
	return db.db.AddError(err)
}

// DB returns `*sql.DB`
func (db *GORMDB) DB() (*sql.DB, error) {
	return db.db.DB()
}

// SetupJoinTable setup join table schema
func (db *GORMDB) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return db.db.SetupJoinTable(model, field, joinTable)
}

// Use use plugin
func (db *GORMDB) Use(plugin gorm.Plugin) error {
	return db.db.Use(plugin)
}

func (db *GORMDB) ToSQL(queryFn func(tx *gorm.DB) *gorm.DB) string {
	return db.db.ToSQL(queryFn)
}

// Model specify the model you would like to run db operations
//    // update all users's name to `hello`
//    db.Model(&User{}).Update("name", "hello")
//    // if user's primary key is non-blank, will use it as condition, then will only update the user's name to `hello`
//    db.Model(&user).Update("name", "hello")
func (db *GORMDB) Model(value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Model(value))
}

// Clauses Add clauses
func (db *GORMDB) Clauses(conds ...clause.Expression) GORMDBIOCInterface {
	return fromDB(db.db.Clauses(conds...))
}

// Table specify the table you would like to run db operations
func (db *GORMDB) Table(name string, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Table(name, args...))
}

// Distinct specify distinct fields that you want querying
func (db *GORMDB) Distinct(args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Distinct(args...))
}

// Select specify fields that you want when querying, creating, updating
func (db *GORMDB) Select(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Select(query, args...))
}

// Omit specify fields that you want to ignore when creating, updating and querying
func (db *GORMDB) Omit(columns ...string) GORMDBIOCInterface {
	return fromDB(db.db.Omit(columns...))
}

// Where add conditions
func (db *GORMDB) Where(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Where(query, args...))
}

// Not add NOT conditions
func (db *GORMDB) Not(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Not(query, args...))
}

// Or add OR conditions
func (db *GORMDB) Or(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Or(query, args...))
}

// Joins specify Joins conditions
//     db.Joins("Account").Find(&user)
//     db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Find(&user)
//     db.Joins("Account", DB.Select("id").Where("user_id = users.id AND name = ?", "someName").Model(&Account{}))
func (db *GORMDB) Joins(query string, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Joins(query, args...))
}

// Group specify the group method on the find
func (db *GORMDB) Group(name string) GORMDBIOCInterface {
	return fromDB(db.db.Group(name))
}

// Having specify HAVING conditions for GROUP BY
func (db *GORMDB) Having(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Having(query, args...))
}

// Order specify order when retrieve records from database
//     db.Order("name DESC")
//     db.Order(clause.OrderByColumn{Column: clause.Column{Name: "name"}, Desc: true})
func (db *GORMDB) Order(value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Order(value))
}

// Limit specify the number of records to be retrieved
func (db *GORMDB) Limit(limit int) GORMDBIOCInterface {
	return fromDB(db.db.Limit(limit))
}

// Offset specify the number of records to skip before starting to return the records
func (db *GORMDB) Offset(offset int) GORMDBIOCInterface {
	return fromDB(db.db.Offset(offset))
}

// Scopes pass current database connection to arguments `func(DB) DB`, which could be used to add conditions dynamically
//     func AmountGreaterThan1000(db *gorm.DB) *gorm.DB {
//         return db.Where("amount > ?", 1000)
//     }
//
//     func OrderStatus(status []string) func (db *gorm.DB) *gorm.DB {
//         return func (db *gorm.DB) *gorm.DB {
//             return db.Scopes(AmountGreaterThan1000).Where("status in (?)", status)
//         }
//     }
//
//     db.Scopes(AmountGreaterThan1000, OrderStatus([]string{"paid", "shipped"})).Find(&orders)
func (db *GORMDB) Scopes(funcs ...func(db *gorm.DB) *gorm.DB) GORMDBIOCInterface {
	return fromDB(db.db.Scopes(funcs...))
}

// Preload preload associations with given conditions
//    db.Preload("Orders", "state NOT IN (?)", "cancelled").Find(&users)
func (db *GORMDB) Preload(query string, args ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Preload(query, args...))
}

func (db *GORMDB) Attrs(attrs ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Attrs(attrs...))
}

func (db *GORMDB) Assign(attrs ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Assign(attrs...))
}

func (db *GORMDB) Unscoped() GORMDBIOCInterface {
	return fromDB(db.db.Unscoped())
}

func (db *GORMDB) Raw(sql string, values ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Raw(sql, values...))
}

func (db *GORMDB) Error() error {
	return db.db.Error
}

// Create insert the value into database
func (db *GORMDB) Create(value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Create(value))
}

// CreateInBatches insert the value in batches into database
func (db *GORMDB) CreateInBatches(value interface{}, batchSize int) GORMDBIOCInterface {
	return fromDB(db.db.CreateInBatches(value, batchSize))
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (db *GORMDB) Save(value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Save(value))
}

// First find first record that match given conditions, order by primary key
func (db *GORMDB) First(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.First(dest, conds...))
}

// Take return a record that match given conditions, the order will depend on the database implementation
func (db *GORMDB) Take(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Take(dest, conds...))
}

// Last find last record that match given conditions, order by primary key
func (db *GORMDB) Last(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Last(dest, conds...))
}

// Find find records that match given conditions
func (db *GORMDB) Find(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Find(dest, conds...))
}

// FindInBatches find records in batches
func (db *GORMDB) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) GORMDBIOCInterface {
	return fromDB(db.db.FindInBatches(dest, batchSize, fc))
}

// FirstOrInit gets the first matched record or initialize a new instance with given conditions (only works with struct or map conditions)
func (db *GORMDB) FirstOrInit(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.FirstOrInit(dest, conds...))
}

// FirstOrCreate gets the first matched record or create a new one with given conditions (only works with struct, map conditions)
func (db *GORMDB) FirstOrCreate(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.FirstOrCreate(dest, conds...))
}

// Update update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
func (db *GORMDB) Update(column string, value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Update(column, value))
}

// Updates update attributes with callbacks, refer: https://gorm.io/docs/update.html#Update-Changed-Fields
func (db *GORMDB) Updates(values interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Updates(values))
}

func (db *GORMDB) UpdateColumn(column string, value interface{}) GORMDBIOCInterface {
	return fromDB(db.db.UpdateColumn(column, value))
}

func (db *GORMDB) UpdateColumns(values interface{}) GORMDBIOCInterface {
	return fromDB(db.db.UpdateColumns(values))
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
func (db *GORMDB) Delete(value interface{}, conds ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Delete(value, conds...))
}

func (db *GORMDB) Count(count *int64) GORMDBIOCInterface {
	return fromDB(db.db.Count(count))
}

func (db *GORMDB) Row() *sql.Row {
	return db.db.Row()
}

func (db *GORMDB) Rows() (*sql.Rows, error) {
	return db.db.Rows()
}

// Scan scan value to a struct
func (db *GORMDB) Scan(dest interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Scan(dest))
}

// Pluck used to query single column from a model as a map
//     var ages []int64
//     db.Model(&users).Pluck("age", &ages)
func (db *GORMDB) Pluck(column string, dest interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Pluck(column, dest))
}

func (db *GORMDB) ScanRows(rows *sql.Rows, dest interface{}) error {
	return db.db.ScanRows(rows, dest)
}

// Connection  use a db conn to execute Multiple commands,this conn will put conn pool after it is executed.
func (db *GORMDB) Connection(fc func(db *gorm.DB) error) (err error) {
	return db.db.Connection(fc)
}

// Transaction start a transaction as a block, return error will rollback, otherwise to commit.
func (db *GORMDB) Transaction(fc func(db *gorm.DB) error, opts ...*sql.TxOptions) (err error) {
	return db.db.Transaction(fc, opts...)
}

// Begin begins a transaction
func (db *GORMDB) Begin(opts ...*sql.TxOptions) GORMDBIOCInterface {
	return fromDB(db.db.Begin(opts...))
}

// Commit commit a transaction
func (db *GORMDB) Commit() GORMDBIOCInterface {
	return fromDB(db.db.Commit())
}

// Rollback rollback a transaction
func (db *GORMDB) Rollback() {
	db.db.Rollback()
}

func (db *GORMDB) SavePoint(name string) GORMDBIOCInterface {
	return fromDB(db.db.SavePoint(name))
}

func (db *GORMDB) RollbackTo(name string) GORMDBIOCInterface {
	return fromDB(db.db.RollbackTo(name))
}

// Exec execute raw sql
func (db *GORMDB) Exec(sql string, values ...interface{}) GORMDBIOCInterface {
	return fromDB(db.db.Exec(sql, values...))
}

func (db *GORMDB) Migrator() gorm.Migrator {
	return db.db.Migrator()
}

func (db *GORMDB) AutoMigrate(dst ...interface{}) error {
	return db.db.AutoMigrate(dst...)
}

func (db *GORMDB) Association(column string) *gorm.Association {
	return db.db.Association(column)
}
