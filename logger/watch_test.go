package logger

import (
	"io"
	"os"
	"testing"
)

func TestWatchWriter(t *testing.T) {
	ww := newWatchWriter(os.Stdout)
	ww.AddWatchFunc(func(data []byte) {
		t.Log(string(data))
	})
	_, _ = io.WriteString(ww, "test AddWatchFunc")
}
