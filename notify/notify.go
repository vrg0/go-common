package notify

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/vrg0/go-common/util"
	"golang.org/x/time/rate"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
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
	dstList      []string
	ignore       []string
	re           []string
	limitMap     map[string]*rate.Limiter
	limitMapLock *sync.RWMutex
	client       *http.Client
}

func (n *Notify) SetIgnoreWithLimiter(sub string, r rate.Limit, b int) {
	n.limitMapLock.Lock()
	n.limitMap[sub] = rate.NewLimiter(r, b)
	n.limitMapLock.Unlock()
}

func (n *Notify) SetIgnore(ignore []string) {
	n.ignore = ignore
}

func (n *Notify) SetIgnoreRegexp(re []string) {
	n.re = re
}

func New(dstList []string) *Notify {
	return &Notify{
		dstList:      dstList,
		ignore:       make([]string, 0),
		re:           make([]string, 0),
		limitMap:     make(map[string]*rate.Limiter),
		limitMapLock: new(sync.RWMutex),
		client: &http.Client{
			Timeout: time.Second * 10, //10秒超时
		},
	}
}

func (n *Notify) isIgnore(body string) bool {
	//字符串匹配
	for _, sub := range n.ignore {
		if strings.Contains(body, sub) {
			return true
		}
	}
	//正则匹配
	for _, r := range n.re {
		if ok, _ := regexp.MatchString(r, body); ok {
			return true
		}
	}
	//字符串匹配 + 最小值
	n.limitMapLock.RLock()
	for sub, l := range n.limitMap {
		if strings.Contains(body, sub) && !l.Allow() {
			return true
		}
	}
	n.limitMapLock.RUnlock()

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
		_, e := n.client.Post(dst, "application/json;charset=utf-8", data)
		if e == nil {
			break
		} else if i == 2 {
			return errors.Errorf("send msg error: %s %s %s", dst, msgStr, e.Error())
		}
	}

	return nil
}
