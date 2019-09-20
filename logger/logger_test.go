package logger

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const (
	TestLogFile = "./logger_test.log"
)

func loggerClear() {
	sugar = nil
	defaultLogger = nil
}

func InitDev(t *testing.T) {
	loggerClear()
	if e := Init(false, ""); e != nil {
		t.Error(e)
	}
}

func InitPro(t *testing.T) {
	loggerClear()
	if e := Init(true, TestLogFile); e != nil {
		t.Error(e)
	}
}

func TestInit(t *testing.T) {
	InitDev(t)
	Debug("test Init")

	InitPro(t)
	Info("test Init")
	if e := GetSugar().Sync(); e != nil {
		t.Error(e)
	}
	if fileBytes, e := ioutil.ReadFile(TestLogFile); e != nil {
		t.Error(e)
	} else {
		t.Logf(string(fileBytes))
		if e := os.Remove(TestLogFile); e != nil {
			t.Error(e)
		}
	}
}

func TestDebug(t *testing.T) {
	InitDev(t)
	Debug("test Debug")
}

func TestDebugf(t *testing.T) {
	InitDev(t)
	Debugf("test Debug%s", "f")
}

func TestDebugw(t *testing.T) {
	InitDev(t)
	Debugw("test Debugw", "k1", "v1", "k2", "v2")
}

func TestInfo(t *testing.T) {
	InitDev(t)
	Info("test Info")
}

func TestInfof(t *testing.T) {
	InitDev(t)
	Infof("test Info%s", "f")
}

func TestInfow(t *testing.T) {
	InitDev(t)
	Infow("test Infow", "k1", "v1", "k2", "v2")
}

func TestPanic(t *testing.T) {
	InitDev(t)
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panic("test Panic")
}

func TestPanicf(t *testing.T) {
	InitDev(t)
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panicf("test Panic%s", "f")
}

func TestPanicw(t *testing.T) {
	InitDev(t)
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanicw(t *testing.T) {
	InitDev(t)
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	DPanicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanic(t *testing.T) {
	InitDev(t)
	DPanic("test DPanic")
}

func TestDPanicf(t *testing.T) {
	InitDev(t)
	DPanicf("test DPanic%s", "f")
}

func TestError(t *testing.T) {
	InitDev(t)
	Error("test Error")
}

func TestErrorf(t *testing.T) {
	InitDev(t)
	Errorf("test Error%s", "f")
}

func TestErrorw(t *testing.T) {
	InitDev(t)
	Errorw("test Errorw", "k1", "v1", "k2", "v2")
}

func TestFatal(t *testing.T) {
	InitDev(t)
	if os.Getenv("TEST_FATAL") == "1" {
		Fatal("test Fatal")
	}

	stdAll := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	cmd.Stdout = stdAll
	cmd.Stderr = stdAll
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		outBytes, _ := ioutil.ReadAll(stdAll)
		t.Logf(string(outBytes))
		return
	} else {
		t.Fatal("test Fatal err")
	}
}

func TestFatalf(t *testing.T) {
	InitDev(t)
	if os.Getenv("TEST_FATALF") == "1" {
		Fatalf("test Fatal%s", "f")
	}

	stdAll := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalf")
	cmd.Env = append(os.Environ(), "TEST_FATALF=1")
	cmd.Stdout = stdAll
	cmd.Stderr = stdAll
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		outBytes, _ := ioutil.ReadAll(stdAll)
		t.Logf(string(outBytes))
		return
	} else {
		t.Fatal("test Fatalf err")
	}
}

func TestFatalw(t *testing.T) {
	InitDev(t)
	if os.Getenv("TEST_FATALW") == "1" {
		Fatalw("test Fatalw", "k1", "v1", "k2", "v2")
	}

	stdAll := new(bytes.Buffer)
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalw")
	cmd.Env = append(os.Environ(), "TEST_FATALW=1")
	cmd.Stdout = stdAll
	cmd.Stderr = stdAll
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		outBytes, _ := ioutil.ReadAll(stdAll)
		t.Logf(string(outBytes))
		return
	} else {
		t.Fatal("test Fatalw err")
	}
	InitDev(t)
}

func TestWarn(t *testing.T) {
	InitDev(t)
	Warn("test Warn")
}

func TestWarnf(t *testing.T) {
	InitDev(t)
	Warnf("test Warn%s", "f")
}

func TestWarnw(t *testing.T) {
	InitDev(t)
	Warnw("test Warnw", "k1", "v1", "k2", "v2")
}

func TestGetDefaultLogger(t *testing.T) {
	InitDev(t)
	logger := GetDefaultLogger()
	if logger == nil {
		t.Error("defaultLog equals nil")
	} else {
		logger.Print("test GetDefaultLogger")
	}
}

func TestGetSugar(t *testing.T) {
	InitDev(t)
	logger := GetSugar()
	if logger == nil {
		t.Error("sugar equals nil")
	} else {
		logger.Info("test GetSugar")
	}
}
