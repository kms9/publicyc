package okafka

import (
	"encoding/json"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
)

// Producer 生产者
type Producer struct {
	config      *Config
	Client      sarama.SyncProducer
	IsConnected bool
}

// handleClose 监控关闭
func (p *Producer) handleClose() {
	if !p.IsConnected {
		return
	}
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	if _, ok := <-ch; ok {
		p.Client.Close()
	}
}

// NewProducer 创建producer
func (p *Producer) NewProducer(mqConfig *sarama.Config) {
	client, err := sarama.NewSyncProducer(p.config.Address, mqConfig)
	if err != nil {
		b, _ := json.Marshal(p.config)
		p.config._logger.Panicf("Kafka producer create fail. err: %v config: %s", err, string(b))
	}
	p.Client = client
	p.IsConnected = true
	p.config._logger.Info("Kafka producer connect SUCCESS")
	go p.handleClose()
}

// SendMessage 发送消息
func (p *Producer) SendMessage(key string, content string) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: p.config.Topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(content),
	}
	p.config._logger.Info(key, content)
	return p.Client.SendMessage(msg)
}

// Message 生产者消息
type Message struct {
	Key     string
	Content string
}

// SendMessages 批量发送消息
func (p *Producer) SendMessages(msgs []*Message) error {
	producerMsgs := make([]*sarama.ProducerMessage, len(msgs))
	for i, msg := range msgs {
		producerMsgs[i] = &sarama.ProducerMessage{
			Topic: p.config.Topic,
			Key:   sarama.StringEncoder(msg.Key),
			Value: sarama.StringEncoder(msg.Content),
		}
	}
	return p.Client.SendMessages(producerMsgs)
}
