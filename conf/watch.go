package conf

import (
	"github.com/shima-park/agollo"
	"sync"
	"time"
)

type watchNamespaceHandler struct {
	Namespace string
	Handler   func(map[string]string, map[string]string)
}

type watchKeyHandler struct {
	Namespace string
	Key       string
	Handler   func(string, string)
}

var (
	namespaceHandler     = make([]*watchNamespaceHandler, 0)
	namespaceHandlerLock = new(sync.RWMutex)
	keyHandler           = make([]*watchKeyHandler, 0)
	keyHandlerLock       = new(sync.RWMutex)
	startWatchOnce       = new(sync.Once)
)

func mapInterfaceToString(interfaceMap map[string]interface{}) map[string]string {
	stringMap := make(map[string]string)
	for k, v := range interfaceMap {
		if vStr, ok := v.(string); ok {
			stringMap[k] = vStr
		}
	}
	return stringMap
}

func startWatch() {
	watchChan := agollo.Watch()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				time.Sleep(time.Second * 1)
				go startWatch()
			}
		}()

		for {
			select {
			case w := <-watchChan:
				oldValue := mapInterfaceToString(w.OldValue)
				newValue := mapInterfaceToString(w.NewValue)

				//Watch Namespace
				namespaceHandlerLock.RLock()
				for _, v := range namespaceHandler {
					if v.Namespace == w.Namespace {
						oldV := make(map[string]string)
						for k, v := range oldValue {
							oldV[k] = v
						}
						newV := make(map[string]string)
						for k, v := range newValue {
							newV[k] = v
						}
						v.Handler(oldV, newV)
					}
				}
				namespaceHandlerLock.RUnlock()

				//Watch Key
				keyHandlerLock.RLock()
				for _, v := range keyHandler {
					if v.Namespace == w.Namespace {
						v.Handler(oldValue[v.Key], newValue[v.Key])
					}
				}
				keyHandlerLock.RUnlock()
			}
		}
	}()
}

func StartWatch() {
	startWatchOnce.Do(startWatch)
}

func WatchNamespace(namespace string, handler func(oldCfgs map[string]string, newCfgs map[string]string)) {
	//添加处理函数
	namespaceHandlerLock.Lock()
	defer namespaceHandlerLock.Unlock()
	namespaceHandler = append(namespaceHandler, &watchNamespaceHandler{
		Namespace: namespace,
		Handler:   handler,
	})

	//首次加载数据
	handler(make(map[string]string), GetNamespace(namespace))
}

func Watch(namespace string, key string, handler func(oldCfg string, newCfg string)) {
	//加载处理函数
	keyHandlerLock.Lock()
	defer keyHandlerLock.Unlock()
	keyHandler = append(keyHandler, &watchKeyHandler{
		Namespace: namespace,
		Key:       key,
		Handler:   handler,
	})

	//首次加载数据
	kv := GetNamespace(namespace)
	if v, ok := kv[key]; ok {
		handler("", v)
	}
}
