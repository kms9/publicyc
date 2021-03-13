package ohttp

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/kms9/publicyc/pkg/conf"
	onion_log "github.com/kms9/publicyc/pkg/onion-log"
	"github.com/kms9/publicyc/pkg/onion-log/logger"
)

// Config url配置信息
type RequestConfig struct {
	TimeOut     time.Duration `json:"timeout"`
	Path        string        `json:"path"`
	Method      string        `json:"method"`
	Header      []string      `json:"header"`
	Query       []string      `json:"query"`
	PathParams  []string      `json:"pathParams"`
	BodyRequire bool          `json:"bodyRequire"`
}

// Curl 请求
func (c *RequestConfig) CurlWithContext(trace *logger.BaseContentInfo, url string, params map[string]interface{}) ([]byte, int, error) {
	response, err := request(trace, resty.New(), url, c, params)
	if err != nil {
		return nil, 0, err
	}
	return response.Body(), response.StatusCode(), nil
}

type ProjectConfig struct {
	URL     string                    `json:"url"`
	TimeOut time.Duration             `json:"timeout"`
	API     map[string]*RequestConfig `json:"api"`
	_logger *onion_log.Log
}

type Projects map[string]*ProjectConfig

type Config struct {
	Projects    Projects      `json:"projects"`
	SlowRequest time.Duration `json:"slowRequest"`
	Debug       bool          `json:"debug"`
	_logger     *onion_log.Log
}

var defaultSetTimeout = 3 * time.Second

// DefaultConfig ...
func DefaultConfig() *Config {
	return &Config{
		SlowRequest: time.Second,
		Debug:       false,
		Projects:    map[string]*ProjectConfig{},
		_logger:     onion_log.DefaultLogger(),
	}
}

// UseConfig 获取gin的配置参数
func UseConfig(name string) *Config {
	return RawConfig("yc." + name)
}

// RawConfig 查询配置是否有配置
func RawConfig(key string) *Config {
	var config = DefaultConfig()
	if conf.Detail().Get(key) == nil {
		config._logger.Warnf("config key: %s is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, &config); err != nil {
		config._logger.Panicf("request parse conf panic key:%s err:%s", key, err)
	}

	urlKey := fmt.Sprintf("%sUrls", key)
	urls := map[string]string{}
	if conf.Detail().Get(urlKey) == nil {
		config._logger.Warnf("config key: %s is not exists", urlKey)
	}
	if err := conf.Detail().UnmarshalKey(urlKey, &urls); err != nil {
		config._logger.Panicf("request parse conf panic key:%s err:%s", urlKey, err)
	}

	for k, c := range config.Projects {
		if url, exists := urls[k]; exists && url != "" {
			c.URL = url
		}
	}

	return config
}

// Build 返回单例
func (c *Config) Build() *Requests {
	r := &Requests{
		Config: c,
		clients: sync.Pool{
			New: func() interface{} {
				return resty.New()
			},
		},
		_logger: onion_log.DefaultLogger(),
	}
	return r
}
