//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package gorm

import (
	contextx "context"
	"database/sql"

	gorm_iogormx "gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &gORMDB_{}
		},
	})
	gORMDBStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &GORMDB{}
		},
		ParamFactory: func() interface{} {
			var _ paramInterface = &Param{}
			return &Param{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(paramInterface)
			impl := i.(*GORMDB)
			return param.New(impl)
		},
		Metadata: map[string]interface{}{
			"aop":      map[string]interface{}{},
			"autowire": map[string]interface{}{},
		},
	}
	normal.RegisterStructDescriptor(gORMDBStructDescriptor)
	singleton.RegisterStructDescriptor(gORMDBStructDescriptor)
}

type paramInterface interface {
	New(impl *GORMDB) (*GORMDB, error)
}
type gORMDB_ struct {
	Session_         func(config *gorm_iogormx.Session) GORMDBIOCInterface
	WithContext_     func(ctx contextx.Context) GORMDBIOCInterface
	Debug_           func() GORMDBIOCInterface
	Set_             func(key string, value interface{}) GORMDBIOCInterface
	Get_             func(key string) (interface{}, bool)
	InstanceSet_     func(key string, value interface{}) GORMDBIOCInterface
	InstanceGet_     func(key string) (interface{}, bool)
	AddError_        func(err error) error
	DB_              func() (*sql.DB, error)
	SetupJoinTable_  func(model interface{}, field string, joinTable interface{}) error
	Use_             func(plugin gorm_iogormx.Plugin) error
	ToSQL_           func(queryFn func(tx *gorm_iogormx.DB) *gorm_iogormx.DB) string
	Model_           func(value interface{}) GORMDBIOCInterface
	Clauses_         func(conds ...clause.Expression) GORMDBIOCInterface
	Table_           func(name string, args ...interface{}) GORMDBIOCInterface
	Distinct_        func(args ...interface{}) GORMDBIOCInterface
	Select_          func(query interface{}, args ...interface{}) GORMDBIOCInterface
	Omit_            func(columns ...string) GORMDBIOCInterface
	Where_           func(query interface{}, args ...interface{}) GORMDBIOCInterface
	Not_             func(query interface{}, args ...interface{}) GORMDBIOCInterface
	Or_              func(query interface{}, args ...interface{}) GORMDBIOCInterface
	Joins_           func(query string, args ...interface{}) GORMDBIOCInterface
	Group_           func(name string) GORMDBIOCInterface
	Having_          func(query interface{}, args ...interface{}) GORMDBIOCInterface
	Order_           func(value interface{}) GORMDBIOCInterface
	Limit_           func(limit int) GORMDBIOCInterface
	Offset_          func(offset int) GORMDBIOCInterface
	Scopes_          func(funcs ...func(db *gorm_iogormx.DB) *gorm_iogormx.DB) GORMDBIOCInterface
	Preload_         func(query string, args ...interface{}) GORMDBIOCInterface
	Attrs_           func(attrs ...interface{}) GORMDBIOCInterface
	Assign_          func(attrs ...interface{}) GORMDBIOCInterface
	Unscoped_        func() GORMDBIOCInterface
	Raw_             func(sql string, values ...interface{}) GORMDBIOCInterface
	Error_           func() error
	Create_          func(value interface{}) GORMDBIOCInterface
	CreateInBatches_ func(value interface{}, batchSize int) GORMDBIOCInterface
	Save_            func(value interface{}) GORMDBIOCInterface
	First_           func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Take_            func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Last_            func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Find_            func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	FindInBatches_   func(dest interface{}, batchSize int, fc func(tx *gorm_iogormx.DB, batch int) error) GORMDBIOCInterface
	FirstOrInit_     func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	FirstOrCreate_   func(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Update_          func(column string, value interface{}) GORMDBIOCInterface
	Updates_         func(values interface{}) GORMDBIOCInterface
	UpdateColumn_    func(column string, value interface{}) GORMDBIOCInterface
	UpdateColumns_   func(values interface{}) GORMDBIOCInterface
	Delete_          func(value interface{}, conds ...interface{}) GORMDBIOCInterface
	Count_           func(count *int64) GORMDBIOCInterface
	Row_             func() *sql.Row
	Rows_            func() (*sql.Rows, error)
	Scan_            func(dest interface{}) GORMDBIOCInterface
	Pluck_           func(column string, dest interface{}) GORMDBIOCInterface
	ScanRows_        func(rows *sql.Rows, dest interface{}) error
	Connection_      func(fc func(db *gorm_iogormx.DB) error) (err error)
	Transaction_     func(fc func(db *gorm_iogormx.DB) error, opts ...*sql.TxOptions) (err error)
	Begin_           func(opts ...*sql.TxOptions) GORMDBIOCInterface
	Commit_          func() GORMDBIOCInterface
	Rollback_        func()
	SavePoint_       func(name string) GORMDBIOCInterface
	RollbackTo_      func(name string) GORMDBIOCInterface
	Exec_            func(sql string, values ...interface{}) GORMDBIOCInterface
	Migrator_        func() gorm_iogormx.Migrator
	AutoMigrate_     func(dst ...interface{}) error
	Association_     func(column string) *gorm_iogormx.Association
}

func (g *gORMDB_) Session(config *gorm_iogormx.Session) GORMDBIOCInterface {
	return g.Session_(config)
}

func (g *gORMDB_) WithContext(ctx contextx.Context) GORMDBIOCInterface {
	return g.WithContext_(ctx)
}

func (g *gORMDB_) Debug() GORMDBIOCInterface {
	return g.Debug_()
}

func (g *gORMDB_) Set(key string, value interface{}) GORMDBIOCInterface {
	return g.Set_(key, value)
}

func (g *gORMDB_) Get(key string) (interface{}, bool) {
	return g.Get_(key)
}

func (g *gORMDB_) InstanceSet(key string, value interface{}) GORMDBIOCInterface {
	return g.InstanceSet_(key, value)
}

func (g *gORMDB_) InstanceGet(key string) (interface{}, bool) {
	return g.InstanceGet_(key)
}

func (g *gORMDB_) AddError(err error) error {
	return g.AddError_(err)
}

func (g *gORMDB_) DB() (*sql.DB, error) {
	return g.DB_()
}

func (g *gORMDB_) SetupJoinTable(model interface{}, field string, joinTable interface{}) error {
	return g.SetupJoinTable_(model, field, joinTable)
}

func (g *gORMDB_) Use(plugin gorm_iogormx.Plugin) error {
	return g.Use_(plugin)
}

func (g *gORMDB_) ToSQL(queryFn func(tx *gorm_iogormx.DB) *gorm_iogormx.DB) string {
	return g.ToSQL_(queryFn)
}

func (g *gORMDB_) Model(value interface{}) GORMDBIOCInterface {
	return g.Model_(value)
}

func (g *gORMDB_) Clauses(conds ...clause.Expression) GORMDBIOCInterface {
	return g.Clauses_(conds...)
}

func (g *gORMDB_) Table(name string, args ...interface{}) GORMDBIOCInterface {
	return g.Table_(name, args...)
}

func (g *gORMDB_) Distinct(args ...interface{}) GORMDBIOCInterface {
	return g.Distinct_(args...)
}

func (g *gORMDB_) Select(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return g.Select_(query, args...)
}

func (g *gORMDB_) Omit(columns ...string) GORMDBIOCInterface {
	return g.Omit_(columns...)
}

func (g *gORMDB_) Where(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return g.Where_(query, args...)
}

func (g *gORMDB_) Not(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return g.Not_(query, args...)
}

func (g *gORMDB_) Or(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return g.Or_(query, args...)
}

func (g *gORMDB_) Joins(query string, args ...interface{}) GORMDBIOCInterface {
	return g.Joins_(query, args...)
}

func (g *gORMDB_) Group(name string) GORMDBIOCInterface {
	return g.Group_(name)
}

func (g *gORMDB_) Having(query interface{}, args ...interface{}) GORMDBIOCInterface {
	return g.Having_(query, args...)
}

func (g *gORMDB_) Order(value interface{}) GORMDBIOCInterface {
	return g.Order_(value)
}

func (g *gORMDB_) Limit(limit int) GORMDBIOCInterface {
	return g.Limit_(limit)
}

func (g *gORMDB_) Offset(offset int) GORMDBIOCInterface {
	return g.Offset_(offset)
}

func (g *gORMDB_) Scopes(funcs ...func(db *gorm_iogormx.DB) *gorm_iogormx.DB) GORMDBIOCInterface {
	return g.Scopes_(funcs...)
}

func (g *gORMDB_) Preload(query string, args ...interface{}) GORMDBIOCInterface {
	return g.Preload_(query, args...)
}

func (g *gORMDB_) Attrs(attrs ...interface{}) GORMDBIOCInterface {
	return g.Attrs_(attrs...)
}

func (g *gORMDB_) Assign(attrs ...interface{}) GORMDBIOCInterface {
	return g.Assign_(attrs...)
}

func (g *gORMDB_) Unscoped() GORMDBIOCInterface {
	return g.Unscoped_()
}

func (g *gORMDB_) Raw(sql string, values ...interface{}) GORMDBIOCInterface {
	return g.Raw_(sql, values...)
}

func (g *gORMDB_) Error() error {
	return g.Error_()
}

func (g *gORMDB_) Create(value interface{}) GORMDBIOCInterface {
	return g.Create_(value)
}

func (g *gORMDB_) CreateInBatches(value interface{}, batchSize int) GORMDBIOCInterface {
	return g.CreateInBatches_(value, batchSize)
}

func (g *gORMDB_) Save(value interface{}) GORMDBIOCInterface {
	return g.Save_(value)
}

func (g *gORMDB_) First(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.First_(dest, conds...)
}

func (g *gORMDB_) Take(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.Take_(dest, conds...)
}

func (g *gORMDB_) Last(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.Last_(dest, conds...)
}

func (g *gORMDB_) Find(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.Find_(dest, conds...)
}

func (g *gORMDB_) FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm_iogormx.DB, batch int) error) GORMDBIOCInterface {
	return g.FindInBatches_(dest, batchSize, fc)
}

