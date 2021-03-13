package omysql

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kms9/publicyc/pkg/util"
)

type (
	DB = gorm.DB
	// Scope ...
	Scope = gorm.Scope
)

// Open 连接数据库
func Open(c *Config) (*DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", c.User,  c.Password, c.Host, c.Port, c.Db)
	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if util.GetGoEnv() == "development" {
		db.LogMode(true)
	} else {
		db.LogMode(c.Debug)
	}

	db.SingularTable(true)
	db.DB().SetMaxIdleConns(c.MaxIdleConns)
	db.DB().SetMaxOpenConns(c.MaxOpenConns)
	if c.ConnMaxLifetime != 0 {
		db.DB().SetConnMaxLifetime(c.ConnMaxLifetime)
	}

	c._logger.Infof("mysql: %s connect: %s MaxIdleConns=%d MaxOpenConns=%d ConnMaxLifetime=%ds", c.Name, dsn, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxLifetime/time.Second)
	replace := func(processor func() *gorm.CallbackProcessor, callbackName string, interceptors ...Interceptor) {
		old := processor().Get(callbackName)
		var handler = old
		for _, inte := range interceptors {
			handler = inte(callbackName, c)(handler)
		}
		processor().Replace(callbackName, handler)
	}

	replace(
		db.Callback().Delete,
		"gorm:delete",
		c.interceptors...,
	)
	replace(
		db.Callback().Update,
		"gorm:update",
		c.interceptors...,
	)
	replace(
		db.Callback().Create,
		"gorm:create",
		c.interceptors...,
	)
	replace(
		db.Callback().Query,
		"gorm:query",
		c.interceptors...,
	)
	replace(
		db.Callback().RowQuery,
		"gorm:row_query",
		c.interceptors...,
	)

	return db, nil
}
