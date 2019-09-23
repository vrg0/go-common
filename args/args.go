package args

/**
 * 使用args模块必须进行初始化
 */

import (
	"errors"
	"flag"
	"os"
	"strings"
)

var (
	env           string
	idc           string
	logDir        string
	logPath       string
	eosHost       string
	eosNamespace  string
	appName       string
	isInitialized = false
	argNameList = []string{"-env", "-idc", "-log_dir", "-eos_host", "-eos_namespace"}
)

func Init(applicationName string) error {
	//懒汉单例
	if isInitialized {
		return errors.New("the args module have been initialized")
	} else {
		//获取args
		args := make([]string, 0)
		for _, arg := range os.Args {
			for _, subStr := range argNameList {
				if strings.Contains(arg, subStr) {
					args = append(args, arg)
					break
				}
			}
		}

		flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		flagSet.StringVar(&env, "env", "dev", "Env: dev/pro")
		flagSet.StringVar(&idc, "idc", "k8s", "Idc: k8s/3C")
		flagSet.StringVar(&logDir, "log_dir", "/var/log/", "LogPath: /var/log/")
		flagSet.StringVar(&eosHost, "eos_host", "", "EosHost: No default value")
		flagSet.StringVar(&eosNamespace, "eos_namespace", "", "EosNamespace: No default value")
		if e := flagSet.Parse(args); e != nil {
			return e
		}
	}

	//过滤参数
	applicationName = strings.TrimSpace(applicationName) //去除空白
	if applicationName == "" {
		return errors.New("applicationName can not be empty")
	}
	if env == "" {
		return errors.New("env can not be empty")
	}
	if idc == "" {
		return errors.New("idc can not be empty")
	}
	if logDir == "" {
		return errors.New("logDir can not be empty")
	} else if logDir[len(logDir)-1] != '/' {
		logDir += "/"
	}
	if !EnvIsPro() {
		if eosNamespace == "" {
			return errors.New("eos_namespace can not be empty")
		}
		if eosHost == "" {
			return errors.New("eos_host can not be empty")
		} else if eosHost[len(eosHost)-1] != '/' {
			eosHost += "/"
		}
	}
	logPath = logDir + applicationName + "/" + applicationName + ".log"
	appName = applicationName

	//标记
	isInitialized = true

	return nil
}

func GetLogPath() string {
	return logPath
}

func GetLogDir() string {
	return logDir
}

func GetEosNamespace() string {
	return eosNamespace
}

func GetEosHost() string {
	return eosHost
}

func GetEnv() string {
	return env
}

func GetIdc() string {
	return idc
}

func EnvIsPro() bool {
	return env == "pro"
}

func GetAppName() string {
	return appName
}
