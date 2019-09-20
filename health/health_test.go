package health

import (
	"fmt"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	Init("test", "1.2.3.4", "/")
}

func TestSetCallback(t *testing.T) {
	TestInit(t)
	SetCallback(func(msg *ChunkMsg) {
		fmt.Println(string(msg.ToJson()))
	})

	time.Sleep(time.Second*16)
}