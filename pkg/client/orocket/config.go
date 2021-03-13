package orocket

import (
	"github.com/kms9/publicyc/pkg/conf"
	onionLog "github.com/kms9/publicyc/pkg/onion-log"
)

// Config 配置数据
type Config struct {
	Name           string   `json:"name"`
	Endpoint       []string `json:"endpoint"`       // 阿里云的TCP接入地址
	HTTPEndpoint   string   `json:"httpEndpoint"`   // 阿里云的HTTP接入地址
	HTTPStartCheck bool     `json:"httpStartCheck"` // 阿里云的HTTP启动check
	Access         string   `json:"access"`         // 阿里云账号的AccessKeyId
	Secret         string   `json:"secret"`         // 阿里云账号的AccessKeySecret
	Group          string   `json:"group"`          // 阿里云RocketMQ的控制台上获取的GID
	InstanceID     string   `json:"instanceId"`     // Topic所属实例ID，默认实例为空
	Topic          string   `json:"topic"`          // 阿里云创建的topic
	Retries        int      `json:"retries"`        // 重试次数
	LogLevel       string   `json:"logLevel"`       // 设置内置日志打印级别 enum: debug,info,warn,error
	Tag            string   `json:"tag"`            // 消费tag,默认为 "" 有值为: "TagA || TagC"
	_logger        *onionLog.Log
}

// DefaultRocketMQConfig 默认配置
func DefaultRocketMQConfig() *Config {
	return &Config{
		Retries: 2,
		_logger: onionLog.DefaultLogger(),
	}
}

// UseConfig 使用配置
func UseConfig(name string) *Config {
	return RawRocketConfig("yc.rocket." + name)
}

// RawRocketConfig ...
func RawRocketConfig(key string) *Config {
	var config = DefaultRocketMQConfig()
	if conf.Detail().Get(key) == nil {
		onionLog.Panicf("key:%s rocket Config is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, &config); err != nil {
		onionLog.Panicf("unmarshal rocket Config key:%s err:%s", key, err)
	}
	config.Name = key
	return config
}

// BuildPushConsumer 创建consumer
func (c *Config) BuildPushConsumer() (*PushConsumer, error) {
	consumer := &PushConsumer{config: c}
	return consumer.Init()
}

// BuildProducer 创建producer
func (c *Config) BuildProducer() (*Producer, error) {
	producer := &Producer{config: c}
	return producer.Init()
}

//// BuildAliConsumer 基于阿里云http-sdk 消费者
//func (c *Config) BuildAliConsumer() *AliConsumer {
//	producer := &AliConsumer{config: c}
//	return producer.Init()
//}

// BuildAliProducer 基于阿里云http-sdk 生产者
func (c *Config) BuildAliProducer() *AliProducer {
	producer := &AliProducer{config: c}
	return producer.Init()
}
