package okafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/kms9/publicyc/pkg/conf"
	onionLog "github.com/kms9/publicyc/pkg/onion-log"
)

var (
	accessFromUser = 0
	colon          = ":"
)

// Config ..
type Config struct {
	Address  []string `json:"address"`
	TLS      bool     `json:"tls"`
	CertFile string   `json:"certFile"`
	Topic    string   `json:"topic"`
	SASL     struct {
		Enable    bool   `json:"enable"`
		User      string `json:"user"`
		Password  string `json:"password"`
		Handshake bool   `json:"handshake"`
	} `json:"sasl"`
	Cert     *x509.CertPool `json:"certPool,omitempty"`
	Producer struct {
		ReturnSuccess bool  `json:"returnSuccess"` // 是否接收成功消息
		RetryBackoff  int64 `json:"retryBackoff"`
	} `json:"producer"`
	Consumer struct {
		Group               string `json:"group"`
		ReturnErrors        bool   `json:"returnErrors"`
		ReturnNotifications bool   `json:"returnNotifications"`
	} `json:"consumer"`
	_logger *onionLog.Log
}

// loadCert 加载证书
func loadCert(certFileName string) *x509.CertPool {
	workPath, _ := os.Getwd()
	fullPath := filepath.Join(filepath.Join(workPath, "config"), certFileName)
	certBytes, err := ioutil.ReadFile(fullPath)
	if err != nil {
		onionLog.Panicf("kafka cert file read err err:%s", err)
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		onionLog.Panicf("kafka failed to parse root certificate")
	}
	return clientCertPool
}

// DefaultKafkaConfig 默认配置
func DefaultKafkaConfig() *Config {
	c := &Config{}
	c.Producer.ReturnSuccess = true
	c.SASL.Enable = true
	c.SASL.Handshake = true
	c.TLS = true
	c._logger = onionLog.DefaultLogger()
	return c
}

// UseConfig 使用配置
func UseConfig(name string) *Config {
	return RawKafkaConfig("yc.kafka." + name)
}

// RawKafkaConfig ...
func RawKafkaConfig(key string) *Config {
	var config = DefaultKafkaConfig()
	if conf.Detail().Get(key) == nil {
		onionLog.Panicf("key:%s kafaka config is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, config); err != nil {
		onionLog.Panicf("unmarshal kafaka config key:%s err:%s", key, err)
	}
	if config.CertFile != "" {
		config.Cert = loadCert(config.CertFile)
	}
	return config
}

// NewProducerConfig producer config
func (c *Config) NewProducerConfig() *sarama.Config {
	mqConfig := sarama.NewConfig()
	mqConfig.Net.SASL.Enable = c.SASL.Enable
	mqConfig.Net.SASL.User = c.SASL.User
	mqConfig.Net.SASL.Password = c.SASL.Password
	mqConfig.Net.SASL.Handshake = c.SASL.Handshake
	mqConfig.Net.TLS.Enable = c.TLS
	mqConfig.Producer.Return.Successes = c.Producer.ReturnSuccess
	mqConfig.Metadata.Retry.Backoff = time.Duration(c.Producer.RetryBackoff) * time.Millisecond
	if c.Cert != nil {
		mqConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            c.Cert,
			InsecureSkipVerify: true,
		}
	}
	if err := mqConfig.Validate(); err != nil {
		c._logger.Panicf("Kafka producer config invalidate.err: %v", err)
	}
	return mqConfig
}

// BuildNewProducer 创建producer
func (c *Config) BuildNewProducer(mqConfig *sarama.Config) *Producer {
	p := &Producer{
		config: c,
	}
	p.NewProducer(mqConfig)
	return p
}

// BuildProducer 默认配置创建producer
func (c *Config) BuildProducer() *Producer {
	mqConfig := c.NewProducerConfig()
	return c.BuildNewProducer(mqConfig)
}

// NewConsumerConfig consumer config
func (c *Config) NewConsumerConfig() *cluster.Config {
	mqConfig := cluster.NewConfig()
	mqConfig.Net.SASL.Enable = c.SASL.Enable
	mqConfig.Net.SASL.User = c.SASL.User
	mqConfig.Net.SASL.Password = c.SASL.Password
	mqConfig.Net.SASL.Handshake = c.SASL.Handshake
	if c.Cert != nil {
		mqConfig.Net.TLS.Config = &tls.Config{
			RootCAs:            c.Cert,
			InsecureSkipVerify: true,
		}
	}
	mqConfig.Net.TLS.Enable = c.TLS
	mqConfig.Consumer.Return.Errors = c.Consumer.ReturnErrors
	mqConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	mqConfig.Group.Return.Notifications = c.Consumer.ReturnNotifications

	mqConfig.Version = sarama.V0_10_0_0
	if err := mqConfig.Validate(); err != nil {
		c._logger.Panicf("Kafka consumer config invalidate. config: %v. err: %v", *mqConfig, err)
	}
	return mqConfig
}

// BuildNewConsumer 创建consumer
func (c *Config) BuildNewConsumer(mq *cluster.Config) *Consumer {
	consumer := &Consumer{
		config: c,
	}
	consumer.InitConsumerWorker(mq)
	return consumer
}

// BuildConsumer 默认配置创建consumer
func (c *Config) BuildConsumer() *Consumer {
	mqConfig := c.NewConsumerConfig()
	return c.BuildNewConsumer(mqConfig)
}
