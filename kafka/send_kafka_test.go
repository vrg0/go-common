package kafka

import (
	"github.com/Shopify/sarama"
	"testing"
)

func TestInitSender(t *testing.T) {
	if e := InitSender([]string{"192.168.29.186:9092"}, "test222", sarama.V1_0_0_0); e != nil {
		t.Error(e)
	}
}

func TestSendMsg(t *testing.T) {
	TestInitSender(t)
	if e := SendMsg("test", "123"); e != nil {
		t.Error(e)
	}
}
