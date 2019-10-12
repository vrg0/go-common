package health

import (
	"encoding/json"
	"errors"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/vrg0/go-common/util"
	"runtime"
	"sync"
	"time"
)

// 块Msg
type ChunkMsg struct {
	AppName   string     `json:"appName"`   //服务名称
	Ip        string     `json:"ip"`        //服务IP
	EndPoints []EndPoint `json:"endPoints"` //点
}

//点
type EndPoint struct {
	Health    Health      `json:"health"`    //健康信息
	Metrics   MetricsInfo `json:"metrics"`   //基本信息
	Timestamp int64       `json:"timestamp"` //时间戳
}

//常规a信息
type MetricsInfo struct {
	Mem                        int     `json:"mem"`                          //内存总量
	MemFree                    int     `json:"mem.free"`                     //内存可用
	Processors                 int     `json:"processors"`                   //CPU核心数
	SystemloadAverage          float64 `json:"systemload.average"`           //CPU负载、负载/核心数 = CPU使用量
	InstanceUptime             int     `json:"instance.uptime"`              //程序运行时间
	Uptime                     int     `json:"uptime"`                       //系统运行时间
	HeapCommitted              int     `json:"heap.committed"`               //堆提交
	HeapInit                   int     `json:"heap.init"`                    //堆初始化
	HeapUsed                   int     `json:"heap.used"`                    //堆已用
	Heap                       int     `json:"heap"`                         //堆总量
	NonheapCommitted           int     `json:"nonheap.committed"`            //非堆提交
	NonheapInit                int     `json:"nonheap.init"`                 //非堆初始化
	NonheapUsed                int     `json:"nonheap.used"`                 //非堆已用
	Nonheap                    int     `json:"nonheap"`                      //非堆可用
	ThreadsPeak                int     `json:"threads.peak"`                 //峰值线程数
	ThreadsDaemon              int     `json:"threads.daemon"`               //正在运行的线程数
	ThreadsTotalStarted        int     `json:"threads.totalStarted"`         //总过创建的线程数
	Threads                    int     `json:"threads"`                      //线程数（废弃）
	Classes                    int     `json:"classes"`                      //类总数
	ClassesLoaded              int     `json:"classes.loaded"`               //类加载数ok
	ClassesUnloaded            int     `json:"classes.unloaded"`             //类回收数ok
	GcParnewCount              int     `json:"gc.parnew.count"`              //GC次数
	GcParnewTime               int     `json:"gc.parnew.time"`               //GC时间
	GcConcurrentmarksweepCount int     `json:"gc.concurrentmarksweep.count"` //GC次数
	GcConcurrentmarksweepTime  int     `json:"gc.concurrentmarksweep.time"`  //GC时间
}

//健康信息
type Health struct {
	Status  Status `json:"status"`  //服务的状态
	Details Detail `json:"details"` //细节（硬盘、缓存、数据库）
}

//状态
type Status struct {
	Code        string `json:"code"`        //状态码， UP为1， 其他为0
	Description string `json:"description"` //细节
}

//磁盘信息
type Detail struct {
	DiskSpace DiskInfo  `json:"diskSpace"` //硬盘信息
	Redis     RedisInfo `json:"redis"`     //缓存信息
	Db        DBInfo    `json:"db"`        //数据库信息
}

//磁盘信息
type DiskInfo struct {
	Details DiskDetail `json:"details"` //硬盘细节
	Status  Status     `json:"status"`  //硬盘状态
}

//磁盘细节
type DiskDetail struct {
	Total     int64 `json:"total"`     //硬盘总量
	Free      int64 `json:"free"`      //硬盘可用
	Threshold int64 `json:"threshold"` //门限（为使用）
}

//redis信息
type RedisInfo struct {
	Status Status `json:"status"` //缓存状态
}

//数据信息
type DBInfo struct {
	Details map[string]DBDetail `json:"details"` //数据库细节
	Status  Status              `json:"status"`  //数据库状态
}

//数据库
type DBDetail struct {
	Status  Status       `json:"status"`  //状态
	Details DBDetailInfo `json:"details"` //数据库细节
}

type DBDetailInfo struct {
	Database string `json:"database"` //数据库名称
}

//添加Endpoint
func (c *ChunkMsg) addEndPoint(endpoint EndPoint) {
	c.EndPoints = append(c.EndPoints, endpoint)
}

//转化成Json
func (c *ChunkMsg) ToJson() []byte {
	rtn, _ := json.Marshal(c)
	return rtn
}

