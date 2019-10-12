package conf

/**
 * 从agollo配置中心中获取配置
 *
 * 支持kv映射功能
 * 把配置中的 "${key}" 映射成value
 * 支${key}的嵌套，如：${hello}值为world，${AAA}值为ll, 则${he${AAA}o}值为world
 * ${key}的嵌套最多支持8层
 */

import (
	"github.com/shima-park/agollo"
	"github.com/vrg0/go-common/args"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	kvMapReplaceDeep  = 8
	placeHolderPrefix = "${"
	placeHolderSuffix = "}"
)

var (
	defaultConf *Conf = nil
)

type Conf struct {
	kvMap                map[string]string
	ago                  agollo.Agollo
	logger               *log.Logger
	namespaceHandler     []*watchNamespaceHandler
	namespaceHandlerLock *sync.RWMutex
	keyHandler           []*watchKeyHandler
	keyHandlerLock       *sync.RWMutex
}

func New(configServer string, appId string, cluster string, cacheFilePath string, logger *log.Logger) (*Conf, error) {
	rtn := Conf{
		kvMap:                make(map[string]string),
		logger:               logger,
		namespaceHandler:     make([]*watchNamespaceHandler, 0),
		namespaceHandlerLock: new(sync.RWMutex),
		keyHandler:           make([]*watchKeyHandler, 0),
		keyHandlerLock:       new(sync.RWMutex),
	}

	if newAgo, err := agollo.New(
		configServer, appId,
		agollo.Cluster(cluster),                   //集群名称(idc)
		agollo.BackupFile(cacheFilePath),          //缓存文件的路径
		agollo.FailTolerantOnBackupExists(),       //从apollo service读取配置失败时，从缓存文件中读取配置
		agollo.AutoFetchOnCacheMiss(),             //当缓存中找不到namespace时，自动从apollo server拉取namespace
		agollo.LongPollerInterval(time.Second*10), //从apollo server更新数据的轮训时间
	); err != nil {
		return nil, err
	} else {
		rtn.ago = newAgo
	}

	//开启一个协程，从apollo server更新数据
	if logger != nil {
		errLogChan := rtn.ago.Start()
		go func() {
			defer func() {
				if e := recover(); e != nil {
					return
				}
			}()
			for {
				select {
				case e := <-errLogChan:
					rtn.logger.Print(e)
				}
			}
		}()
	} else {
		rtn.ago.Start()
	}

	//开启一个协程，启用watch机制
	rtn.startWatch()

	return &rtn, nil
}

//刷新kvMap
func (c *Conf) RefreshKvMap(kvMapper map[string]string) {
	newKvMap := make(map[string]string)
	for k, v := range kvMapper {
		newKvMap[k] = v
	}
	c.kvMap = newKvMap
}

func (c *Conf) internalKvMapReplace(value string, deep int) (string, bool) {
	if deep >= kvMapReplaceDeep {
		return "", false
	}

	startIndex := strings.Index(value, placeHolderPrefix)
	if startIndex == -1 {
		return value, true
	}

	endIndex := strings.Index(value, placeHolderSuffix)
	if endIndex == -1 {
		return value, true
	}

	innerStr := value[startIndex+len(placeHolderPrefix) : endIndex]
	if replaceValue, ok := c.kvMap[innerStr]; ok {
		value = strings.ReplaceAll(value, placeHolderPrefix+innerStr+placeHolderSuffix, replaceValue)
		return c.internalKvMapReplace(value, deep+1)
	}

	return "", false
}

//kv替换
func (c *Conf) kvMapReplace(value string) (string, bool) {
	return c.internalKvMapReplace(value, 0)
}

//获取指定namespace中的key，失败返回false
//TODO 测试Get不到Key时的返回值
func (c *Conf) Get(namespace string, key string) (string, bool) {
	value := c.ago.Get(key, agollo.WithNamespace(namespace))
	if value == "" {
		return "", false
	}
	if rtn, ok := c.kvMapReplace(value); ok {
		return rtn, true
	} else {
		return "", false
	}
}

//获取指定namespace中的key，失败返回默认值
//TODO 测试Get不到Key时的返回值
func (c *Conf) GetOrDefault(namespace string, key string, defaultValue string) string {
	value := c.ago.Get(key, agollo.WithNamespace(namespace))
	if value == "" {
		return defaultValue
	}
	if rtn, ok := c.kvMapReplace(value); ok {
		return rtn
	} else {
		return defaultValue
	}
}

//获取namespace
func (c *Conf) GetNamespace(namespace string) map[string]string {
	rtn := make(map[string]string)

	configs := c.ago.GetNameSpace(namespace)
	for k, v := range configs {
		if str, ok := v.(string); ok {
			if v, ok := c.kvMapReplace(str); ok {
				rtn[k] = v
			}
		}
	}

	return rtn
}

//初始化，如果初始化失败，则不可使用defaultConf
func init() {
	configServer, ok := args.Get("config_server")
	if !ok {
		return
	}
	appName, ok := args.Get("app_name")
	if !ok {
		return
	}
	idc, ok := args.Get("idc")
	if !ok {
		return
	}
	cacheFilePath := args.GetOrDefault("cache_file_path", os.Args[0]+".cache_file")

	if newConf, err := New(configServer, appName, idc, cacheFilePath, nil); err != nil {
		return
	} else {
		defaultConf = newConf
	}
}

//重新设置DefaultConf
func ResetDefaultConf(configServer string, appName string, idc string, cacheFilePath string, logger *log.Logger) error {
	if newConf, err := New(configServer, appName, idc, cacheFilePath, logger); err != nil {
		return err
	} else {
		defaultConf = newConf
	}
	return nil
}

func RefreshKvMap(kvMapper map[string]string) {
	defaultConf.RefreshKvMap(kvMapper)
}

func Get(namespace string, key string) (string, bool) {
	return defaultConf.Get(namespace, key)
}

func GetOrDefault(namespace string, key string, defaultValue string) string {
	return defaultConf.GetOrDefault(namespace, key, defaultValue)
}

func GetNamespace(namespace string) map[string]string {
	return defaultConf.GetNamespace(namespace)
}

func WatchNamespace(namespace string, handler func(oldCfgs map[string]string, newCfgs map[string]string)) {
	defaultConf.WatchNamespace(namespace, handler)
}

func Watch(namespace string, key string, handler func(oldCfg string, newCfg string)) {
	defaultConf.Watch(namespace, key, handler)
}
