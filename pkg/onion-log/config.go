package onion_log

import (
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/kms9/publicyc/pkg/conf"
	"github.com/kms9/publicyc/pkg/onion-log/hook"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
	"github.com/kms9/publicyc/pkg/util"
)

// Config 日志配置信息
type Config struct {
	Level string // Level 日志初始等级
	Hooks map[string]interface{}
}

// RawConfig 获取配置信息
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if err := conf.Detail().UnmarshalKey(key, &config); err != nil {
		panic(err)
	}
	return config
}

// UseConfig 使用日志的配置文件
func UseConfig(name string) *Config {
	return RawConfig("yc." + name)
}

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		Level: "info",
	}
}

// Build 获取配置
func (c *Config) Build() *Log {
	var log *Log
	var hooks []logrus.Hook
	if len(c.Hooks) > 0 {
		for k, v := range c.Hooks {
			switch k {
			case "ding":
				ding := v.(map[string]interface{})
				if ding["sign"] != nil && ding["access_token"] != nil {
					hooks = append(hooks, &hook.NotifyHook{
						Ding: &hook.Ding{
							Sign:        ding["sign"].(string),
							AccessToken: ding["access_token"].(string),
						},
					})
				}
			}
		}
	}
	log = New(c.Level, util.GetGoEnv(), hooks...)
	_logger = log.With(&logger.BaseContentInfo{UID: "yc", SpanID: uuid.NewV4().String()})
	return log
}
