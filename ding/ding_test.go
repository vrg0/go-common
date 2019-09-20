package ding

import "testing"

var DingUrl = "https://oapi.dingtalk.com/robot/send?access_token=XXXXXXXXXX"

func TestInit(t *testing.T) {
	Init(true)
}

func TestSendLink(t *testing.T) {
	TestInit(t)
	if e := SendLink("测试", "吧啦吧啦吧", "https://www.baidu.com", DingUrl); e != nil {
		t.Error(e)
	}
}

func TestSendText(t *testing.T) {
	TestInit(t)
	if e := SendText("测试吧啦吧啦吧", DingUrl); e != nil {
		t.Error(e)
	}
}