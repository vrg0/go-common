package kafka

import (
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"log"
)

type Recver struct {
	cluster []string
	config  *sarama.Config
	logger  *log.Logger
}

type ConsumerCallback func(msg *sarama.ConsumerMessage)

func NewRecver(version sarama.KafkaVersion, cluster []string, logger *log.Logger) *Recver {
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                       //设置offset
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin //设置消费方式

	return &Recver{
		cluster: cluster,
		config:  config,
		logger:  logger,
	}
}

func (recver *Recver) print(args ...interface{}) {
	if recver.logger != nil {
		recver.logger.Print(args...)
	}
}

type coreConsumer struct {
	callback ConsumerCallback
}

func (c *coreConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *coreConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *coreConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.callback(message)
		session.MarkMessage(message, "")
	}
	return nil
}

type Consumer struct {
	cancel context.CancelFunc
	client sarama.ConsumerGroup
}

func (consumer *Consumer) Close() {
	consumer.cancel()
	_ = consumer.client.Close()
}

func (recver *Recver) NewConsumer(groupId string, topics []string, callback ConsumerCallback) (*Consumer, error) {
	if callback == nil {
		return nil, errors.New("callback can not be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(recver.cluster, groupId, recver.config)
	if err != nil {
		return nil, err
	}

	rtn := &Consumer{
		cancel: cancel,
		client: client,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				recver.print(err)
			}
		}()

		for {
			if err := client.Consume(ctx, topics, &coreConsumer{callback: callback}); err != nil {
				recver.print(err)
				continue
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return rtn, nil
}
