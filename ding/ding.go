package ding

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vrg0/go-common/util"
	"net/http"
	"strings"
)

/**
 * 使用前必须进行初始化
 */

type Message struct {
	MsgType string `json:"msgtype"`
	Text    Text   `json:"text"`
	Link    Link   `json:"link"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageUrl string `json:"messageUrl"`
}

type Text struct {
	Content string `json:"content"`
}

var (
	priEnvIsPro = false
)

func Init(envIsPro bool) {
	priEnvIsPro = envIsPro
}

func sendMsg(message *Message, dst string) error {
	msgBytes, _ := json.Marshal(message)
	msgStr := util.BytesString(msgBytes)

	for i := 0; i < 3; i++ {
		data := strings.NewReader(msgStr)
		_, e := http.Post(dst, "application/json;charset=utf-8", data)
		if e == nil {
			break
		} else if i == 2 {
			return errors.Errorf("send msg error: %s %s %s", dst, msgStr, e.Error())
		}
	}

	return nil
}

func SendText(body string, dsts ...string) error {
	//开发模式不发叮
	if !priEnvIsPro {
		return nil
	}

	if len(dsts) == 0 {
		return errors.New("dsts can not empty")
	}

	msg := &Message{
		MsgType: "text",
		Text:    Text{Content: body},
	}

	errStr := ""
	for _, dst := range dsts {
		if e := sendMsg(msg, dst); e != nil {
			errStr += e.Error() + "\n"
		}
	}

	if errStr != "" {
		return errors.New(errStr)
	} else {
		return nil
	}
}

func SendLink(title string, body string, url string, dsts ...string) error {
	//开发模式不发叮
	if !priEnvIsPro {
		return nil
	}
	if len(dsts) == 0 {
		return errors.New("dsts can not empty")
	}

	msg := &Message{
		MsgType: "link",
		Link: Link{
			Title:      title,
			Text:       body,
			MessageUrl: url,
		},
	}

	for _, dst := range dsts {
		_ = sendMsg(msg, dst)
	}

	return nil
}
