package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Recver struct {
	cluster []string
	config  *sarama.Config
	logger  *log.Logger
}

func (recver *Recver) Print(args ...interface{}) {
	if recver.logger != nil {
		recver.logger.Print(args...)
	}
}

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

type ConsumerCallback func(msg *sarama.ConsumerMessage)

type consumer struct {
	callback ConsumerCallback
}

//新建一个接收器，logger可为nil
func newConsumer(callback ConsumerCallback) *consumer {
	return &consumer{
		callback: callback,
	}
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		c.callback(message)
		session.MarkMessage(message, "")
	}
	return nil
}

func (recver *Recver) ListenAndRecvMsg(groupId string, topics []string, callback ConsumerCallback) error {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(recver.cluster, groupId, recver.config)
	if err != nil {
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				recver.Print(err)
			}
		}()

		consumerObj := newConsumer(callback)
		for {
			if err := client.Consume(ctx, topics, consumerObj); err != nil {
				recver.Print(err)
				continue
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		cancel()
		_ = client.Close()
	case <-sigterm:
		cancel()
		_ = client.Close()
	}

	return nil
}
