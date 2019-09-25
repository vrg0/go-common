package registry

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

var (
	testConsulRegister Registry = nil
)

//正常初始化
func TestConsulRegistry_Init(t *testing.T) {
	if testConsulRegister == nil {
		dc := "dc1"
		cluster := []string{"127.0.0.1:8500"}
		r := newConsulRegistry()
		if e := r.Init(map[string]interface{}{"dc": dc, "cluster": cluster}); e != nil {
			t.Error(e)
		} else {
			testConsulRegister = r
		}
	}
}

//默认初始化
func TestConsulRegistry_Init2(t *testing.T) {
	r := newConsulRegistry()
	if e := r.Init(map[string]interface{}{}); e != nil {
		t.Error(e)
	}
}

//服务注册，正常
func TestConsulRegistry_Register(t *testing.T) {
	//开启http-service
	TestConsulRegistry_Init(t)
	http.HandleFunc("/tech/health/check", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "Pong!\n")
	})
	go func() { _ = http.ListenAndServe("192.168.29.88:8989", nil) }()

	//注册服务，预期会阻塞运行，直到 健康检查的状态为 passing。预期返回nil
	node := NewNode("aaaa", "192.168.29.88:8989", "192.168.29.88", 8989)
	if e := testConsulRegister.Register(node); e != nil {
		t.Error(e)
	}
}

//服务注册，异常
func TestConsulRegistry_Register2(t *testing.T) {
	TestConsulRegistry_Init(t)

	//健康检查的状态为 passing。 预期所以会阻塞10秒，返回 timeout
	node := NewNode("aaaa", "192.168.29.88:9898", "192.168.29.88", 8989)
	if e := testConsulRegister.Register(node); e != nil {
		if e.Error() != "register timeout" {
			t.Error(e)
		}
	}
}

//服务注销
func TestConsulRegistry_Deregister(t *testing.T) {
	TestConsulRegistry_Register(t)

	//会阻塞执行，直到服务节点被删除，或者整个服务被删除。预期返回nil
	node := NewNode("aaaa", "192.168.29.88:8989", "192.168.29.88", 8989)
	if e := testConsulRegister.Deregister(node); e != nil {
		t.Error(e)
	}
}

//获取service
func TestConsulRegistry_GetService(t *testing.T) {
	//预期：获取到aaaa服务有一个节点，id为192.168.29.88:8989
	TestConsulRegistry_Register(t)
	if service, ok := testConsulRegister.GetService("aaaa"); !ok {
		t.Error("get service fatal")
	} else {
		if len(service.NodeMap) != 1 {
			t.Error("service.NodeMap is not eq 1")
		}
		if _, ok := service.NodeMap["192.168.29.88:8989"]; !ok {
			t.Error("service.NodeMap is not exists")
		}
	}
}

//获取 service Map
func TestConsulRegistry_GetServiceMap(t *testing.T) {
	//预期：aaaa服务存在
	TestConsulRegistry_Register(t)
	if serviceMap, ok := testConsulRegister.GetServiceMap(); !ok {
		t.Error("get service map fatal")
	} else {
		if _, ok := serviceMap["aaaa"]; !ok {
			t.Error("aaaa is not exists")
		}

		for _, v := range serviceMap {
			for _, n := range v.NodeMap {
				fmt.Println(n)
			}
			fmt.Println("------------")
		}
	}
}

//观察者接口的实现
type testObserver struct{}
func (testObserver) DeleteService(serviceName string) {
	fmt.Println(serviceName)
}
func (testObserver) UpdateNodes(service *Service) {
	for _, node := range service.NodeMap {
		fmt.Println(node)
	}
}

//设置观察者
func TestConsulRegistry_SetObserver(t *testing.T) {
	//预期：注册节点成功后，打印aaaa相关的信息
	TestConsulRegistry_Register(t)
	o := new(testObserver)
	if e := testConsulRegister.SetObserver(o); e != nil {
		t.Error(e)
	}

	node := NewNode("aaaa", "192.168.29.88:8989", "192.168.29.88", 8989)
	if e := testConsulRegister.Deregister(node); e != nil {
		t.Error(e)
	}
}
