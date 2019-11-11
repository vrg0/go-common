package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"testing"
)

func TestNewRecver(t *testing.T) {
	x := NewRecver(sarama.V1_0_0_0, []string{"192.168.29.186:9092"}, log.New(os.Stdout, "", 0))
	if x == nil {
		t.Error("new recver")
	}
}

func TestRecver_ListenAndRecvMsg(t *testing.T) {
	x := NewRecver(sarama.V1_0_0_0, []string{"192.168.29.186:9092"}, log.New(os.Stdout, "", 0))
	_, err := x.NewConsumer("123", []string{"test"}, func(msg *sarama.ConsumerMessage) {
		fmt.Println(string(msg.Key), string(msg.Value))
	})

	if err != nil {
		t.Error(err)
	}

	select {}
}
