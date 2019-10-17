package logger

import (
	"io"
	"os"
	"testing"
)

func TestWatchWriter(t *testing.T) {
	ww := NewHookWriter(os.Stdout)
	ww.AddHookFunc(func(data []byte) {
		t.Log(string(data))
	})
	_, _ = io.WriteString(ww, "test AddHookFunc")
}
