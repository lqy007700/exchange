package mq

import (
	"asset-service/config"
	"fmt"
	"github.com/Shopify/sarama"
	"go-micro.dev/v4/logger"
	"log"
	"time"
)

type KafkaClient struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	topics   []string
}

type MsgHandler func(msg *sarama.ConsumerMessage) error

func NewKafkaClient() (*KafkaClient, error) {
	conf := sarama.NewConfig()

	conf.Consumer.Return.Errors = true
	conf.Producer.Return.Errors = true
	conf.Producer.Return.Successes = true
	conf.Consumer.Offsets.AutoCommit.Enable = false

	logger.Infof("kafka brokers: %v", config.Conf.Kafka.Brokers)
	producer, err := sarama.NewSyncProducer(config.Conf.Kafka.Brokers, conf)
	if err != nil {
		logger.Fatalf("Error occurred while creating kafka producer: %+v", err)
		return nil, err
	}

	consumer, err := sarama.NewConsumer(config.Conf.Kafka.Brokers, conf)
	if err != nil {
		logger.Fatalf("Error occurred while creating kafka consumer: %+v", err)
		return nil, err
	}

	return &KafkaClient{
		producer: producer,
		consumer: consumer,
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
	isPanic := true
	for {
		partitionConsumer, err := kc.consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)

		if err != nil {
			// 初次启动时
			if isPanic {
				panic(fmt.Sprintf("Error occurred while consuming message: %+v", err))
			}

			logger.Errorf("Error occurred while consuming message,again in 5 seconds: %+v", err)
			// Sleep for a while before trying to reconnect
			time.Sleep(time.Second * 5)
			continue
		}
		isPanic = false

		for msg := range partitionConsumer.Messages() {
			err = handler(msg)
			if err != nil {
				logger.Errorf("Error occurred while handling message: %+v", err)
				continue
			}
			fmt.Printf("Consumed message offset %d\n", msg.Offset)
		}

		// If we reach here, it means the partitionConsumer has been closed and we need to reinitialize it.
		log.Println("Kafka connection closed, trying to reconnect...")
	}
}

func (kc *KafkaClient) Close() {
	if err := kc.producer.Close(); err != nil {
		logger.Errorf("close kafka producer error: %v", err)
	}

	if err := kc.consumer.Close(); err != nil {
		logger.Errorf("close kafka consumer error: %v", err)
	}
}
