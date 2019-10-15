package notify

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vrg0/go-common/util"
	"net/http"
	"regexp"
	"strings"
)

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

type Notify struct {
	dstList []string
	ignore  []string
	re      []string
}

func (n *Notify) SetIgnore(ignore []string) {
	n.ignore = ignore
}

func (n *Notify) SetIgnoreRegexp(re []string) {
	n.re = re
}

func New(dstList []string) *Notify {
	return &Notify{
		dstList: dstList,
		ignore:  make([]string, 0),
	}
}

func (n *Notify) isIgnore(body string) bool {
	for _, sub := range n.ignore {
		if strings.Contains(body, sub) {
			return true
		}
	}
	for _, r := range n.re {
		if ok, _ := regexp.MatchString(r, body); ok {
			return true
		}
	}
	return false
}

func (n *Notify) SendText(body string) {
	if n.isIgnore(body) {
		return
	}
	msg := &Message{
		MsgType: "text",
		Text:    Text{Content: body},
	}

	for _, dst := range n.dstList {
		_ = n.sendMsg(msg, dst)
	}
}

func (n *Notify) SendLink(title string, body string, url string) {
	if n.isIgnore(body) {
		return
	}
	msg := &Message{
		MsgType: "link",
		Link: Link{
			Title:      title,
			Text:       body,
			MessageUrl: url,
		},
	}

	for _, dst := range n.dstList {
		_ = n.sendMsg(msg, dst)
	}
}

func (n *Notify) sendMsg(message *Message, dst string) error {
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
