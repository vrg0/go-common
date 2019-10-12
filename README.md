# 使用说明

## args模块

args模块在init的过程中会读取[环境变量]和[命令行参数]的kv对，其中key不区分大小写。

当[环境变量]和[命令行参数]中的key重复时，默认读取[命令行参数]中的key。

命令行参数格式：key=value | -key=value | --key=value

```go
//获取参数，如果失败则返回默认值
value := args.GetOrDefault("key", "defaultValue")

//获取参数，如果失败则返回("", false)
value, ok := args.Get("key")
```

## logger模块

logger模块在init过程中会使用args模块读取env和log_path参数。

当env获取失败时，默认值为"dev"

当env =="dev"时，日志默认打印到stdout，日志等级默认为debug。

当env != "dev"时，日志默认打印到os.Argv[0]+".log"，日志等级默认为info。

```go
//打印日志
logger.Info("info")

//重新设置默认日志
logger.ResetDefaultLogger("/log/path", zapcore.DebugLevel)

//日志hook
logger.SetHookFunc(func(data []byte) {
  //fmt.Println(string(data))
})

//获取标准库日志对象
sl := logger.GetStandardLogger()

//日志hook & 标准库日志
sl := logger.GetStandardLogger()
sl.SetPrefix("_LogHook_ ")
logger.SetHookFunc(func(data []byte) {
  if strings.Contains(string(data), "_LogHook_")  {
    //fmt.Println(string(data))
  }
})
sl.Print("log hook")
```

## conf模块

logger模块在init过程中会使用args模块读取config_serve、app_name、idc、cache_file_path参数。

config_serve：apollo服务器的地址，如：localhost:8080 | agollo.xxx.com

app_name：对应apollo中的AppId

idc：对应apollo中的Cluster

cache_file_path：指定apollo缓存文件的路径，默认为os.Argv[0]+".cache_file"

当config_serve或app_name或idc读取失败时，则初始化失败，此时不可使用conf模块的defaultConf。

```go
//获取配置，失败返回("", false)
value, ok := Get("namespace", "xxx")

//获取配置，失败返回"defaultValue"
value := Get("namespace", "xxx", "defaultValue")

//获取namespace
kvMap := GetNamespace("namespace")

//kv映射功能
//把配置中的 "${key}" 映射成value
//${key}的嵌套，如：${hello}值为world，${AAA}值为ll, 则${he${AAA}o}值为world
//${key}的嵌套最多支持8层
RefreshKvMap(map[string]string{"123":"abc"})

//监控某个namespace的变化，当发生变化后调用回调函数
WatchNamespace("application",
  func(oldCfgs map[string]string, newCfgs map[string]string) {
  	//pass
})

//监控某个key的变化，当发生变化后调用回调函数
Watch("application", "key",
   func(oldCfg string, newCfg string)) {
  	//pass
})
```