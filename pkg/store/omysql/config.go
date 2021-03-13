package omysql

import (
	"time"

	"github.com/kms9/publicyc/pkg/conf"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/util/otime"

	//"github.com/kms9/publicyc/pkg/util/otime"
)

// config options
type Config struct {
	Name string
	// 连接地址
	Host string `json:"host"`
	// 端口
	Port int `json:"port"`
	// 用户名
	User string `json:"user"`
	// 密码
	Password string `json:"password"`
	// 数据库名称
	Db string `json:"db"`

	// 最大空闲连接数
	MaxIdleConns int `json:"maxIdleConns"`
	// 最大活动连接数
	MaxOpenConns int `json:"maxOpenConns"`
	// 连接的最大存活时间
	ConnMaxLifetime time.Duration `json:"connMaxLifetime"`
	// 慢日志阈值
	SlowThreshold time.Duration `json:"slowThreshold"`
	// 关闭指标采集
	DisableMetric bool `json:"disableMetric"`
	// Debug开关
	Debug         bool `json:"debug"`

	// 日志
	_logger *onion_log.Log
	// sql回调函数
	interceptors []Interceptor

}

// UseConfig 标准配置，规范配置文件头
func UseConfig(name string) *Config {
	return RawConfig("yc.mysql." + name)
}

// RawConfig 解析配置
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Detail().Get(key) == nil {
		onion_log.Panicf("key:%s mysqlConfig is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, config); err != nil {
		onion_log.Panicf("unmarshal mysqlConfig key:%s err:%s", key, err)
	}

	config.Name = key
	return config
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Debug:           false,
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: otime.Duration("300s"),
		SlowThreshold:   otime.Duration("500ms"),
		_logger:         onion_log.DefaultLogger(),
	}
}

// WithInterceptor ...
func (c *Config) WithInterceptor(intes ...Interceptor) *Config {
	if c.interceptors == nil {
		c.interceptors = make([]Interceptor, 0)
	}
	c.interceptors = append(c.interceptors, intes...)
	return c
}

// Build ...
func (c *Config) Build() *DB {
	if !c.DisableMetric {
		c = c.WithInterceptor(metricInterceptor)
	}
	db, err := Open(c)
	if err != nil {
		c._logger.Panicf("connect mysql err %s", err)
	}

	if err := db.DB().Ping(); err != nil {
		c._logger.Panicf("ping mysql err %s", err)
	}
	return db
}
