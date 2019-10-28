package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

type Sender struct {
	producer sarama.SyncProducer
}

var (
	defaultSender *Sender = nil
)

func NewSender(cluster []string, id string, kafkaVersion sarama.KafkaVersion) (*Sender, error) {
	if id == "" {
		id = "kafka_sender"
	}

	config := sarama.NewConfig()
	config.ClientID = id
	config.Version = kafkaVersion
	config.Producer.RequiredAcks = sarama.WaitForAll          //等待服务器所有副本都保存成功后的响应
	config.Producer.Partitioner = sarama.NewRandomPartitioner //随机的分区类型
	config.Producer.Return.Successes = true                   //是否等待成功和失败后的响应
	config.Producer.Timeout = 3 * time.Second                 //3秒超时

	//使用给定代理地址和配置创建一个同步生产者
	producer, e := sarama.NewSyncProducer(cluster, config)
	if e != nil {
		return nil, errors.New("new sync producer err")
	}

	rtn := &Sender{
		producer: producer,
	}

	return rtn, nil
}

func (s *Sender) SendMsg(topic string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	if partition, offset, e := s.producer.SendMessage(msg); e != nil {
		errStr := fmt.Sprintf("SendMsg err %s : partition:%d, offset:%d, %s:%s",
			e.Error(), partition, offset, topic, value)
		return errors.New(errStr)
	}

	return nil
}

func InitSender(cluster []string, id string, kafkaVersion sarama.KafkaVersion) error {
	if defaultSender == nil {
		if sender, e := NewSender(cluster, id, kafkaVersion); e != nil {
			return e
		} else {
			defaultSender = sender
		}
	}
	return nil
}

func SendMsg(topic string, value string) error {
	return defaultSender.SendMsg(topic, value)
}
