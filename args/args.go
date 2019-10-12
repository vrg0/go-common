package args

/**
 * 从[环境变量]或[命令行参数]中获取配置
 * key不区分大小写，value区分大小写
 * [命令行参数]会覆盖掉[环境变量]中的配置
 */

import (
	"os"
	"strings"
)

var (
	argMap = make(map[string]string)
)

func argUnMarshal(arg string) (string, string, bool) {
	kv := strings.Split(arg, "=")
	if len(kv) != 2 {
		return "", "", false
	}

	key := kv[0]
	if key[0] != '-' {
		return "", "", false
	}
	if key[1] == '-' {
		key = key[2:]
	} else {
		key = key[1:]
	}
	key = strings.ToLower(key)

	value := kv[1]

	return key, value, true
}

func envUnMarshal(env string) (string, string) {
	kv := strings.Split(env, "=")
	return strings.ToLower(kv[0]), kv[1]
}

func init() {
	//解析环境变量
	for _, env := range os.Environ() {
		k, v := envUnMarshal(env)
		argMap[k] = v
	}

	//解析命令行参数
	for _, arg := range os.Args {
		if k, v, ok := argUnMarshal(arg); ok {
			argMap[k] = v
		}
	}
}

//获取参数，如果失败则返回默认值
func GetOrDefault(key string, defaultValue string) string {
	key = strings.ToLower(key)
	if value, ok := argMap[key]; ok {
		return value
	} else {
		return defaultValue
	}
}

//获取参数，成功返回(value, true)，失败返回("", false)
func Get(key string) (string, bool) {
	key = strings.ToLower(key)
	if value, ok := argMap[key]; ok {
		return value, true
	} else {
		return "", false
	}
}
