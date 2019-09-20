package registry

import (
	"errors"
)

//节点
type Node struct {
	ServiceName string            //服务名称
	Id          string            //节点编号
	Address     string            //地址
	Port        int               //端口
	Meta        map[string]string //元数据
	Status      string            //状态
}

//服务
type Service struct {
	Name    string              //服务名称
	NodeMap map[string]*Node    //节点映射
	TagMap  map[string]struct{} //标签
}

//服务映射
type ServiceMap map[string]*Service

//观察者
type Observer interface {
	//删除节点事件
	DeleteService(serviceName string)

	//更新节点事件
	UpdateNodes(service *Service)
}

//注册器
type Registry interface {
	//初始化
	Init(options map[string]interface{}) error

	//服务注册
	Register(node *Node) error

	//服务注销
	Deregister(node *Node) error

	//获取服务映射，成功返回true，失败返回false
	GetServiceMap() (ServiceMap, bool)

	//获取服务，成功返回true，失败返回false
	GetService(serviceName string) (*Service, bool)

	//设置观察者
	SetObserver(observer Observer) error
}

var (
	//注册器映射
	newRegistryMap = make(map[string]func() Registry, 0)

	//默认注册器
	defaultRegistry Registry = nil
)

const (
	Unknown  = "unknown"
	Passing  = "passing"
	Critical = "critical"
)

//新建注册器
func New(name string) (Registry, error) {
	if newRegistry, ok := newRegistryMap[name]; ok {
		return newRegistry(), nil
	} else {
		return nil, errors.New("not found registry")
	}
}

//新建节点
func NewNode(serviceName string, id string, address string, port int) *Node {
	return &Node{
		ServiceName: serviceName,
		Id:          id,
		Address:     address,
		Port:        port,
		Status:      Unknown,
		Meta:        make(map[string]string, 0),
	}
}

//新建服务
func NewService(name string) *Service {
	return &Service{
		Name:    name,
		NodeMap: make(map[string]*Node, 0),
		TagMap:  make(map[string]struct{}, 0),
	}
}

//深拷贝节点
func DeepCopyNode(node *Node) *Node {
	rtn := NewNode(node.ServiceName, node.Id, node.Address, node.Port)
	rtn.Status = node.Status
	for k, v := range node.Meta {
		rtn.Meta[k] = v
	}
	return rtn
}

//深拷贝服务
func DeepCopyService(service *Service) *Service {
	rtn := NewService(service.Name)
	for _, node := range service.NodeMap {
		rtn.NodeMap[node.Id] = DeepCopyNode(node)
	}
	for k, v := range service.TagMap {
		rtn.TagMap[k] = v
	}
	return rtn
}

//深拷贝服务列表
func DeepCopyServiceMap(serviceMap ServiceMap) ServiceMap {
	rtn := make(ServiceMap)
	for name, service := range serviceMap {
		rtn[name] = DeepCopyService(service)
	}
	return rtn
}

//初始化
func Init(name string, options map[string]interface{}) error {
	//懒汉单例
	if defaultRegistry != nil {
		return nil
	}

	if newRegistry, e := New(name); e != nil {
		return e
	} else {
		defaultRegistry = newRegistry
	}

	return defaultRegistry.Init(options)
}

//服务注册
func Register(node *Node) error {
	return defaultRegistry.Register(node)
}

//服务注销
func Deregister(node *Node) error {
	return defaultRegistry.Deregister(node)
}

//获取服务映射
func GetServiceMap() (ServiceMap, bool) {
	return defaultRegistry.GetServiceMap()
}

//获取服务
func GetService(serviceName string) (*Service, bool) {
	return defaultRegistry.GetService(serviceName)
}

//设置观察者
func SetObserver(observer Observer) error {
	return defaultRegistry.SetObserver(observer)
}