func (g *gORMDB_) FirstOrInit(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.FirstOrInit_(dest, conds...)
}

func (g *gORMDB_) FirstOrCreate(dest interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.FirstOrCreate_(dest, conds...)
}

func (g *gORMDB_) Update(column string, value interface{}) GORMDBIOCInterface {
	return g.Update_(column, value)
}

func (g *gORMDB_) Updates(values interface{}) GORMDBIOCInterface {
	return g.Updates_(values)
}

func (g *gORMDB_) UpdateColumn(column string, value interface{}) GORMDBIOCInterface {
	return g.UpdateColumn_(column, value)
}

func (g *gORMDB_) UpdateColumns(values interface{}) GORMDBIOCInterface {
	return g.UpdateColumns_(values)
}

func (g *gORMDB_) Delete(value interface{}, conds ...interface{}) GORMDBIOCInterface {
	return g.Delete_(value, conds...)
}

func (g *gORMDB_) Count(count *int64) GORMDBIOCInterface {
	return g.Count_(count)
}

func (g *gORMDB_) Row() *sql.Row {
	return g.Row_()
}

func (g *gORMDB_) Rows() (*sql.Rows, error) {
	return g.Rows_()
}

func (g *gORMDB_) Scan(dest interface{}) GORMDBIOCInterface {
	return g.Scan_(dest)
}

