package registry

import (
	"testing"
)

type testRegistry struct{}

func newTestRegistry() Registry {
	return &testRegistry{}
}

func init() {
	if _, ok := newRegistryMap["test"]; !ok {
		newRegistryMap["test"] = newTestRegistry
	}
}

func (testRegistry) Init(options map[string]interface{}) error {
	return nil
}

func (testRegistry) Register(node *Node) error {
	return nil
}

func (testRegistry) Deregister(node *Node) error {
	return nil
}

func (testRegistry) GetServiceMap() (ServiceMap, bool) {
	return nil, true
}

func (testRegistry) GetService(serviceName string) (*Service, bool) {
	return nil, true
}

func (testRegistry) SetObserver(observer Observer) error {
	return nil
}

func defaultRegistryClear() {
	defaultRegistry = nil
}

func TestInit(t *testing.T) {
	defaultRegistryClear()
	if e := Init("testXXX", nil); e == nil {
		t.Error("init error")
	}

	defaultRegistryClear()
	if e := Init("test", nil); e != nil {
		t.Error(e)
	}
}

func TestNew(t *testing.T) {
	if _, e := New("test"); e != nil {
		t.Error(e)
	}

	if _, e := New("testXXX"); e == nil {
		t.Error("registry error")
	}
}

func TestDeepCopyNode(t *testing.T) {
	oldNode := &Node{
		ServiceName: "ttt",
		Id:          "123",
		Address:     "1.1.1.1",
		Port:        12345,
		Meta:        make(map[string]string),
		Status:      Passing,
	}

	newNode := DeepCopyNode(oldNode)

	if &oldNode.Meta == &newNode.Meta {
		t.Error("node.Meta copy fatal")
	}
}

func TestDeepCopyService(t *testing.T) {
	oldNode := &Node{
		ServiceName: "ttt",
		Id:          "123",
		Address:     "1.1.1.1",
		Port:        12345,
		Meta:        make(map[string]string),
		Status:      Passing,
	}

	oldService := &Service{
		Name:    oldNode.ServiceName,
		NodeMap: map[string]*Node{oldNode.Id: oldNode},
		TagMap:  map[string]struct{}{},
	}

	newService := DeepCopyService(oldService)

	if &oldService.TagMap == &newService.TagMap {
		t.Error("service.TagMap fatal")
	}

	if &oldService.NodeMap == &newService.NodeMap {
		t.Error("service.NodeMap fatal")
	}

	if oldService.NodeMap[oldNode.Id] == newService.NodeMap[oldNode.Id] {
		t.Error("service.NodeMap[] fatal")
	}
}

func TestDeepCopyServiceMap(t *testing.T) {
	oldNode := &Node{
		ServiceName: "ttt",
		Id:          "123",
		Address:     "1.1.1.1",
		Port:        12345,
		Meta:        make(map[string]string),
		Status:      Passing,
	}

	oldService := &Service{
		Name:    oldNode.ServiceName,
		NodeMap: map[string]*Node{oldNode.Id: oldNode},
		TagMap:  map[string]struct{}{},
	}

	oldServiceMap := ServiceMap{oldService.Name: oldService}

	newServiceMap := DeepCopyServiceMap(oldServiceMap)

	if &oldServiceMap == &newServiceMap {
		t.Error("ServiceMap fatal")
	}

	if oldServiceMap[oldService.Name] == newServiceMap[oldService.Name] {
		t.Error("ServiceMap[] fatal")
	}
}

func TestDeregister(t *testing.T) {
	TestInit(t)
	if e := Deregister(nil); e != nil {
		t.Error(e)
	}
}

func TestGetService(t *testing.T) {
	TestInit(t)
	if _, ok := GetService("ttt"); !ok {
		t.Error("GetService")
	}
}

func TestGetServiceMap(t *testing.T) {
	TestInit(t)
	if _, ok := GetServiceMap(); !ok {
		t.Error("GetServiceMap")
	}
}

func TestSetObserver(t *testing.T) {
	TestInit(t)
	if e := SetObserver(nil); e != nil {
		t.Log(e)
	}
}

func TestRegister(t *testing.T) {
	TestInit(t)
	if e := Register(nil); e != nil {
		t.Error(e)
	}
}

func TestNewNode(t *testing.T) {
	node := NewNode("ttt", "1", "1.1.1.1", 1111)
	if node == nil {
		t.Error("node")
	}
	if node.ServiceName != "ttt" {
		t.Error("node.ServiceName")
	}
	if node.Id != "1" {
		t.Error("node.Id")
	}
	if node.Address != "1.1.1.1" {
		t.Error("node.Address")
	}
	if node.Port != 1111 {
		t.Error("node.Port")
	}
	if node.Status != Unknown {
		t.Error("node.Status")
	}
	if node.Meta == nil || len(node.Meta) != 0 {
		t.Error("node.Meta")
	}
}

func TestNewService(t *testing.T) {
	service := NewService("ttt")
	if service == nil {
		t.Error("service")
	}
	if service.Name != "ttt" {
		t.Error("service.ServiceName")
	}
	if service.NodeMap == nil || len(service.NodeMap) != 0 {
		t.Error("service.NodeMap")
	}
	if service.TagMap == nil || len(service.TagMap) != 0 {
		t.Error("service.TagMap")
	}
}
