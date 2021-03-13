package okafka

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/kms9/publicyc/pkg/worker"
)

var _ worker.Worker = &Consumer{}

// Consumer 消费者
type Consumer struct {
	config              *Config
	Client              *cluster.Consumer
	IsConnected         bool
	sig                 chan int64
	msgProcess          func(msg *sarama.ConsumerMessage, more bool) error
	errProcess          func(err error, more bool)
	notificationProcess func(notification *cluster.Notification, more bool)
}

// InitConsumerWorker 加载消费worker
func (c *Consumer) InitConsumerWorker(mqConfig *cluster.Config) *Consumer {
	client, err := cluster.NewConsumer(c.config.Address, c.config.Consumer.Group, []string{c.config.Topic}, mqConfig)
	if err != nil {
		c.config._logger.Panicf("Create kafka consumer error: %v. address:%s topic:%s group:%s", err, c.config.Address, c.config.Topic, c.config.Consumer.Group)
	}
	c.Client = client
	c.errProcess = func(err error, more bool) {
		if more {
			c.config._logger.Warnf("Kafka consumer error: %v", err.Error())
		}
	}
	c.notificationProcess = func(notification *cluster.Notification, more bool) {
		if more {
			c.config._logger.Infof("Kafka consumer rebalance: %v", notification)
		}
	}
	c.sig = make(chan int64, 1)
	c.config._logger.Info("kafka consumer success!!!")
	c.IsConnected = true
	return c
}

// LoadMsgProcess 加载处理消息函数
func (c *Consumer) LoadMsgProcess(process func(msg *sarama.ConsumerMessage, more bool) error) *Consumer {
	c.msgProcess = process
	return c
}

// LoadErrProcess 加载错误处理函数
func (c *Consumer) LoadErrProcess(process func(err error, more bool)) *Consumer {
	c.errProcess = process
	return c
}

// LoadNotificationProcess 加载通知处理函数
func (c *Consumer) LoadNotificationProcess(process func(notification *cluster.Notification, more bool)) *Consumer {
	c.notificationProcess = process
	return c
}

// Name 消费者名称
func (c *Consumer) Name() string {
	return fmt.Sprintf("kafka: [topic %s group %s]", c.config.Topic, c.config.Consumer.Group)
}

// Start 开启消费
func (c *Consumer) Start() error {
	if c.msgProcess == nil {
		return errors.New("Kafka consumer msgProcess func unloaded")
	}
	go c.consume()
	return nil
}

// IndependentStart 不使用框开启消费
func (c *Consumer) IndependentStart() {
	if c.msgProcess == nil {
		c.config._logger.Panic("Kafka consumer msgProcess func unloaded")
		return
	}
	go c.consume()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt, syscall.SIGTERM)
	if _, ok := <-ch; ok {
		_ = c.Stop()
	}
}

// Stop 停止消费
func (c *Consumer) Stop() error {
	c.config._logger.Error("Stop kafka consumer server...")
	c.sig <- 1
	return c.Client.Close()
}

// consume 消费消息
func (c *Consumer) consume() {
	for {
		select {
		case msg, more := <-c.Client.Messages():
			err := c.msgProcess(msg, more)
			if err != nil {
				c.config._logger.Errorf("kafka process msg error: %v", err)
			}
			if err == nil && more {
				c.Client.MarkOffset(msg, "")
			}
		case err, more := <-c.Client.Errors():
			c.errProcess(err, more)
		case ntf, more := <-c.Client.Notifications():
			c.notificationProcess(ntf, more)
		case <-c.sig:
			return
		}
	}
}
