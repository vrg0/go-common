package logger

import (
	"io"
	"sync"
)

type HookFunc func(data []byte)

type hookWriter struct {
	writer           io.Writer
	hookFuncList     []HookFunc
	hookFuncListLock *sync.RWMutex
}

func newHookWriter(w io.Writer) *hookWriter {
	return &hookWriter{
		writer:           w,
		hookFuncList:     make([]HookFunc, 0),
		hookFuncListLock: new(sync.RWMutex),
	}
}

func (hw *hookWriter) Write(p []byte) (n int, err error) {
	if len(hw.hookFuncList) != 0 {
		hw.hookFuncListLock.RLock()
		for _, handler := range hw.hookFuncList {
			handler(p)
		}
		hw.hookFuncListLock.RUnlock()
	}

	return hw.writer.Write(p)
}

func (hw *hookWriter) AddHookFunc(hookFunc HookFunc) {
	hw.hookFuncListLock.Lock()
	hw.hookFuncList = append(hw.hookFuncList, hookFunc)
	hw.hookFuncListLock.Unlock()
}
