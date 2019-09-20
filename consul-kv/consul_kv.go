package consul_kv

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
	"time"
)

var defaultClient *ConsulKV = nil

func Init(dc string, cluster []string) error {
	if defaultClient == nil {
		defaultClient = New(dc, cluster)
		if defaultClient == nil {
			return errors.New("new consul_kv error")
		}
	}
	return nil
}

func GetValue(key string) (string, error) {
	return defaultClient.GetValue(key)
}

func SetValue(key string, value string) error {
	return defaultClient.SetValue(key, value)
}

func DelValue(key string) error {
	return defaultClient.DelValue(key)
}

func List(prefix string) ([]KVPair, error) {
	return defaultClient.List(prefix)
}

func WatchPrefix(prefix string, handler func(_ uint64,pairs api.KVPairs)) {
	defaultClient.WatchPrefix(prefix, handler)
}

type ConsulKV struct {
	kvClients    []*api.Client
	watchClients []*api.Client
}

//kv对
type KVPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//新建客户端
func newConsulClient(address string, dc string, waitTime time.Duration) *api.Client {
	config := api.DefaultConfig()
	config.Address = address
	config.Datacenter = dc
	config.WaitTime = waitTime
	client, _ := api.NewClient(config)
	return client
}

//新建
//注：集群模式未实现
func New(dc string, cluster []string) *ConsulKV {
	if len(cluster) == 0 {
		return nil
	}

	rtn := &ConsulKV{}
	for _, endpoint := range cluster {
		rtn.kvClients = append(rtn.kvClients, newConsulClient(endpoint, dc, time.Second*1))
		rtn.watchClients = append(rtn.watchClients, newConsulClient(endpoint, dc, time.Second*60*5))
	}
	return rtn
}

func (ckv *ConsulKV) GetValue(key string) (string, error) {
	rtn := ""
	for idx, kvClient := range ckv.kvClients {
		kv := kvClient.KV()
		kvPair, _, e := kv.Get(key, &api.QueryOptions{})
		if e != nil {
			if idx == len(ckv.kvClients)-1 {
				return "", e
			} else {
				continue
			}
		}
		if kvPair == nil {
			rtn = ""
			break
		} else {
			rtn = string(kvPair.Value)
			break
		}
	}

	return rtn, nil
}

func (ckv *ConsulKV) SetValue(key string, value string) error {
	for idx, kvClient := range ckv.kvClients {
		kv := kvClient.KV()
		_, e := kv.Put(&api.KVPair{Key: key, Value: []byte(value)}, &api.WriteOptions{})
		if e != nil {
			if idx == len(ckv.kvClients)-1 {
				return e
			}
		} else {
			break
		}
	}
	return nil
}

func (ckv *ConsulKV) DelValue(key string) error {
	for idx, kvClient := range ckv.kvClients {
		kv := kvClient.KV()
		_, e := kv.Delete(key, &api.WriteOptions{})
		if e != nil {
			if idx == len(ckv.kvClients)-1 {
				return e
			}
		} else {
			break
		}
	}
	return nil
}

func (ckv *ConsulKV) List(prefix string) ([]KVPair, error) {
	rtn := make([]KVPair, 0)

	for idx, kvClient := range ckv.kvClients {
		kv := kvClient.KV()
		kvPairs, _, e := kv.List(prefix, &api.QueryOptions{})
		if e != nil {
			if idx == len(ckv.kvClients)-1 {
				return rtn, e
			}
		} else {
			for i := len(kvPairs) - 1; i >= 0; i-- {
				rtn = append(rtn, KVPair{Key: kvPairs[i].Key, Value: string(kvPairs[i].Value)})
			}
			break
		}
	}

	return rtn, nil
}

func (ckv *ConsulKV) WatchPrefix(prefix string, handler func(uint64, api.KVPairs)) {
	parse, _ := watch.Parse(map[string]interface{}{"type": "keyprefix", "prefix": prefix})
	parse.Handler = func(u uint64, i interface{}) {
		handler(u, i.(api.KVPairs))
	}
	go ckv.watch(parse, 0)
}

func (ckv *ConsulKV) watch(parse *watch.Plan, sentry int) {
	defer func() {
		if e := recover(); e != nil {
			time.Sleep(time.Second * 1)
			if sentry < len(ckv.watchClients)-1 {
				sentry++
			} else {
				sentry = 0
			}
			ckv.watch(parse, sentry)
		}
	}()
	_ = parse.RunWithClientAndLogger(ckv.watchClients[sentry], nil)
}
