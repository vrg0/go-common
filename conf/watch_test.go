package conf

import (
	"fmt"
	"testing"
	"time"
)

func TestWatchNamespace(t *testing.T) {
	tConf.WatchNamespace("blacklist.url", func(oldCfgs map[string]string, newCfgs map[string]string) {
		fmt.Println("---------------------------")
		fmt.Println(oldCfgs)
		fmt.Println(newCfgs)
	})
	time.Sleep(time.Second*1)
}

func TestWatchKey(t *testing.T) {
	tConf.Watch("blacklist.url", "kkk", func(oldCfg string, newCfg string) {
		fmt.Println("---------------------------")
		fmt.Println(oldCfg)
		fmt.Println(newCfg)
	})
	time.Sleep(time.Second*1)
}
