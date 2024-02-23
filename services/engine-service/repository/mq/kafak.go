package mq

import (
	"context"
	"engine-service/config"
	"fmt"
	"github.com/Shopify/sarama"
	"go-micro.dev/v4/logger"
	"time"
)

// Producer 包含Kafka生产者的配置和实例
type Producer struct {
	producer sarama.AsyncProducer
}

// NewProducer 创建一个新的Kafka生产者实例
func NewProducer(brokers []string) (*Producer, error) {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Compression = sarama.CompressionSnappy
	conf.Producer.Flush.Frequency = 500 * time.Millisecond
	conf.Producer.Flush.Messages = 100

	producer, err := sarama.NewAsyncProducer(config.Conf.Kafka.Brokers, conf)
	if err != nil {
		return nil, err
	}

	kp := &Producer{producer: producer}
	go func() {
		for {
			select {
			case success := <-kp.producer.Successes():
				fmt.Printf("Message sent successfully: %v\n", success)
			case err := <-kp.producer.Errors():
				fmt.Println("Failed to send message:", err)
			}
		}
	}()
	return kp, nil
}

// ProduceMessage 异步发送消息
func (p *Producer) ProduceMessage(topic string, message string) {
	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(message)}
	p.producer.Input() <- msg
}

func (p *Producer) Close() error {
	if err := p.producer.Close(); err != nil {
		logger.Errorf("Failed to shut down producer cleanly: %v", err)
		return err
	}
	logger.Infof("Producer shutdown cleanly")
	return nil
}

// Consumer 包含Kafka消费者组的配置和实例
type Consumer struct {
	group sarama.ConsumerGroup
}

// NewConsumer 创建一个新的Kafka消费者组实例
func NewConsumer(groupID string) (*Consumer, error) {
	conf := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(config.Conf.Kafka.Brokers, groupID, conf)
	if err != nil {
		return nil, err
	}

	kc := &Consumer{group: consumer}
	return kc, nil
}

func (c *Consumer) Consume(ctx context.Context, topics []string, handler *ConsumerHandler) error {
	go func() {
		<-ctx.Done()
		if err := c.group.Close(); err != nil {
			logger.Errorf("Failed to shut down consumer group cleanly: %v", err)
		} else {
			logger.Info("Consumer group shutdown cleanly")
		}
	}()

	for {
		if err := c.group.Consume(ctx, topics, handler); err != nil {
			logger.Infof("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type ConsumerHandler struct {
	HandleMessage func(message []byte)
	Ctx           context.Context
}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if ok {
				h.HandleMessage(msg.Value)
				sess.MarkMessage(msg, "")
			}
		case <-h.Ctx.Done():
			logger.Info("ConsumerHandler ConsumeClaim Done")
			return h.Ctx.Err()
		}
	}
}
