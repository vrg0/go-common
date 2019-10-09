package logger

import (
	"io"
	"sync"
)

type WatchFunc func(data []byte)

type watchWriter struct {
	writer io.Writer
	watchFuncList []WatchFunc
	watchFuncListLock *sync.RWMutex
}

func newWatchWriter(w io.Writer) *watchWriter {
	return &watchWriter{
		writer:w,
		watchFuncList:make([]WatchFunc, 0),
		watchFuncListLock:new(sync.RWMutex),
	}
}

func (ww *watchWriter) Write(p []byte) (n int, err error) {
	if len(ww.watchFuncList) != 0 {
		ww.watchFuncListLock.RLock()
		for _, handler := range ww.watchFuncList {
			handler(p)
		}
		ww.watchFuncListLock.RUnlock()
	}

	return ww.writer.Write(p)
}

func (ww *watchWriter) AddWatchFunc(watchFunc WatchFunc) {
	ww.watchFuncListLock.Lock()
	ww.watchFuncList = append(ww.watchFuncList, watchFunc)
	ww.watchFuncListLock.Unlock()
}
