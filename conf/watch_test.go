package conf

import (
	"fmt"
	"testing"
	"time"
)

func TestStartWatch(t *testing.T) {
	TestInit(t)
	StartWatch()
}

func TestWatchNamespace(t *testing.T) {
	TestStartWatch(t)
	WatchNamespace("blacklist.url", func(oldCfgs map[string]string, newCfgs map[string]string) {
		fmt.Println("---------------------------")
		fmt.Println(oldCfgs)
		fmt.Println(newCfgs)
	})
	time.Sleep(time.Second*1)
}

func TestWatchKey(t *testing.T) {
	TestStartWatch(t)
	Watch("blacklist.url", "kkk", func(oldCfg string, newCfg string) {
		fmt.Println("---------------------------")
		fmt.Println(oldCfg)
		fmt.Println(newCfg)
	})
	time.Sleep(time.Second*1)
}