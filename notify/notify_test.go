package notify

import "testing"

var DingUrl = []string{"https://oapi.dingtalk.com/robot/send?access_token=XXXXXXXXXX"}

func TestNotify_SendLink(t *testing.T) {
	n := New(DingUrl)
	n.SendText("test")
	n.SendLink("test", "body", "www.baidu.com")
}
