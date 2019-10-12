package conf

import "testing"

func TestGet(t *testing.T) {
	v, ok := Get("blacklist.url", "kkk")
	if !ok {

	} else {
		t.Log(v)
	}
}

/*
func TestGetNamespace(t *testing.T) {

}

func TestGetOrDefault(t *testing.T) {

}

func TestRefreshKvMap(t *testing.T) {

}
 */

/*
import (
	"log"
	"os"
	"testing"
)

//1、配置中心存在时，会从配置中心中取出数据。
//2、配置中心断开时，会从缓存文件中取出数据。
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

	if v, ok := kvMap["REDIS_SERVICE_HOST"]; !ok || v != "1.1.1.1" {
		t.Error("refreshKvMap fatal")
	}

	if v, ok := kvMap["REDIS_SERVICE_PORT_6379"]; !ok || v != "6379" {
		t.Error("refreshKvMap fatal")
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
*/
