package consul_kv

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"testing"
	"time"
)

var consulTestKvObj *ConsulKV = nil
var DC = "dc1"
var CLUSTER = []string{"127.0.0.1:8500"}

func TestNew(t *testing.T) {
	if consulTestKvObj == nil {
		consulTestKvObj = New(DC, CLUSTER)
		if consulTestKvObj == nil {
			t.Error("new ConsulKV err")
		}
	}
}

func TestConsulKV_SetValue(t *testing.T) {
	TestNew(t)
	if e := consulTestKvObj.SetValue("testXXX", "xxx"); e != nil {
		t.Error(e)
	}
}

func TestConsulKV_DelValue(t *testing.T) {
	TestConsulKV_SetValue(t)
	if e := consulTestKvObj.DelValue("testXXX"); e != nil {
		t.Error(e)
	}
}

func TestConsulKV_GetValue(t *testing.T) {
	TestConsulKV_SetValue(t)
	if v, e := consulTestKvObj.GetValue("testXXX"); e != nil {
		t.Error(e)
	} else {
		t.Log(v)
	}
}

func TestConsulKV_List(t *testing.T) {
	TestConsulKV_SetValue(t)
	if v, e := consulTestKvObj.List("test"); e != nil {
		t.Error(e)
	} else {
		t.Log(v)
	}
}

func TestConsulKV_WatchPrefix(t *testing.T) {
	TestConsulKV_SetValue(t)
	consulTestKvObj.WatchPrefix("test", func(_ uint64, pairs api.KVPairs) {
		for _, v := range pairs {
			fmt.Println(v.Key, string(v.Value))
		}
	})

	time.Sleep(time.Second*1)

	if e := consulTestKvObj.SetValue("testXXX", "yyy"); e != nil {
		t.Error(e)
	}

	time.Sleep(time.Second*1)
}

func TestInit(t *testing.T) {
	if e := Init(DC, CLUSTER); e != nil {
		t.Error(e)
	}
}

func TestSetValue(t *testing.T) {
	TestInit(t)
	if e := SetValue("testXXX", "xxx"); e != nil {
		t.Error(e)
	}
}

func TestGetValue(t *testing.T) {
	TestSetValue(t)
	if v, e := GetValue("testXXX"); e != nil {
		t.Error(e)
	} else {
		fmt.Println(v)
	}
}

func TestDelValue(t *testing.T) {
	TestSetValue(t)
	if e := DelValue("testXXX"); e != nil {
		t.Error(e)
	}
}

func TestList(t *testing.T) {
	TestSetValue(t)
	if v, e := List("test"); e != nil {
		t.Error(e)
	} else {
		fmt.Println(v)
	}
}

func TestWatchPrefix(t *testing.T) {
	TestSetValue(t)
	WatchPrefix("test", func(_ uint64, pairs api.KVPairs) {
		for _, v := range pairs {
			fmt.Println(v.Key, string(v.Value))
		}
	})

	time.Sleep(time.Second*1)

	if e := SetValue("testXXX", "yyy"); e != nil {
		t.Error(e)
	}

	time.Sleep(time.Second*1)
}
