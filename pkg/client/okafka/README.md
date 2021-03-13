## 使用方式
### 配置
```yaml
yc:
    kafka:
        pk:
          address:
            - 172.0.0.1:9093
            - 172.0.0.2:9093
            - 172.0.0.3:9093
          tls: true
          certFile: ca-cert # ca证书文件
          topic: ladder_pk_record_test
          sasl:
            enable: true
            user: username
            password: password
            handshake: true
          producer:
            returnSuccess: true # 是否接收成功消息
            retryBackoff: 2000 # 毫秒
          consumer:
            group: CID_alikafka_ladder_consumer_test
            returnErrors: true
            returnNotifications: true
```

### 生产者
#### 简单方式（使用已有配置）
```go
package setup

import (
	"sync"

	"github.com/kms9/publicyc/pkg/client/okafka"
)

var kafkaOnce sync.Once
var Kafka *okafka.Producer

// StartKafka 启动并初始化
func StartKafka() {
	kafkaOnce.Do(func() {
		Kafka = okafka.UseConfig("pk").BuildProducer()
	})
}


```

#### 指定配置
```go
package setup

import (
	"sync"

	"github.com/kms9/publicyc/pkg/client/okafka"
)

var kafkaOnce sync.Once
var Kafka *okafka.Producer

// StartKafka 启动并初始化
func StartKafka() {
	kafkaOnce.Do(func() {
		config := okafka.UseConfig("pk")
		mqConfig := config.NewProducerConfig()
		// TODO: 设置自定义配置
		Kafka = config.BuildNewProducer(mqConfig)
	})
}

```

#### 发送消息
```go
// 单条消息
Kafka.sendMessage("test", "Hello World")

// 多条消息
msgs := make([]*okafka.Message, 10)
for i := 0; i < 10; i++ {
    msgs[i] = &okafka.Message{
        Key:     "test",
        Content: "hello world" + strconv.Itoa(i),
    }
}
err := Kafka.SendMessages(msgs)
```

### 消费者
#### 简单方式（使用已有配置）
```go
func msgProcess(msg *sarama.ConsumerMessage, more bool) error {
	if more {
		onion_log.Infof("studyDataCenter Partition:%d, Offset:%d, Key:%s, Value:%s \n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
	return nil
}
func (eng *Engine) serveConsumer() error {
	consumer := okafka.UseConfig("pk").BuildConsumer().LoadMsgProcess(msgProcess)
	return eng.Job(consumer)
}

```
#### 自定义配置
```go
func (eng *Engine) serveConsumer() error {
    config := okafka.UseConfig("pk")
    mqConfig := config.NewConsumerConfig()
    // TODO: 自定义配置
    consumer := config.BuildNewConsumer(mqConfig)
    consumer = consumer.LoadMsgProcess(msgProcess)
	return eng.Job(consumer)
}

```
