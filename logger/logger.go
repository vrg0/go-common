package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

type Logger struct {
	sugar  *zap.SugaredLogger
	logger *zap.Logger
	writer *hookWriter
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
	rtn.writer = NewHookWriter(writer)

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

func (l *Logger) GetSugaredLogger() *zap.SugaredLogger {
	return l.logger.Sugar()
}

func (l *Logger) SetHookFunc(hookFunc HookFunc) {
	l.writer.AddHookFunc(hookFunc)
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

