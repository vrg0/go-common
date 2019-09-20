package registry

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"strconv"
	"sync"
	"time"
)

//consul注册器
type consulRegistry struct {
	masterClient         *api.Client   //主客户端
	watchClients         []*api.Client //监视器客户端
	servicesWatchParse   *watch.Plan   //服务列表监视器
	serviceWatchParseMap *sync.Map     //服务监视器
	observerMap          *sync.Map     //观察者映射
	serviceMap           ServiceMap    //服务映射
	serviceMapLock       *sync.RWMutex //服务映射锁
	//	serviceWatchParseMap map[string]*watch.Plan //服务监视器
}

//新建注册器
func newConsulRegistry() Registry {
	return &consulRegistry{
		masterClient:         nil,
		watchClients:         make([]*api.Client, 0),
		servicesWatchParse:   nil,
		serviceWatchParseMap: new(sync.Map),
		observerMap:          new(sync.Map),
		serviceMap:           make(map[string]*Service),
		serviceMapLock:       new(sync.RWMutex),
	}
}

const (
	RegisterRetryCount = 2
)

//内部初始化
func init() {
	if _, ok := newRegistryMap["consul"]; !ok {
		newRegistryMap["consul"] = newConsulRegistry
	}
}

//新建consul客户端
func newConsulClient(address string, dc string, waitTime time.Duration) *api.Client {
	config := api.DefaultConfig()
	config.Address = address
	config.Datacenter = dc
	config.WaitTime = waitTime
	client, _ := api.NewClient(config)
	return client
}

//初始化
//参数名         类型         格式		默认值
//cluster		[]string	ip:port		["127.0.0.1:8500"]
//dc			string		string		"dc1"
func (cr *consulRegistry) Init(options map[string]interface{}) error {
	//参数过滤
	cluster := []string{"127.0.0.1:8500"}
	dc := "dc1"
	if v, ok := options["cluster"]; ok {
		if _, ok := v.([]string); ok {
			cluster = v.([]string)
		} else {
			return errors.New("consul  cluster []string")
		}
	}
	if v, ok := options["dc"]; ok {
		if _, ok := v.(string); ok {
			dc = v.(string)
		} else {
			return errors.New("consul dc string")
		}
	}

	cr.masterClient = newConsulClient(cluster[0], dc, time.Second*1)
	for _, address := range cluster {
		cr.watchClients = append(cr.watchClients, newConsulClient(address, dc, time.Second*300))
	}

	//服务初始化
	parse, _ := watch.Parse(map[string]interface{}{"type": "services"})
	parse.Handler = cr.handlerServices
	cr.servicesWatchParse = parse
	go cr.watch(cr.servicesWatchParse, 0)

	return nil
}

//服务注册
//此函数会等待服务状态为passing，最多会阻塞10秒
func (cr *consulRegistry) Register(node *Node) error {
	node = DeepCopyNode(node)

	check := &api.AgentServiceCheck{
		HTTP:                           "http://" + node.Address + ":" + strconv.Itoa(node.Port) + "/tech/health/check",
		Interval:                       "5s",
		Timeout:                        "10s",
		DeregisterCriticalServiceAfter: "12h",
		CheckID:                        node.Id,
	}
	registration := api.AgentServiceRegistration{
		ID:      node.Id,
		Name:    node.ServiceName,
		Port:    node.Port,
		Address: node.Address,
		Check:   check,
		Meta:    node.Meta,
	}

	//可以重试一次
	for i := 0; i < RegisterRetryCount; i++ {
		if e := cr.masterClient.Agent().ServiceRegister(&registration); e != nil && i == RegisterRetryCount-1 {
			return e
		}
	}

	//轮训等待服务可用
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		if cr.nodeExistsAndPassing(node.ServiceName, node.Id) {
			return nil
		}
	}

	return errors.New("register timeout")
}

//服务注销
//此函数会等待节点删除，最多会阻塞10秒
func (cr *consulRegistry) Deregister(node *Node) error {
	node = DeepCopyNode(node)

	for i := 0; i < RegisterRetryCount; i++ {
		if e := cr.masterClient.Agent().ServiceDeregister(node.Id); e != nil && i == RegisterRetryCount-1 {
			return e
		}
	}

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		if !cr.nodeExists(node.ServiceName, node.Id) {
			return nil
		}
	}

	return errors.New("deregister timeout")
}

//获取服务映射，成功返回true，失败返回false
func (cr *consulRegistry) GetServiceMap() (ServiceMap, bool) {
	//读锁
	cr.serviceMapLock.RLock()
	defer cr.serviceMapLock.RUnlock()

	if cr.serviceMap == nil {
		return nil, false
	}

	return DeepCopyServiceMap(cr.serviceMap), true
}

