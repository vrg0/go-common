package args

import (
	"os"
	"sync"
	"testing"
)

var testInitOnce = sync.Once{}

func TestInit(t *testing.T) {
	testInitOnce.Do(func() {
		os.Args = append(os.Args, "-eos_namespace=xxx")
		os.Args = append(os.Args, "-eos_host=http://xxx.xx")
	})
	if e := Init("test"); e != nil {
		if e.Error() != "the args module have been initialized" {
			t.Error(e)
		}
	}
}

func TestGetEnv(t *testing.T) {
	TestInit(t)
	t.Log(GetEnv())
}

func TestGetEosHost(t *testing.T) {
	TestInit(t)
	t.Log(GetEosHost())
}

func TestGetEosNamespace(t *testing.T) {
	TestInit(t)
	t.Log(GetEosNamespace())
}

func TestGetIdc(t *testing.T) {
	TestInit(t)
	t.Log(GetIdc())
}

func TestGetLogDir(t *testing.T) {
	TestInit(t)
	t.Log(GetLogDir())
}

func TestGetLogPath(t *testing.T) {
	TestInit(t)
	t.Log(GetLogPath())
}

func TestEnvIsPro(t *testing.T) {
	TestInit(t)
	t.Log(EnvIsPro())
}

func TestGetAppName(t *testing.T) {
	TestInit(t)
	t.Log(GetAppName())
}