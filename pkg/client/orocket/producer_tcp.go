package orocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

// Producer 生产者
type Producer struct {
	config *Config
	client rocketmq.Producer
}

// Init ..
func (c *Producer) Init() (*Producer, error) {
	client, err := rocketmq.NewProducer(
		producer.WithNameServer(c.config.Endpoint),
		producer.WithCreateTopicKey(c.config.Topic),
		producer.WithGroupName(c.config.Group),
		producer.WithNamespace(c.config.InstanceID),
		producer.WithCredentials(primitive.Credentials{
			AccessKey: c.config.Access,
			SecretKey: c.config.Secret,
		}),
		producer.WithRetry(c.config.Retries),
	)
	c.client = client
	rlog.SetLogLevel(c.config.LogLevel)
	return c, err
}

// Start 暂时只支持一个topic消费
func (c *Producer) Start() error {
	return c.client.Start()
}

// Stop ..
func (c *Producer) Stop() error {
	return c.client.Shutdown()
}

// Name ..
func (c *Producer) Name() string {
	return fmt.Sprintf("rocket-tcp-producer: [topic %s group %s]", c.config.Topic, c.config.Group)
}

// Client ..
func (c *Producer) Client() rocketmq.Producer {
	return c.client
}

func (c *Producer) mergeMsg(property map[string]string, msg ...interface{}) ([]*primitive.Message, error) {
	notifies := make([]*primitive.Message, 0)
	for _, m := range msg {
		b, err := json.Marshal(m)
		if err != nil {
			return notifies, err
		}
		notify := primitive.NewMessage(c.config.Topic, b)
		if property != nil {
			notify.WithProperties(property)
		}
		notifies = append(notifies, notify)
	}
	return notifies, nil
}

// SendSync 批量同步发送
func (c *Producer) SendSync(ctx context.Context, property map[string]string, msg ...interface{}) error {
	notifies, err := c.mergeMsg(property, msg...)
	if err != nil {
		return err
	}
	_, err = c.Client().SendSync(ctx, notifies...)
	return err
}

// SendAsync 批量异步发送
func (c *Producer) SendAsync(ctx context.Context, mq func(ctx context.Context, result *primitive.SendResult, err error), property map[string]string, msg ...interface{}) error {
	notifies, err := c.mergeMsg(property, msg...)
	if err != nil {
		return err
	}
	err = c.Client().SendAsync(ctx, mq, notifies...)
	return err
}
