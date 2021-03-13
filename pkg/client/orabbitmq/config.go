package orabbitmq

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kms9/publicyc/pkg/conf"
	onionLog "github.com/kms9/publicyc/pkg/onion-log"
)

var (
	accessFromUser = 0
	colon          = ":"
)

type Config struct {
	Name            string        `json:"name"`
	AccessKeyId     string        `json:"accessKeyId"`
	AccessKeySecret string        `json:"accessKeySecret"`
	ResourceOwnerId string        `json:"resourceOwnerId"`
	Endpoint        string        `json:"endpoint"`
	Vhost           string        `json:"vhost"`
	Exchange        string        `json:"exchange"`
	RouterKey       string        `json:"routerKey"`
	RouterType      string        `json:"type"`
	Queue           string        `json:"queue"`
	ConsumerName    string        `json:"consumerName"`
	ReconnectDelay  time.Duration `json:"reconnectDelay"` // 连接断开后多久重连
	ResendDelay     time.Duration `json:"resendDelay"`    // 消息发送失败后，多久重发
	ResendTime      int64         `json:"resendTime"`     // 消息重发次数
	_logger         *onionLog.Log
}

// DefaultRabbitMQConfig 默认配置
func DefaultRabbitMQConfig() Config {
	return Config{
		ReconnectDelay: 10 * time.Second,
		ResendDelay:    5 * time.Second,
		ResendTime:     3,
		ConsumerName:   "consumerOne",
		_logger:        onionLog.DefaultLogger(),
	}
}

// UseConfig 使用配置
func UseConfig(name string) Config {
	return RawRabbitConfig("yc.amqp." + name)
}

// RawRabbitConfig ...
func RawRabbitConfig(key string) Config {
	var config = DefaultRabbitMQConfig()
	if conf.Detail().Get(key) == nil {
		onionLog.Panicf("key:%s rabbitmqConfig is not exists", key)
	}
	if err := conf.Detail().UnmarshalKey(key, &config); err != nil {
		onionLog.Panicf("unmarshal rabbitmqConfig key:%s err:%s", key, err)
	}
	config.Name = key
	return config
}

func hmacSha1(keyStr string, message string) string {
	key := []byte(keyStr)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

func getUserName(ak string, instanceId string) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(accessFromUser))
	buffer.WriteString(colon)
	buffer.WriteString(instanceId)
	buffer.WriteString(colon)
	buffer.WriteString(ak)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func getPassword(sk string) string {
	now := time.Now()
	currentMillis := strconv.FormatInt(now.UnixNano()/1000000, 10)
	var buffer bytes.Buffer
	buffer.WriteString(strings.ToUpper(hmacSha1(currentMillis, sk)))
	buffer.WriteString(colon)
	buffer.WriteString(currentMillis)
	return base64.StdEncoding.EncodeToString(buffer.Bytes())
}

func (c Config) getUrl() string {
	var buf bytes.Buffer
	ak := c.AccessKeyId
	sk := c.AccessKeySecret
	instanceId := c.ResourceOwnerId // instanceId

	userName := getUserName(ak, instanceId)
	password := getPassword(sk)
	buf.WriteString("amqp://")
	buf.WriteString(userName)
	buf.WriteString(":")
	buf.WriteString(password)
	endpoint := c.Endpoint
	vhost := c.Vhost
	buf.WriteString(fmt.Sprintf("@%s/%s", endpoint, vhost))
	url := buf.String()
	return url
}

// Build 构建RabbitMQ客户端
func (c Config) Build() *RabbitMQ {
	rabbitMQ := &RabbitMQ{
		config: &c,
	}
	rabbitMQ.Connection(c.Name)
	return rabbitMQ
}
