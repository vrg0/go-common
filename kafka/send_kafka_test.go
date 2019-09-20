package kafka

import "testing"

var testSender *Sender = nil

func TestInitSender(t *testing.T) {
	if e := InitSender([]string{"127.0.0.1:9092"}); e != nil {
		t.Error(e)
	}
}

func TestSendMsg(t *testing.T) {
	TestInitSender(t)
	if e := SendMsg("test", "123"); e != nil {
		t.Error(e)
	}
}

func TestNewSender(t *testing.T) {
	if sender, e := NewSender([]string{"127.0.0.1:9092"}); e != nil {
		t.Error(e)
	} else {
		testSender = sender
	}
}

func TestSender_SendMsg(t *testing.T) {
	TestNewSender(t)
	if e := testSender.SendMsg("test2", "234"); e != nil {
		t.Error(e)
	}
}