//获取服务，成功返回true，失败返回false
func (cr *consulRegistry) GetService(serviceName string) (*Service, bool) {
	//读锁
	cr.serviceMapLock.RLock()
	defer cr.serviceMapLock.RUnlock()

	if service, ok := cr.serviceMap[serviceName]; ok {
		return DeepCopyService(service), true
	} else {
		return nil, false
	}
}

//设置观察者
func (cr *consulRegistry) SetObserver(observer Observer) error {
	cr.observerMap.Store(observer, observer)
	return nil
}

//节点存在并且状态为passing
func (cr *consulRegistry) nodeExistsAndPassing(serviceName string, nodeId string) bool {
	//读锁
	cr.serviceMapLock.RLock()
	defer cr.serviceMapLock.RUnlock()

	if service, ok := cr.serviceMap[serviceName]; ok {
		if node, ok := service.NodeMap[nodeId]; ok && node.Status == Passing {
			return true
		}
	}
	return false
}

//判断节点是否存在
func (cr *consulRegistry) nodeExists(serviceName string, nodeId string) bool {
	//读锁
	cr.serviceMapLock.RLock()
	defer cr.serviceMapLock.RUnlock()

	if service, ok := cr.serviceMap[serviceName]; ok {
		if _, ok := service.NodeMap[nodeId]; ok {
			return true
		}
	}

	return false
}

//监控服务列表的变化
func (cr *consulRegistry) handlerServices(_ uint64, data interface{}) {
	services := data.(map[string][]string)

	//锁
	cr.serviceMapLock.Lock()
	defer cr.serviceMapLock.Unlock()

	//监控全部服务
	for service, tags := range services {
		//如果服务不存在，则创建服务、监听服务
		parse, _ := watch.Parse(map[string]interface{}{"type": "service", "service": service})
		parse.Handler = cr.handlerService
		if _, ok := cr.serviceWatchParseMap.LoadOrStore(service, parse); !ok {
			go cr.watch(parse, 0)
		}

		//获取标签
		tagMap := make(map[string]struct{})
		for _, tag := range tags {
			tagMap[tag] = struct{}{}
		}

		//更新服务
		newService := new(Service)
		if oldService, ok := cr.serviceMap[service]; ok {
			newService = DeepCopyService(oldService)
		} else {
			newService = NewService(service)
		}
		newService.TagMap = tagMap
		cr.serviceMap[service] = newService
	}

	//处理被删除的服务
	cr.serviceWatchParseMap.Range(func(key interface{}, value interface{}) bool {
		serviceName := key.(string)
		parse := value.(*watch.Plan)
		if _, ok := services[serviceName]; !ok {
			//删除服务
			parse.Stop()
			cr.serviceWatchParseMap.Delete(serviceName)
			delete(cr.serviceMap, serviceName)

			//执行观察者
			cr.observerMap.Range(func(_ interface{}, value interface{}) bool {
				observer := value.(Observer)
				//清空服务中的节点和标签
				observer.UpdateNodes(NewService(serviceName))

				//删除服务
				observer.DeleteService(serviceName)
				return true
			})
		}
		return true
	})
}

//观察服务的变化
func (cr *consulRegistry) handlerService(_ uint64, data interface{}) {
	nodeList := data.([]*api.ServiceEntry)
	//确保至少有一个node
	if len(nodeList) == 0 {
		return
	}
	//获取服务名称
	serviceName := nodeList[0].Service.Service
	//忽略掉 consul 服务
	if serviceName == "consul" {
		return
	}

	cr.serviceMapLock.Lock()
	defer cr.serviceMapLock.Unlock()

	nodeMap := make(map[string]*Node)
	for _, node := range nodeList {
		newNode := NewNode(serviceName, node.Service.ID, node.Service.Address, node.Service.Port)
		newMeta := make(map[string]string)
		for k, v := range node.Service.Meta {
			newMeta[k] = v
		}
		newNode.Meta = newMeta
		for _, health := range node.Checks {
			if node.Service.ID == health.ServiceID {
				if health.Status == api.HealthPassing {
					newNode.Status = Passing
				} else {
					newNode.Status = Critical
				}
			}
		}
		nodeMap[node.Service.ID] = newNode
	}
	cr.serviceMap[serviceName].NodeMap = nodeMap

	cr.observerMap.Range(func(_ interface{}, value interface{}) bool {
		observer := value.(Observer)
		service := DeepCopyService(cr.serviceMap[serviceName])
		observer.UpdateNodes(service)
		return true
	})
}

//监控
func (cr *consulRegistry) watch(parse *watch.Plan, sentry int) {
	defer func() {
		if e := recover(); e != nil {
			time.Sleep(time.Second * 1)
			if sentry < len(cr.watchClients)-1 {
				sentry++
			} else {
				sentry = 0
			}
			cr.watch(parse, sentry)
		}
	}()

	_ = parse.RunWithClientAndLogger(cr.watchClients[sentry], nil)
}
