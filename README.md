# 使用说明

## args模块

args模块在init的过程中会读取[环境变量]和[命令行参数]的kv对，其中key不区分大小写。

当[环境变量]和[命令行参数]中的key重复时，默认读取[命令行参数]中的key。

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

当env != "dev"时，日志默认打印到os.argv[0]+".log"，日志登记默认为info。

```go
//打印日志
logger.Info("info")

//重新设置默认日志
logger.ResetDefaultLogger("/log/path", zapcore.DebugLevel)

//日志hook
logger.SetWatchFunc(func(data []byte) {
  //fmt.Println(string(data))
})

//获取标准库日志对象
sl := logger.GetStandardLogger()

//日志hook & 标准库日志
sl := logger.GetStandardLogger()
sl.SetPrefix("_LogHook_")
logger.SetWatchFunc(func(data []byte) {
  if strings.Contains(string(data), "_LogHook")  {
    //fmt.Println(string(data))
  }
})
sl.Print("log hook")
```

