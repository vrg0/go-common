package conf

import (
	"log"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	appId := os.Getenv("TEST_APP_ID")
	domainName := os.Getenv("TEST_DOMAIN_NAME")
	env := os.Getenv("TEST_ENV")
	idc := os.Getenv("TEST_IDC")

	logger := log.New(os.Stdout, "", 0)
	if e := Init(appId, domainName, env, idc, "", logger); e != nil {
		t.Error(e)
	}
}

func TestRefreshKvMap(t *testing.T) {
	kvMap := make(map[string]string)
	kvMap["REDIS_SERVICE_HOST"] = "1.1.1.1"
	kvMap["REDIS_SERVICE_PORT_6379"] = "6379"
	if e := RefreshKvMap(kvMap); e != nil {
		t.Error(e)
	}
}

func TestGet(t *testing.T) {
	TestInit(t)
	TestRefreshKvMap(t)
	t.Log("1:", Get("aaa", "bbb"))
	t.Log("2:", Get("blacklist.url", "aaaa"))
	t.Log("3:", Get("application", "redis.master.addr"))
}

func TestGetWithDefault(t *testing.T) {
	TestInit(t)
	TestRefreshKvMap(t)
	t.Log("1:", GetWithDefault("bbb"))
	t.Log("2:", GetWithDefault("aaaa"))
	t.Log("3:", GetWithDefault("redis.master.addr"))
}

func TestGetNamespace(t *testing.T) {
	TestInit(t)
	TestRefreshKvMap(t)
	kv := GetNamespace("application")
	for k, v := range kv {
		t.Log(k, v)
	}
}
