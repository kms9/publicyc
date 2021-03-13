package orocket

import (
	"fmt"
	"time"

	mq_http_sdk "github.com/aliyunmq/mq-http-go-sdk"
)

// AliProducer 基于阿里云sdk的生产者
type AliProducer struct {
	config *Config

	client   mq_http_sdk.MQClient
	producer mq_http_sdk.MQProducer
}

// Init 初始化
func (a *AliProducer) Init() *AliProducer {
	a.client = mq_http_sdk.NewAliyunMQClient(
		a.config.HTTPEndpoint,
		a.config.Access, a.config.Secret,
		"",
	)
	return a
}

// Start ...
func (a *AliProducer) Start() error {
	a.producer = a.client.GetProducer(a.config.InstanceID, a.config.Topic)
	if a.config.HTTPStartCheck {
		// ping 测试,如果配置异常 会直接启动失败
		_, err := a.producer.PublishMessage(mq_http_sdk.PublishMessageRequest{
			MessageBody: `{"msg": "aliyun producer init"}`,
		})
		return err
	}
	return nil
}

// Stop ...
func (a *AliProducer) Stop() error {
	return nil
}

// Name ...
func (a *AliProducer) Name() string {
	return fmt.Sprintf("rocket-http-producer: [topic %s group %s instanceId %s]", a.config.Topic, a.config.Group, a.config.InstanceID)
}

// BuildMsg 构建消息
func (a *AliProducer) BuildMsg(body string, properties map[string]string) *AliPubMsg {
	if properties == nil {
		properties = map[string]string{}
	}
	return &AliPubMsg{
		producer: a.producer,
		msg: &mq_http_sdk.PublishMessageRequest{
			MessageBody: body,
			Properties:  properties,
		},
	}
}

// AliPubMsg 增加消息
type AliPubMsg struct {
	producer mq_http_sdk.MQProducer
	msg      *mq_http_sdk.PublishMessageRequest
}

// AddDelayEnd 延时到xxx时间开始发送消息
func (a *AliPubMsg) AddDelayEnd(t time.Time) *AliPubMsg {
	a.msg.StartDeliverTime = t.UnixNano() / 1e6
	return a
}

// AddTag 增加标签
func (a *AliPubMsg) AddTag(tag string) *AliPubMsg {
	a.msg.MessageTag = tag
	return a
}

// AddKey 增加messageKey 用的不多
func (a *AliPubMsg) AddKey(key string) *AliPubMsg {
	a.msg.MessageKey = key
	return a
}

// AddSharding 增加sharding 保证有序
func (a *AliPubMsg) AddSharding(shardingKey string) *AliPubMsg {
	a.msg.ShardingKey = shardingKey
	return a
}

// AliMsgRes 返回响应消息-作为扩展字段
type AliMsgRes struct {
	// 消息ID
	MessageId string `xml:"MessageId" json:"message_id"`
	// 消息体MD5
	MessageBodyMD5 string `xml:"MessageBodyMD5" json:"message_body_md5"`
}

// Send 发送消息
func (a *AliPubMsg) Send() (AliMsgRes, error) {
	ret, err := a.producer.PublishMessage(*a.msg)
	if err != nil {
		return AliMsgRes{}, err
	}
	return AliMsgRes{MessageBodyMD5: ret.MessageBodyMD5, MessageId: ret.MessageId}, nil
}
