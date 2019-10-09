package logger

import (
	"github.com/natefinch/lumberjack"
	"github.com/vrg0/go-common/args"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"sync/atomic"
	"unsafe"
)

type Logger struct {
	sugar  *zap.SugaredLogger
	logger *zap.Logger
	writer *watchWriter
}

var (
	defaultLogger *Logger = nil
)

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

	defaultLogger = New(logPath, level)
}

func ResetDefaultLogger(logPath string, level zapcore.Level) {
	if logPath == "" {
		logPath = "/dev/stdout'"
	}

	newLogger := New(logPath, level)
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&defaultLogger)), unsafe.Pointer(newLogger))
}

// 新建Logger对象，成功返回对象指针，失败返回nil
func New(logPath string, level zapcore.Level) *Logger {
	//参数过滤
	if logPath == "" {
		logPath = "/dev/stdout"
	}

	rtn := Logger{}

	var writer io.Writer
	switch logPath {
	case "/dev/stdout":
		writer = os.Stdout
	case "/dev/stderr":
		writer = os.Stderr
	default:
		writer = &lumberjack.Logger{
			Filename:   logPath, //日志路径
			MaxSize:    1024,    //日志大小，单位MB
			MaxBackups: 30,      //日志文件最多保存备份
			MaxAge:     7,       //日志文件最多保存多少天
			LocalTime:  true,    //打印本地时间
			Compress:   false,   //日志备份不进行压缩，压缩会导致占用过多cpu
		}
	}
	rtn.writer = newWatchWriter(writer)

	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(config)
	writeSyncer := zapcore.AddSync(rtn.writer)
	logger := zap.New(zapcore.NewCore(encoder, writeSyncer, level))
	rtn.logger = logger
	rtn.sugar = logger.Sugar()

	return &rtn
}

func (l *Logger) GetStandardLogger() *log.Logger {
	return zap.NewStdLog(l.logger)
}

func (l *Logger) SetWatchFunc(watchFunc WatchFunc) {
	l.writer.AddWatchFunc(watchFunc)
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugar.Debug(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugar.Info(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugar.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugar.Error(args...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.sugar.DPanic(args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sugar.Panic(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugar.Fatal(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.sugar.DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugar.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugar.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func DPanic(args ...interface{}) {
	defaultLogger.DPanic(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	defaultLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	defaultLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultLogger.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	defaultLogger.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultLogger.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultLogger.Fatalf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultLogger.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Errorw(msg, keysAndValues...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.DPanicw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultLogger.Fatalw(msg, keysAndValues...)
}

func SetWatchFunc(watchFunc WatchFunc) {
	defaultLogger.writer.AddWatchFunc(watchFunc)
}

func GetStandardLogger() *log.Logger {
	return defaultLogger.GetStandardLogger()
}
