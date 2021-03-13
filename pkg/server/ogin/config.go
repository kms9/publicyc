package ogin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kms9/publicyc/pkg/conf"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
)

// Config HTTP conf
type Config struct {
	Name    string
	Host    string
	Port    int
	Mode    string
	_logger *onion_log.Log
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Name:    "gin",
		Host:    "",
		Port:    10086,
		Mode:    gin.ReleaseMode,
		_logger: onion_log.DefaultLogger(),
	}
}

// UseConfig 获取gin的配置参数
func UseConfig(name string) *Config {
	return RawConfig("yc.server." + name)
}

// RawConfig 查询配置是否有配置
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Detail().Get(key) == nil {
		config._logger.Warnf("config key: %s is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, &config); err != nil {
		config._logger.Panicf("http server parse conf panic key:%s err:%s", key, err)
	}
	return config
}

// Build create server instance, then initialize it with necessary interceptor
func (c *Config) Build() *Server {
	server := newServer(c)
	server.Use(
		ErrorTrace(c._logger),
		LogMiddle(c._logger, c.Name),
	)

	return server
}

// Address 地址
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
