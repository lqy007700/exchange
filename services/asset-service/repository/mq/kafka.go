package mq

import (
	"github.com/go-micro/plugins/v4/broker/kafka"
	"go-micro.dev/v4/broker"
)

// Service 封装了与 Kafka 交互的功能
type Service struct {
	producer broker.Broker
	consumer broker.Broker
}

// NewService 创建一个新的 Kafka 服务实例
func NewService() *Service {
	s := &Service{}
	err := s.Init()
	if err != nil {
		return nil
	}
	return s
}

// Init 初始化 Kafka 服务
func (s *Service) Init() error {

	b := kafka.NewBroker(broker.Addrs([]string{"127.0.0.1:9092"}...), broker.Addrs())

	err := b.Connect()
	if err != nil {
		return err
	}
	s.producer = b
	s.consumer = b

	return nil
}

// Close 关闭 Kafka 服务
func (s *Service) Close() {
	if s.producer != nil {
		s.producer.Disconnect()
	}
	if s.consumer != nil {
		s.consumer.Disconnect()
	}
}

// Producer 返回 Kafka 生产者
func (s *Service) Producer() broker.Broker {
	return s.producer
}

// Consumer 返回 Kafka 消费者
func (s *Service) Consumer() broker.Broker {
	return s.consumer
}
