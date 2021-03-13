package orocket

import (
	"context"
	"errors"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

// PushConsumer 消费者
type PushConsumer struct {
	config     *Config
	client     rocketmq.PushConsumer
	msgProcess func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error)
}

// Init ..
func (c *PushConsumer) Init() (*PushConsumer, error) {
	pushConsumer, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer(c.config.Endpoint),
		consumer.WithGroupName(c.config.Group),
		consumer.WithNamespace(c.config.InstanceID),
		consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.config.Access,
			SecretKey: c.config.Secret,
		}),
		consumer.WithRetry(c.config.Retries),
	)
	c.client = pushConsumer
	rlog.SetLogLevel(c.config.LogLevel)
	return c, err
}

// LoadMsgProcess 加载处理消息函数
func (c *PushConsumer) LoadMsgProcess(process func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error)) *PushConsumer {
	c.msgProcess = process
	return c
}

// Start 暂时只支持一个topic消费
func (c *PushConsumer) Start() error {
	if c.msgProcess == nil {
		return errors.New("rocket mq msgProcess not load")
	}
	selector := consumer.MessageSelector{}
	if c.config.Tag != "" {
		selector.Type = consumer.TAG
		selector.Expression = c.config.Tag
	}
	if err := c.client.Subscribe(c.config.Topic, selector, c.msgProcess); err != nil {
		return err
	}
	return c.client.Start()
}

// Stop ..
func (c *PushConsumer) Stop() error {
	return c.client.Shutdown()
}

// Name ..
func (c *PushConsumer) Name() string {
	return fmt.Sprintf("rocket-tcp-consumer: [topic %s group %s]", c.config.Topic, c.config.Group)
}

// Client ..
func (c *PushConsumer) Client() rocketmq.PushConsumer {
	return c.client
}
