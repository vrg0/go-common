package logger

import (
	"bytes"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestDebug(t *testing.T) {
	Debug("test Debug")
}

func TestDebugf(t *testing.T) {
	Debugf("test Debug%s", "f")
}

func TestDebugw(t *testing.T) {
	Debugw("test Debugw", "k1", "v1", "k2", "v2")
}

func TestInfo(t *testing.T) {
	Info("test Info")
}

func TestInfof(t *testing.T) {
	Infof("test Info%s", "f")
}

func TestInfow(t *testing.T) {
	Infow("test Infow", "k1", "v1", "k2", "v2")
}

func TestPanic(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panic("test Panic")
}

func TestPanicf(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panicf("test Panic%s", "f")
}

func TestPanicw(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	Panicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanicw(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	DPanicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanic(t *testing.T) {
	DPanic("test DPanic")
}

func TestDPanicf(t *testing.T) {
	DPanicf("test DPanic%s", "f")
}

func TestError(t *testing.T) {
	Error("test Error")
}

func TestErrorf(t *testing.T) {
	Errorf("test Error%s", "f")
}

func TestErrorw(t *testing.T) {
	Errorw("test Errorw", "k1", "v1", "k2", "v2")
}

func TestFatal(t *testing.T) {
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
}

func TestWarn(t *testing.T) {
	Warn("test Warn")
}

func TestWarnf(t *testing.T) {
	Warnf("test Warn%s", "f")
}

func TestWarnw(t *testing.T) {
	Warnw("test Warnw", "k1", "v1", "k2", "v2")
}

func TestResetDefaultLogger(t *testing.T) {
	ResetDefaultLogger("/dev/stderr", zapcore.InfoLevel)
	Info("reset default logger")
}

func TestSetHookFunc(t *testing.T) {
	SetHookFunc(func(data []byte) {
		t.Log("data hook", string(data))
	})

	Info("test hook")
}

func TestGetStandardLogger(t *testing.T) {
	sl := GetStandardLogger()
	sl.SetPrefix("_TestStandardLogger_")
	SetHookFunc(func(data []byte) {
		if strings.Contains(string(data), "_TestStandardLogger_") {
			t.Log("hook _TestStandardLogger_")
		}
	})
	sl.Print("test get standard logger")
}
