package conf

import (
	"github.com/shima-park/agollo"
	"log"
	"strings"
	"time"
)

/**
 * 使用conf模块必须进行初始化
 *
 * kv映射功能
 * 把配置中的 "${key}" 映射成value
 * 支${key}的嵌套，如：${hello}值为world，${AAA}值为ll, 则${he${AAA}o}值为world
 * ${key}的嵌套最多支持8层
 */

var (
	kvMap                      = make(map[string]string)
	internalLogger *log.Logger = nil
)

const (
	kvMapReplaceDeep  = 8
	placeHolderPrefix = "${"
	placeHolderSuffix = "}"
)

func internalLogPrint(v ...interface{}) {
	if internalLogger != nil {
		internalLogger.Print(v...)
	}
}

//domainName: 域名，如 baidu.com
//logger: conf模块的日志，传nil表示不记录日志
//cacheFilePath: 缓存文件路径，当路径为空时，启用默认路径 ./agollo.cache
func Init(appID string, domainName string, env string, cluster string, cacheFilePath string, logger *log.Logger) error {
	//缓存文件默认路径
	if cacheFilePath == "" {
		cacheFilePath = "./agollo.cache"
	}
	internalLogger = logger

	configServerURL := "apollo-" + env + "." + domainName + "/"
	if e := agollo.Init(configServerURL, appID,
		agollo.Cluster(cluster),                   //集群名称(idc)
		agollo.BackupFile(cacheFilePath),          //缓存文件的路径
		agollo.FailTolerantOnBackupExists(),       //从apollo service读取配置失败时，从缓存文件中读取配置
		agollo.AutoFetchOnCacheMiss(),             //当缓存中找不到namespace时，自动从apollo server拉取namespace
		agollo.LongPollerInterval(time.Second*10), //从apollo server更新数据的轮训时间
	); e != nil {
		return e
	}

	errLogChan := agollo.Start() //开启一个协程，从apollo server更新数据
	go func() {
		defer func() {
			if e := recover(); e != nil {
				return
			}
		}()
		for {
			select {
			case e := <-errLogChan:
				internalLogPrint(e)
			}
		}
	}()

	return nil
}

//刷新kvMap
func RefreshKvMap(kvMapper map[string]string) error {
	newKvMap := make(map[string]string)
	for k, v := range kvMapper {
		newKvMap[k] = v
	}
	kvMap = newKvMap
	return nil
}

func internalKvMapReplace(value string, deep int) (string, bool) {
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
	if replaceValue, ok := kvMap[innerStr]; ok {
		value = strings.ReplaceAll(value, placeHolderPrefix+innerStr+placeHolderSuffix, replaceValue)
		return internalKvMapReplace(value, deep+1)
	}

	return "", false
}

//kv替换
func kvMapReplace(value string) (string, bool) {
	return internalKvMapReplace(value, 0)
}

//获取指定namespace中的key
func Get(namespace string, key string) string {
	value := agollo.Get(key, agollo.WithNamespace(namespace))
	if rtn, ok := kvMapReplace(value); ok {
		return rtn
	} else {
		return ""
	}
}

//获取application namespace中的key
func GetWithDefault(key string) string {
	if rtn, ok := kvMapReplace(agollo.Get(key)); ok {
		return rtn
	} else {
		return ""
	}
}

//获取namespace
func GetNamespace(namespace string) map[string]string {
	rtn := make(map[string]string)

	configs := agollo.GetNameSpace(namespace)
	for k, v := range configs {
		if str, ok := v.(string); ok {
			if v, ok := kvMapReplace(str); ok {
				rtn[k] = v
			}
		}
	}

	return rtn
}
