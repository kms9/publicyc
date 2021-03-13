package opg

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Handler ...
type Handler func(*Scope)

// Interceptor ...
type Interceptor func(string, *Config) func(next Handler) Handler

func metricInterceptor(op string, c *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.HasError() {
				if scope.DB().Error != gorm.ErrRecordNotFound {
					c._logger.Errorf("mysql err: ", scope.DB().Error)
				} else {
					// c._logger.Warnf("not record")
				}
			} else {
				c._logger.Debug("OK")
			}

			if c.SlowThreshold > time.Duration(0) && c.SlowThreshold < cost {
				c._logger.With(nil).WithFields(map[string]interface{}{
					"sql":     scope.SQL,
					"sqlVal":  scope.SQLVars,
					"table":   scope.TableName(),
					"latency": cost / time.Millisecond,
				}).Warn("slowSql")
			}
		}
	}
}
