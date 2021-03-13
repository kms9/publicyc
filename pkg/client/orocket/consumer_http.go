package orocket

import (
	"fmt"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
)

// AliConsumer 基于阿里云http-sdk的消费者
type AliConsumer struct {
	config *Config

	client   mq_http_sdk.MQClient
	consumer mq_http_sdk.MQConsumer

	done       chan bool
	msgProcess func() error
}

// Init 初始化
func (a *AliConsumer) Init() *AliConsumer {
	a.client = mq_http_sdk.NewAliyunMQClient(
		a.config.Endpoint[0],
		a.config.Access, a.config.Secret,
		"",
	)
	a.consumer = a.client.GetConsumer(
		a.config.InstanceID,
		a.config.Topic,
		a.config.Group,
		a.config.Tag,
	)
	return a
}

// LoadMsgProcess TODO: 加载处理函数
func (a *AliConsumer) LoadMsgProcess() *AliConsumer {
	return a
}

// Start TODO: 开始消费
func (a *AliConsumer) Start() error {
	return nil
}

// Stop 停止消费
func (a *AliConsumer) Stop() error {
	a.done <- true
	return nil
}

// Name ...
func (a *AliConsumer) Name() string {
	return fmt.Sprintf("rocket-http-consumer: [topic %s group %s]", a.config.Topic, a.config.Group)
}
