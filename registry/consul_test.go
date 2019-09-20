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

func TestConsulRegistry_Register(t *testing.T) {
	TestConsulRegistry_Init(t)
	http.HandleFunc("/tech/health/check", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, "Pong!\n")
	})
	go func() { _ = http.ListenAndServe("192.168.29.88:8989", nil) }()

	node := NewNode("aaaa", "192.168.29.88:8989", "192.168.29.88", 8989)
	if e := testConsulRegister.Register(node); e != nil {
		t.Error(e)
	}
}

func TestConsulRegistry_Deregister(t *testing.T) {
	TestConsulRegistry_Register(t)

	node := NewNode("aaaa", "192.168.29.88:8989", "192.168.29.88", 8989)
	if e := testConsulRegister.Deregister(node); e != nil {
		t.Error(e)
	}
}

func TestConsulRegistry_GetService(t *testing.T) {
	TestConsulRegistry_Register(t)
	if service, ok := testConsulRegister.GetService("aaaa"); !ok {
		t.Error("get service fatal")
	} else {
		for _, n := range service.NodeMap {
			fmt.Println(n)
		}
	}
}

func TestConsulRegistry_GetServiceMap(t *testing.T) {
	TestConsulRegistry_Register(t)
	if serviceMap, ok := testConsulRegister.GetServiceMap(); !ok {
		t.Error("get service map fatal")
	} else {
		for _, v := range serviceMap {
			for _, n := range v.NodeMap {
				fmt.Println(n)
			}
			fmt.Println("------------")
		}
	}
}

type testObserver struct{}

func (testObserver) DeleteService(serviceName string) {
	fmt.Println(serviceName)
}
func (testObserver) UpdateNodes(service *Service) {
	for _, node := range service.NodeMap {
		fmt.Println(node)
	}
}

func TestConsulRegistry_SetObserver(t *testing.T) {
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
