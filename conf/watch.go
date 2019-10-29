package conf

import (
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

func mapInterfaceToString(interfaceMap map[string]interface{}) map[string]string {
	stringMap := make(map[string]string)
	for k, v := range interfaceMap {
		if vStr, ok := v.(string); ok {
			stringMap[k] = vStr
		}
	}
	return stringMap
}

func (c *Conf) startWatch() {
	watchChan := c.ago.Watch()
	go func() {
		defer func() {
			if e := recover(); e != nil {
				time.Sleep(time.Second * 1)
				go c.startWatch()
			}
		}()

		for {
			select {
			case w := <-watchChan:
				oldValue := mapInterfaceToString(w.OldValue)
				newValue := mapInterfaceToString(w.NewValue)

				for _, v := range c.namespaceHandler {
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

				for _, v := range c.keyHandler {
					if v.Namespace == w.Namespace && oldValue[v.Key] != newValue[v.Key] {
						v.Handler(oldValue[v.Key], newValue[v.Key])
					}
				}
			}
		}
	}()
}

func (c *Conf) WatchNamespace(namespace string, handler func(oldCfgs map[string]string, newCfgs map[string]string)) {
	//添加处理函数
	newNamespaceHandler := make([]*watchNamespaceHandler, 0)
	for _, watchHandler := range c.namespaceHandler {
		newNamespaceHandler = append(newNamespaceHandler, watchHandler)
	}
	newNamespaceHandler = append(newNamespaceHandler, &watchNamespaceHandler{
		Namespace: namespace,
		Handler:   handler,
	})
	c.namespaceHandler = newNamespaceHandler

	//首次加载数据
	handler(make(map[string]string), c.GetNamespace(namespace))
}

func (c *Conf) Watch(namespace string, key string, handler func(oldCfg string, newCfg string)) {
	//加载处理函数
	newKeyHandler := make([]*watchKeyHandler, 0)
	for _, watchHandler := range c.keyHandler {
		newKeyHandler = append(newKeyHandler, watchHandler)
	}
	newKeyHandler = append(newKeyHandler, &watchKeyHandler{
		Namespace: namespace,
		Key:       key,
		Handler:   handler,
	})
	c.keyHandler = newKeyHandler

	//首次加载数据
	kv := c.GetNamespace(namespace)
	if v, ok := kv[key]; ok {
		handler("", v)
	}
}
