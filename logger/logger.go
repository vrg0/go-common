package logger

/**
 * 使用log模块必须进行初始化
 */

import (
	"github.com/natefinch/lumberjack"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var sugar *zap.SugaredLogger = nil
var defaultLogger *log.Logger = nil

//当envIsPro为true时，logPath不能为空，日志会打印到指定文件
//等envIsPro为false时，logPath参数无效，日志会打印到标准输出
func Init(envIsPro bool, logPath string) error {
	if defaultLogger != nil || sugar != nil {
		return errors.New("the logger module have been initialized")
	}

	if envIsPro && logPath == "" {
		return errors.New("when envIsPro equals true, logPath can not be empty")
	}

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(config)
	var logger *zap.Logger

	if envIsPro { //生产环境，打印到文件，级别Info
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logPath, //日志路径
			MaxSize:    1024,    //日志大小，单位MB
			MaxBackups: 30,      //日志文件最多保存备份
			MaxAge:     7,       //日志文件最多保存多少天
			LocalTime:  true,    //打印本地时间
			Compress:   false,   //日志备份不进行压缩，压缩会导致占用过多cpu
		})
		logger = zap.New(zapcore.NewCore(encoder, writeSyncer, zap.InfoLevel))
	} else { //测试环境，打印到标准输出，级别Debug
		writeSyncer := zapcore.AddSync(os.Stdout)
		logger = zap.New(zapcore.NewCore(encoder, writeSyncer, zap.DebugLevel))
	}

	defer func() {
		_ = logger.Sync()
	}()

	sugar = logger.Sugar()
	defaultLogger = zap.NewStdLog(logger)

	return nil
}

func GetSugar() *zap.SugaredLogger {
	return sugar
}

func GetDefaultLogger() *log.Logger {
	return defaultLogger
}

func Debug(args ...interface{}) {
	sugar.Debug(args...)
}

func Info(args ...interface{}) {
	sugar.Info(args...)
}

func Warn(args ...interface{}) {
	sugar.Warn(args...)
}

func Error(args ...interface{}) {
	sugar.Error(args...)
}

func DPanic(args ...interface{}) {
	sugar.DPanic(args...)
}

func Panic(args ...interface{}) {
	sugar.Panic(args...)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	sugar.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	sugar.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	sugar.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	sugar.Fatalw(msg, keysAndValues...)
}