func (g *gORMDB_) Pluck(column string, dest interface{}) GORMDBIOCInterface {
	return g.Pluck_(column, dest)
}

func (g *gORMDB_) ScanRows(rows *sql.Rows, dest interface{}) error {
	return g.ScanRows_(rows, dest)
}

func (g *gORMDB_) Connection(fc func(db *gorm_iogormx.DB) error) (err error) {
	return g.Connection_(fc)
}

func (g *gORMDB_) Transaction(fc func(db *gorm_iogormx.DB) error, opts ...*sql.TxOptions) (err error) {
	return g.Transaction_(fc, opts...)
}

func (g *gORMDB_) Begin(opts ...*sql.TxOptions) GORMDBIOCInterface {
	return g.Begin_(opts...)
}

func (g *gORMDB_) Commit() GORMDBIOCInterface {
	return g.Commit_()
}

func (g *gORMDB_) Rollback() {
	g.Rollback_()
}

func (g *gORMDB_) SavePoint(name string) GORMDBIOCInterface {
	return g.SavePoint_(name)
}

func (g *gORMDB_) RollbackTo(name string) GORMDBIOCInterface {
	return g.RollbackTo_(name)
}

func (g *gORMDB_) Exec(sql string, values ...interface{}) GORMDBIOCInterface {
	return g.Exec_(sql, values...)
}

func (g *gORMDB_) Migrator() gorm_iogormx.Migrator {
	return g.Migrator_()
}

func (g *gORMDB_) AutoMigrate(dst ...interface{}) error {
	return g.AutoMigrate_(dst...)
}

func (g *gORMDB_) Association(column string) *gorm_iogormx.Association {
	return g.Association_(column)
}

