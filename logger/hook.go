package logger

import (
	"io"
)

//当返回false时，则不会调用hw.writer.Writer()
type HookFunc func(data []byte) bool

type hookWriter struct {
	writer           io.Writer
	hookFuncList     []HookFunc
}

func NewHookWriter(w io.Writer) *hookWriter {
	return &hookWriter{
		writer:           w,
		hookFuncList:     make([]HookFunc, 0),
	}
}

func (hw *hookWriter) Write(p []byte) (n int, err error) {
	if len(hw.hookFuncList) != 0 {
		for _, handler := range hw.hookFuncList {
			if !handler(p) {
				return
			}
		}
	}

	return hw.writer.Write(p)
}

func (hw *hookWriter) AddHookFunc(hookFunc HookFunc) {
	newFuncList := make([]HookFunc, 0)
	for _, handler := range hw.hookFuncList {
		newFuncList = append(newFuncList, handler)
	}
	newFuncList = append(newFuncList, hookFunc)
	hw.hookFuncList = newFuncList
}
