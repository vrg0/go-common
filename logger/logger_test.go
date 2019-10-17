package logger

import (
	"bytes"
	"github.com/vrg0/go-common/args"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var tLog *Logger = nil

func init() {
	env := args.GetOrDefault("env", "dev")
	var logPath string
	var level zapcore.Level
	if env == "dev" {
		logPath = args.GetOrDefault("log_path", "/dev/stdout")
		level = zapcore.DebugLevel
	} else {
		logPath = args.GetOrDefault("log_path", os.Args[0]+".log")
		level = zapcore.InfoLevel
	}

	tLog = New(logPath, level)
}

func TestDebug(t *testing.T) {
	tLog.Debug("test Debug")
}

func TestDebugf(t *testing.T) {
	tLog.Debugf("test Debug%s", "f")
}

func TestDebugw(t *testing.T) {
	tLog.Debugw("test Debugw", "k1", "v1", "k2", "v2")
}

func TestInfo(t *testing.T) {
	tLog.Info("test Info")
}

func TestInfof(t *testing.T) {
	tLog.Infof("test Info%s", "f")
}

func TestInfow(t *testing.T) {
	tLog.Infow("test Infow", "k1", "v1", "k2", "v2")
}

func TestPanic(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	tLog.Panic("test Panic")
}

func TestPanicf(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	tLog.Panicf("test Panic%s", "f")
}

func TestPanicw(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	tLog.Panicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanicw(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Log(e)
		}
	}()
	tLog.DPanicw("test Panicw", "k1", "v1", "k2", "v2")
}

func TestDPanic(t *testing.T) {
	tLog.DPanic("test DPanic")
}

func TestDPanicf(t *testing.T) {
	tLog.DPanicf("test DPanic%s", "f")
}

func TestError(t *testing.T) {
	tLog.Error("test Error")
}

func TestErrorf(t *testing.T) {
	tLog.Errorf("test Error%s", "f")
}

func TestErrorw(t *testing.T) {
	tLog.Errorw("test Errorw", "k1", "v1", "k2", "v2")
}

func TestFatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		tLog.Fatal("test Fatal")
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
		tLog.Fatalf("test Fatal%s", "f")
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
		tLog.Fatalw("test Fatalw", "k1", "v1", "k2", "v2")
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
	tLog.Warn("test Warn")
}

func TestWarnf(t *testing.T) {
	tLog.Warnf("test Warn%s", "f")
}

func TestWarnw(t *testing.T) {
	tLog.Warnw("test Warnw", "k1", "v1", "k2", "v2")
}

/*
func TestResetDefaultLogger(t *testing.T) {
	ResetDefaultLogger("/dev/stderr", zapcore.InfoLevel)
	Info("reset default logger")
}
*/

func TestSetHookFunc(t *testing.T) {
	tLog.SetHookFunc(func(data []byte) bool {
		t.Log("data hook", string(data))
		return true
	})

	tLog.Info("test hook")
}

func TestGetStandardLogger(t *testing.T) {
	sl := tLog.GetStandardLogger()
	sl.SetPrefix("_TestStandardLogger_")
	tLog.SetHookFunc(func(data []byte) bool {
		if strings.Contains(string(data), "_TestStandardLogger_") {
			t.Log("hook _TestStandardLogger_")
		}
		return true
	})
	sl.Print("test get standard logger")
}
