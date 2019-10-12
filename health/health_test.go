package health

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	if e := Start("test", "/"); e != nil {
		t.Error(e)
	}
}

func TestSetCallback(t *testing.T) {
	TestInit(t)
	SetCallback(func(msg *ChunkMsg) {
		fmt.Println(string(msg.ToJson()))
	})

	time.Sleep(time.Second * 16)
}
