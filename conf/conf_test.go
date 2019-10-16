package conf

import (
	"github.com/vrg0/go-common/args"
	"os"
	"testing"
)

var tConf *Conf

func init() {
	configServer, ok := args.Get("config_server")
	if !ok {
		return
	}
	appName, ok := args.Get("app_name")
	if !ok {
		return
	}
	idc, ok := args.Get("idc")
	if !ok {
		return
	}
	cacheFilePath := args.GetOrDefault("cache_file_path", os.Args[0]+".cache_file")

	if newConf, err := New(configServer, appName, idc, cacheFilePath, nil); err != nil {
		return
	} else {
		tConf = newConf
	}
}

func TestGet(t *testing.T) {
	v, ok := tConf.Get("blacklist.url", "kkk")
	if !ok {
		t.Error("not found conf")
	} else {
		t.Log(v)
	}
}

func TestGet2(t *testing.T) {
	v, ok := tConf.Get("blacklist.url", "none")
	if ok {
		t.Error(v)
	} else {
		t.Log(v)
	}
}

func TestGetNamespace(t *testing.T) {
	v := tConf.GetNamespace("blacklist.url")
	t.Log(v)
}

func TestGetOrDefault(t *testing.T) {
	v := tConf.GetOrDefault("blacklist.url", "none", "123")
	if v != "123" {
		t.Log(v)
	}

	v2 := tConf.GetOrDefault("blacklist.url", "kkk", "345")
	if v2 == "345" {
		t.Log(v2)
	}
}

func TestRefreshKvMap(t *testing.T) {
	tConf.RefreshKvMap(map[string]string{"123":"aaa"})
	v := tConf.GetOrDefault("blacklist.url", "kv", "xxx")
	if v != "/aaa" {
		t.Error(v)
	} else {
		t.Log(v)
	}
}