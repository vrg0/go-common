package logger

import (
	"io"
	"os"
	"testing"
)

func TestWatchWriter(t *testing.T) {
	ww := NewHookWriter(os.Stdout)
	ww.AddHookFunc(func(data []byte) bool {
		t.Log(string(data))
		return true
	})
	_, _ = io.WriteString(ww, "test AddHookFunc")
}
