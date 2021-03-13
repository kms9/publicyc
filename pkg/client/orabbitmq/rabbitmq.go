package orabbitmq

import (
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

var (
	errNotConnected  = errors.New("not connected to the producer")
	errAlreadyClosed = errors.New("already closed: not connected to the producer")
)

type RabbitMQ struct {
	name          string
	connection    *amqp.Connection
	channel       *amqp.Channel
	done          chan bool
	notifyClose   chan *amqp.Error
	notifyConfirm chan amqp.Confirmation
	IsConnected   bool
	config        *Config
	process       func(message []byte, d amqp.Delivery)
}

// Connection 连接rabbitMQ
func (rabbitMQ *RabbitMQ) Connection(name string) {
	rabbitMQ.name = name
	rabbitMQ.done = make(chan bool)
	rabbitMQ.connect()
	go rabbitMQ.handleReconnect()
	if !rabbitMQ.IsConnected {
		rabbitMQ.config._logger.Error("rabbitMQ connection failed")
	}
}

func (rabbitMQ *RabbitMQ) failOnError(err error, msg string) {
	if err != nil {
		rabbitMQ.config._logger.Errorf("%s: %s", msg, err)
	}
}

// 如果连接失败会不断重连
// 如果连接断开会重新连接
func (rabbitMQ *RabbitMQ) handleReconnect() {
	for {
		select {
		case <-rabbitMQ.done:
			rabbitMQ.config._logger.Warn("RabbitMQ connection is closed")
			return
		case <-rabbitMQ.notifyClose:
			rabbitMQ.config._logger.Error("RabbitMQ disconnection")
		}
		rabbitMQ.IsConnected = false
		rabbitMQ.config._logger.Info("Attempting to connect RabbitMQ ", rabbitMQ.name)
		for !rabbitMQ.connect() {
			rabbitMQ.config._logger.Error(rabbitMQ.name, " Failed to connect. Retrying...")
			time.Sleep(rabbitMQ.config.ReconnectDelay)
		}
	}
}

// 连接RabbitMQ
func (rabbitMQ *RabbitMQ) connect() bool {
	conn, err := amqp.Dial(rabbitMQ.config.getUrl())
	if err != nil {
		rabbitMQ.failOnError(err, "failed to connect to RabbitMQ")
		return false
	}
	ch, err := conn.Channel()
	if err != nil {
		rabbitMQ.failOnError(err, "failed to create channel")
		return false
	}
	ch.Confirm(false)

	rabbitMQ.changeConnection(conn, ch)
	rabbitMQ.IsConnected = true
	rabbitMQ.config._logger.Info(rabbitMQ.name, " RabbitMQ Connected!")
	return true
}

// 监听Rabbit channel的状态
func (rabbitMQ *RabbitMQ) changeConnection(connection *amqp.Connection, channel *amqp.Channel) {
	rabbitMQ.connection = connection
	rabbitMQ.channel = channel

	rabbitMQ.notifyClose = make(chan *amqp.Error)
	rabbitMQ.notifyConfirm = make(chan amqp.Confirmation)
	// channels没有必要主动关闭。如果没有协程使用它，它会被垃圾收集器收拾
	rabbitMQ.channel.NotifyClose(rabbitMQ.notifyClose)
	rabbitMQ.channel.NotifyPublish(rabbitMQ.notifyConfirm)
	// rabbitMQ.done <- true
	if rabbitMQ.process != nil {
		go rabbitMQ.consumer()
	}
}

// Close 关闭连接/信道
func (rabbitMQ *RabbitMQ) Close() error {
	if !rabbitMQ.IsConnected {
		return errors.New("not connected to the rabbitMQ")
	}
	rabbitMQ.done <- true
	err := rabbitMQ.channel.Close()
	if err != nil {
		return err
	}
	err = rabbitMQ.connection.Close()
	if err != nil {
		return err
	}
	close(rabbitMQ.done)
	rabbitMQ.IsConnected = false
	return nil
}

// LoadConsumer 加载消费函数
func (rabbitMQ *RabbitMQ) LoadConsumer(process func(message []byte, d amqp.Delivery)) *RabbitMQ {
	rabbitMQ.process = process
	return rabbitMQ
}

// 消费消息
func (rabbitMQ *RabbitMQ) consumer() {
	exchange := rabbitMQ.config.Exchange
	routerType := rabbitMQ.config.RouterType
	queueName := rabbitMQ.config.Queue
	ch := rabbitMQ.channel
	err := ch.ExchangeDeclare(exchange, routerType, true, false, false, false, nil)
	rabbitMQ.failOnError(err, "Failed to Declare a exchange")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	rabbitMQ.failOnError(err, "Failed to declare a queue"+q.Name)

	err = ch.QueueBind(queueName, queueName, exchange, false, nil)
	rabbitMQ.failOnError(err, "Failed to bind a queue")
	msgs, err := ch.Consume(
		queueName,                    // queue
		rabbitMQ.config.ConsumerName, // consumer
		false,                        // auto-ack
		false,                        // exclusive
		false,                        // no-local
		false,                        // no-wait
		nil,                          // args
	)
	rabbitMQ.failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		func() {
			defer func() {
				if err := recover(); err != nil {
					rabbitMQ.config._logger.Error("rabbitMQ handle message error: ", err)
				}
			}()
			rabbitMQ.process(d.Body, d)
		}()
	}
}

// Push 多次重传的发消息
func (rabbitMQ *RabbitMQ) Push(data []byte) error {
	if !rabbitMQ.IsConnected {
		return errNotConnected
	}
	var currentTime = 0
	for {
		err := rabbitMQ.UnsafePush(data)
		if err != nil {
			rabbitMQ.config._logger.Error("Push failed. Retrying...")
			currentTime += 1
			if int64(currentTime) < rabbitMQ.config.ResendTime {
				continue
			} else {
				return err
			}
		}
		ticker := time.NewTicker(rabbitMQ.config.ResendDelay)
		select {
		case confirm := <-rabbitMQ.notifyConfirm:
			if confirm.Ack {
				return nil
			}
		case <-ticker.C:
		}
		rabbitMQ.config._logger.Error("rabbitMQ message Push didn't confirm. Retrying...")
	}
}

// UnsafePush 发送出去，不管是否接受的到
func (rabbitMQ *RabbitMQ) UnsafePush(data []byte) error {
	if !rabbitMQ.IsConnected {
		return errors.New("")
	}
	return rabbitMQ.channel.Publish(
		rabbitMQ.config.Exchange,  // Exchange
		rabbitMQ.config.RouterKey, // Routing key
		false,                     // Mandatory
		false,                     // Immediate
		amqp.Publishing{
			DeliveryMode: 2,
			ContentType:  "text/plain",
			Body:         data,
			Timestamp:    time.Now(),
		},
	)
}

// Start ..
func (rabbitMQ *RabbitMQ) Start() error {
	if !rabbitMQ.IsConnected {
		return errors.New("not connected to the rabbitMQ")
	}
	if rabbitMQ.process == nil {
		return errors.New("process is not loaded")
	}
	go rabbitMQ.consumer()
	return nil
}

// Close ..
func (rabbitMQ *RabbitMQ) Stop() error {
	return rabbitMQ.Close()
}

// Name 消费者名称
func (rabbitMQ *RabbitMQ) Name() string {
	return fmt.Sprintf("rabbitMQ: [name %s]", rabbitMQ.name)
}