var (
	defaultStatus    = Status{Code: "UP", Description: ""} //默认状态
	startTime        = time.Now()                          //程序启动时间
	callbackMap      = new(sync.Map)                       //回调函数列表
	initOnce         = new(sync.Once)                      //init单例
	internalAppName  = ""                                  //应用名称
	internalIp       = ""                                  //应用IP
	internalDiskPath = ""                                  //要监控的硬盘
)

//diskPath：要监测的磁盘的挂载路径
func Start(appName string, diskPath string) error {
	//初始化
	internalAppName = appName
	ip, err := util.LocalIp()
	if err != nil {
		return errors.New("get local ip fatal")
	} else {
		internalIp = ip
	}
	internalDiskPath = diskPath

	//单例
	initOnce.Do(internalInit)

	return nil
}

func getEndpoint() EndPoint {
	diskFree, diskTotal := getDisk(internalDiskPath)
	runtimeInfo := runtime.MemStats{}
	runtime.ReadMemStats(&runtimeInfo)

	health := Health{
		Status: defaultStatus,
		Details: Detail{
			DiskSpace: DiskInfo{
				Details: DiskDetail{
					Free:      diskFree,
					Total:     diskTotal,
					Threshold: 0,
				},
				Status: defaultStatus,
			},
			Redis: RedisInfo{
				Status: defaultStatus,
			},
			Db: DBInfo{
				Status:  defaultStatus,
				Details: make(map[string]DBDetail),
			},
		},
	}

	memFree, memTotal := getMem()
	cpuCore, cpuSum := getCpu()
	cpuSum /= 100
	metricsInfo := MetricsInfo{
		Uptime:            (int(time.Now().Unix()) - int(startTime.Unix())) * 1000,
		InstanceUptime:    0,
		Mem:               memTotal,
		MemFree:           memFree,
		Heap:              int(runtimeInfo.HeapSys / 1024),
		HeapUsed:          int(runtimeInfo.HeapAlloc / 1024),
		ClassesLoaded:     0,
		ClassesUnloaded:   0,
		NonheapInit:       0,
		NonheapUsed:       int(runtimeInfo.StackInuse / 1024),
		ThreadsPeak:       runtime.NumGoroutine(),
		SystemloadAverage: cpuSum,
		Processors:        cpuCore,
	}

	rtn := EndPoint{
		Health:    health,
		Metrics:   metricsInfo,
		Timestamp: time.Now().UnixNano() / 1e6,
	}
	return rtn
}

func getMem() (int, int) {
	stat, e := mem.VirtualMemory()
	if e != nil {
		return 0, 0
	} else {
		return int(stat.Total-stat.Used) / 1024, int(stat.Total) / 1024
	}
}

func getDisk(path string) (int64, int64) {
	stat, _ := disk.Usage(path)
	if stat != nil {
		return int64(stat.Free), int64(stat.Total)
	} else {
		return 0, 0
	}
}

func getCpu() (int, float64) {
	percentage, e := cpu.Percent(0, false)
	if e != nil {
		return 0, 0
	}
	sum := float64(0)
	for _, v := range percentage {
		sum += v
	}
	return len(percentage), sum
}

//内部初始化
func internalInit() {
	go func() {
		//防止panic
		defer func() {
			if e := recover(); e != nil {
				time.Sleep(time.Second * 1)
				internalInit()
			}
		}()

		//每间隔1秒，进行一次数据采集。
		//没间隔5秒，调用一次回调函数。
		count := 0
		tick := time.NewTicker(time.Second * 1)
		chunk := &ChunkMsg{
			AppName:   internalAppName,
			Ip:        internalIp,
			EndPoints: make([]EndPoint, 0),
		}
		for {
			select {
			case <-tick.C:
				//添加 endpoint
				endpoint := getEndpoint()
				chunk.EndPoints = append(chunk.EndPoints, endpoint)

				count++
				if count == 5 {
					count = 0
					//调用Callback
					callbackMap.Range(func(_ interface{}, value interface{}) bool {
						handler := value.(func(*ChunkMsg))
						newChunk := &ChunkMsg{}
						if e := util.DeepCopy(newChunk, chunk); e != nil {
							return false
						}
						handler(newChunk)
						return true
					})

					//清空ChunkMsg
					chunk.EndPoints = make([]EndPoint, 0)
				}
			}
		}
	}()
}

func SetCallback(handler func(*ChunkMsg)) {
	callbackMap.Store(&handler, handler)
}
