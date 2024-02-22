package mq

import (
	"context"
	"engine-service/config"
	"fmt"
	"github.com/Shopify/sarama"
	"go-micro.dev/v4/logger"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	group    sarama.ConsumerGroup
	topics   []string
	handler  MsgHandler
	close    chan struct{}
}

type MsgHandler func(msg *sarama.ConsumerMessage) error

func NewKafkaClient() (*KafkaClient, error) {
	conf := sarama.NewConfig()

	conf.Consumer.Return.Errors = true
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Consumer.Offsets.AutoCommit.Enable = true

	producer, err := sarama.NewSyncProducer(config.Conf.Kafka.Brokers, conf)
	if err != nil {
		return nil, err
	}

	group, err := sarama.NewConsumerGroup(config.Conf.Kafka.Brokers, "engine", conf)
	if err != nil {
		return nil, err
	}

	consumer, err := sarama.NewConsumer(config.Conf.Kafka.Brokers, conf)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
		group:    group,
		close:    make(chan struct{}, 1),
	}, nil
}

func (kc *KafkaClient) Produce(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	_, _, err := kc.producer.SendMessage(msg)
	return err
}

func (kc *KafkaClient) Consume(topic string, handler MsgHandler) {
	partitionConsumer, err := kc.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		// 初次启动时
		logger.Errorf("Error occurred while consuming message,again in 5 seconds: %+v", err)
		panic(fmt.Sprintf("Error occurred while consuming message: %+v", err))
	}

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			handler(msg)
			logger.Infof("Consumed message offset %d\n", msg.Offset)
		}
	}
}

func (kc *KafkaClient) Group(topic string) {
	c := &ConsumerGroupHandler{
		kc: kc,
	}
	for {
		err := kc.group.Consume(context.Background(), []string{topic}, c)
		if err != nil {
			logger.Fatalf("Error occurred while consuming message: %+v", err)
		}
	}
}

func (kc *KafkaClient) Close() {
	kc.close <- struct{}{}
	if err := kc.producer.Close(); err != nil {
		logger.Errorf("close kafka producer error: %v", err)
	}

	if err := kc.consumer.Close(); err != nil {
		logger.Errorf("close kafka consumer error: %v", err)
	}
}

type ConsumerGroupHandler struct {
	kc *KafkaClient
}

func (c *ConsumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)
		err := c.kc.handler(msg)
		if err != nil {
			logger.Errorf("Error occurred while handling message: %+v", err)
			return err
		}
		logger.Infof("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)

		// 确认消息已被处理，提交偏移量
		session.MarkMessage(msg, "")
	}

	logger.Infof("consume message success")
	return nil
}

func (kc *KafkaClient) GroupMsg(topic string, handler MsgHandler) {
	kc.handler = handler
	kc.Group(topic)
}
