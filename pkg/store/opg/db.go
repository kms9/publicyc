package opg

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kms9/publicyc/pkg/util"
)

type (
	DB = gorm.DB
	// Scope ...
	Scope = gorm.Scope
)

// Open 连接数据库
func Open(c *Config) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s sslmode=%s dbname=%s",
		c.Host, c.Port, c.User, c.Password, c.SslMode, c.Db,
	)
	db, err := gorm.Open("postgres", connStr)
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

	c._logger.Infof("pg: %s connect: %s MaxIdleConns=%d MaxOpenConns=%d ConnMaxLifetime=%ds", c.Name, connStr, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxLifetime/time.Second)
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

	if len(c.UpdateColumn) > 0 {
		// 更新时间
		db.Callback().Update().Replace("gorm:update_time_stamp", func(scope *Scope) {
			if _, ok := scope.Get("gorm:update_column"); !ok {
				now := time.Now()
				for _, v := range c.UpdateColumn {
					err = scope.SetColumn(v, &now)
				}
			}
		})
	}
	return db, nil
}
