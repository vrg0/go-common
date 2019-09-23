package args

import (
	"os"
	"strings"
	"testing"
)

var argNames = []string{"-env", "-idc", "-log_dir", "-eos_namespace", "-eos_host"}

//清空配置
func clearArgs() {
	newArgs := make([]string, 0)
	for _, v := range os.Args {
		flag := false
		for _, subStr := range argNames {
			if strings.Contains(v, subStr) {
				flag = true
				break
			}
		}
		if !flag {
			newArgs = append(newArgs, v)
		}
	}
	os.Args = newArgs
	isInitialized = false
}

//正常生产环境初始化
func TestInit(t *testing.T) {
	clearArgs()
	os.Args = append(os.Args, "-env=pro")
	os.Args = append(os.Args, "-idc=dc1")
	os.Args = append(os.Args, "-log_dir=/xxx/log")
	if e := Init("test"); e != nil {
		t.Error(e)
	}
}

//默认值测试环境初始化
func TestInit2(t *testing.T) {
	clearArgs()
	os.Args = append(os.Args, "-eos_namespace=xxx")
	os.Args = append(os.Args, "-eos_host=http://xxx.xx")
	if e := Init("test2"); e != nil {
		t.Error(e)
	}
}

//异常初始化，忽略-eos_namespace
func TestInit3(t *testing.T) {
	clearArgs()
	os.Args = append(os.Args, "-eos_host=http://xxx.xx")
	//	os.Args = append(os.Args, "-env=pro") //pro默认无eos_namespace参数
	os.Args = append(os.Args, "-idc=dc1")
	os.Args = append(os.Args, "-log_dir=/xxx/log")
	if e := Init("test3"); e != nil {
		if e.Error() != "eos_namespace can not be empty" {
			t.Error(e)
		}
	}
}

//异常初始化，忽略-eos_host
func TestInit4(t *testing.T) {
	clearArgs()
	os.Args = append(os.Args, "-eos_namespace=xxx")
	//	os.Args = append(os.Args, "-env=pro") //pro默认无eos_namespace参数
	os.Args = append(os.Args, "-idc=dc1")
	os.Args = append(os.Args, "-log_dir=/xxx/log")
	if e := Init("test4"); e != nil {
		if e.Error() != "eos_host can not be empty" {
			t.Error(e)
		}
	}
}

//初始化单例测试
func TestInit5(t *testing.T) {
	clearArgs()
	os.Args = append(os.Args, "-eos_namespace=xxx")
	os.Args = append(os.Args, "-eos_host=http://xxx.xx")
	if e := Init("test5"); e != nil {
		t.Error(e)
	}

	if e := Init("test5"); e != nil {
		if e.Error() != "the args module have been initialized" {
			t.Error(e)
		}
	}
}

func TestGetEnv(t *testing.T) {
	TestInit(t)
	if GetEnv() != "pro" {
		t.Error("GetEnv pro")
	}

	TestInit2(t)
	if GetEnv() != "dev" {
		t.Error("GetEnv dev")
	}
}

func TestGetEosHost(t *testing.T) {
	TestInit(t)
	if GetEosHost() != "" {
		t.Log("GetEosHost pro")
	}

	TestInit2(t)
	if GetEosHost() != "http://xxx.xx/" {
		t.Log("GetEosHost dev")
	}
}

func TestGetEosNamespace(t *testing.T) {
	TestInit(t)
	if GetEosNamespace() != "" {
		t.Log("GetEosNamespace pro")
	}

	TestInit2(t)
	if GetEosNamespace() != "xxx" {
		t.Log("GetEosNamespace dev")
	}
}

func TestGetIdc(t *testing.T) {
	TestInit(t)
	if GetIdc() != "dc1" {
		t.Log("GetIdc pro")
	}

	TestInit2(t)
	if GetIdc() != "k8s" {
		t.Log("GetIdc dev")
	}
}

func TestGetLogDir(t *testing.T) {
	TestInit(t)
	if GetLogDir() != "/xxx/log/" {
		t.Log("GetLogDir pro")
	}

	TestInit2(t)
	if GetLogDir() != "/var/log/" {
		t.Log("GetLogDir dev")
	}
}

func TestGetLogPath(t *testing.T) {
	TestInit(t)
	if GetLogPath() != "/xxx/log/test/test.log" {
		t.Log("GetLogPath pro")
	}

	TestInit2(t)
	if GetLogPath() != "/var/log/test2/test2.log" {
		t.Log("GetLogPath dev")
	}
}

func TestEnvIsPro(t *testing.T) {
	TestInit(t)
	if !EnvIsPro() {
		t.Log("env is pro")
	}

	TestInit2(t)
	if EnvIsPro() {
		t.Log("env is not pro")
	}
}

func TestGetAppName(t *testing.T) {
	TestInit(t)
	if GetAppName() != "test" {
		t.Log("get app name pro")
	}

	TestInit2(t)
	if GetAppName() != "test2" {
		t.Log("get app name dev")
	}
}