type GORMDBIOCInterface interface {
	Session(config *gorm_iogormx.Session) GORMDBIOCInterface
	WithContext(ctx contextx.Context) GORMDBIOCInterface
	Debug() GORMDBIOCInterface
	Set(key string, value interface{}) GORMDBIOCInterface
	Get(key string) (interface{}, bool)
	InstanceSet(key string, value interface{}) GORMDBIOCInterface
	InstanceGet(key string) (interface{}, bool)
	AddError(err error) error
	DB() (*sql.DB, error)
	SetupJoinTable(model interface{}, field string, joinTable interface{}) error
	Use(plugin gorm_iogormx.Plugin) error
	ToSQL(queryFn func(tx *gorm_iogormx.DB) *gorm_iogormx.DB) string
	Model(value interface{}) GORMDBIOCInterface
	Clauses(conds ...clause.Expression) GORMDBIOCInterface
	Table(name string, args ...interface{}) GORMDBIOCInterface
	Distinct(args ...interface{}) GORMDBIOCInterface
	Select(query interface{}, args ...interface{}) GORMDBIOCInterface
	Omit(columns ...string) GORMDBIOCInterface
	Where(query interface{}, args ...interface{}) GORMDBIOCInterface
	Not(query interface{}, args ...interface{}) GORMDBIOCInterface
	Or(query interface{}, args ...interface{}) GORMDBIOCInterface
	Joins(query string, args ...interface{}) GORMDBIOCInterface
	Group(name string) GORMDBIOCInterface
	Having(query interface{}, args ...interface{}) GORMDBIOCInterface
	Order(value interface{}) GORMDBIOCInterface
	Limit(limit int) GORMDBIOCInterface
	Offset(offset int) GORMDBIOCInterface
	Scopes(funcs ...func(db *gorm_iogormx.DB) *gorm_iogormx.DB) GORMDBIOCInterface
	Preload(query string, args ...interface{}) GORMDBIOCInterface
	Attrs(attrs ...interface{}) GORMDBIOCInterface
	Assign(attrs ...interface{}) GORMDBIOCInterface
	Unscoped() GORMDBIOCInterface
	Raw(sql string, values ...interface{}) GORMDBIOCInterface
	Error() error
	Create(value interface{}) GORMDBIOCInterface
	CreateInBatches(value interface{}, batchSize int) GORMDBIOCInterface
	Save(value interface{}) GORMDBIOCInterface
	First(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Take(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Last(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Find(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	FindInBatches(dest interface{}, batchSize int, fc func(tx *gorm_iogormx.DB, batch int) error) GORMDBIOCInterface
	FirstOrInit(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	FirstOrCreate(dest interface{}, conds ...interface{}) GORMDBIOCInterface
	Update(column string, value interface{}) GORMDBIOCInterface
	Updates(values interface{}) GORMDBIOCInterface
	UpdateColumn(column string, value interface{}) GORMDBIOCInterface
	UpdateColumns(values interface{}) GORMDBIOCInterface
	Delete(value interface{}, conds ...interface{}) GORMDBIOCInterface
	Count(count *int64) GORMDBIOCInterface
	Row() *sql.Row
	Rows() (*sql.Rows, error)
	Scan(dest interface{}) GORMDBIOCInterface
	Pluck(column string, dest interface{}) GORMDBIOCInterface
	ScanRows(rows *sql.Rows, dest interface{}) error
	Connection(fc func(db *gorm_iogormx.DB) error) (err error)
	Transaction(fc func(db *gorm_iogormx.DB) error, opts ...*sql.TxOptions) (err error)
	Begin(opts ...*sql.TxOptions) GORMDBIOCInterface
	Commit() GORMDBIOCInterface
	Rollback()
	SavePoint(name string) GORMDBIOCInterface
	RollbackTo(name string) GORMDBIOCInterface
	Exec(sql string, values ...interface{}) GORMDBIOCInterface
	Migrator() gorm_iogormx.Migrator
	AutoMigrate(dst ...interface{}) error
	Association(column string) *gorm_iogormx.Association
}

var _gORMDBSDID string

func GetGORMDB(p *Param) (*GORMDB, error) {
	if _gORMDBSDID == "" {
		_gORMDBSDID = util.GetSDIDByStructPtr(new(GORMDB))
	}
	i, err := normal.GetImpl(_gORMDBSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*GORMDB)
	return impl, nil
}

func GetGORMDBIOCInterface(p *Param) (GORMDBIOCInterface, error) {
	if _gORMDBSDID == "" {
		_gORMDBSDID = util.GetSDIDByStructPtr(new(GORMDB))
	}
	i, err := normal.GetImplWithProxy(_gORMDBSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(GORMDBIOCInterface)
	return impl, nil
}

func GetGORMDBSingleton(p *Param) (*GORMDB, error) {
	if _gORMDBSDID == "" {
		_gORMDBSDID = util.GetSDIDByStructPtr(new(GORMDB))
	}
	i, err := singleton.GetImpl(_gORMDBSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*GORMDB)
	return impl, nil
}

func GetGORMDBIOCInterfaceSingleton(p *Param) (GORMDBIOCInterface, error) {
	if _gORMDBSDID == "" {
		_gORMDBSDID = util.GetSDIDByStructPtr(new(GORMDB))
	}
	i, err := singleton.GetImplWithProxy(_gORMDBSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(GORMDBIOCInterface)
	return impl, nil
}

type ThisGORMDB struct {
}

func (t *ThisGORMDB) This() GORMDBIOCInterface {
	thisPtr, _ := GetGORMDBIOCInterfaceSingleton(nil)
	return thisPtr
}
