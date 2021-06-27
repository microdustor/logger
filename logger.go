package logger

import (
	"os"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"github.com/ivpusic/grpool"
	"strings"
	"io"
)

func GetDefaultLogger(maxsize int, flag int, numWorkers int, jobQueueLen int, depth int) *Logger {
	logger := NewLogger2("/tmp/message", maxsize, flag, numWorkers, jobQueueLen, depth)
	logger.SetLevel(loadLogLevel("info"))
	return logger
}

func NewLogger(flag int, numWorkers int, jobQueueLen int, depth int) *Logger {
	Logger := NewLogger3(os.Stdout, flag, numWorkers, jobQueueLen, depth)
	return Logger
}

func NewLogger2(lfn string, maxsize int, flag int, numWorkers int, jobQueueLen int, depth int) *Logger {
	jack := &lumberjack.Logger{
		Filename: lfn,
		MaxSize:  maxsize, // megabytes
		MaxAge:   1,
		Compress: true,
	}

	logger := NewLogger3(jack, flag, numWorkers, jobQueueLen, depth)

	return logger
}

func NewLogger3(w io.Writer, flag int, numWorkers int, jobQueueLen int, depth int) *Logger {
	logger := new(Logger)
	logger.depth = depth
	if logger.depth <= 0 {
		logger.depth = 2
	}
	logger.std = log.New(os.Stdout, "[Std]", flag)
	logger.err = log.New(w, "[Err] ", flag)
	logger.warn = log.New(w, "[War] ", flag)
	logger.info = log.New(w, "[Inf] ", flag)
	logger.debug = log.New(w, "[Deb] ", flag)
	logger.SetLevel(LevelInformational)

	logger.p = grpool.NewPool(numWorkers, jobQueueLen)

	return logger
}

func (ll *Logger) SetLevel(l int) {
	ll.level = l
}

// 统一设置日志前缀
func (ll *Logger) SetPrefix(prefix string) {
	ll.err.SetPrefix(prefix)
	ll.warn.SetPrefix(prefix)
	ll.info.SetPrefix(prefix)
	ll.debug.SetPrefix(prefix)
}

func (ll *Logger) Fatal(format string, v ...interface{}) {
	ll.p.JobQueue <- func() {
		ll.err.Output(ll.depth, fmt.Sprintf(format, v...))
		ll.std.Output(ll.depth, fmt.Sprintf(format, v...))
	}
	log.Fatalf(format, v...)
}

func (ll *Logger) Error(format string, v ...interface{}) {
	if LevelError > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.err.Output(ll.depth, fmt.Sprintf(format, v...))
		ll.std.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) Warn(format string, v ...interface{}) {
	if LevelWarning > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.warn.Output(ll.depth, fmt.Sprintf(format, v...))
		ll.std.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) Info(format string, v ...interface{}) {
	if LevelInformational > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.info.Output(ll.depth, fmt.Sprintf(format, v...))
		ll.std.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) Debug(format string, v ...interface{}) {
	if LevelDebug > ll.level {
		return
	}
	ll.p.JobQueue <- func() {
		ll.debug.Output(ll.depth, fmt.Sprintf(format, v...))
		ll.std.Output(ll.depth, fmt.Sprintf(format, v...))
	}
}

func (ll *Logger) SetJack(lfn string, maxsize int) {
	jack := &lumberjack.Logger{
		Filename: lfn,
		MaxSize:  maxsize, // megabytes
	}

	ll.err.SetOutput(jack)
	ll.warn.SetOutput(jack)
	ll.info.SetOutput(jack)
	ll.debug.SetOutput(jack)
}

func (ll *Logger) SetFlag(flag int) {
	ll.err.SetFlags(flag)
	ll.warn.SetFlags(flag)
	ll.debug.SetFlags(flag)
}

func (ll *Logger) Stats() (int, int) {
	return cap(ll.p.JobQueue), len(ll.p.JobQueue)
}

// ================= DefaultLogger ======================

func loadLogLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return 4
	case "info":
		return 3
	case "warn":
		return 2
	case "error":
		return 1
	default:
		return 3
	}
}

func Fatal(format string, v ...interface{}) {
	DefaultLogger.Fatal(format, v...)
}

func SetJack(lfn string, maxsize int) {
	DefaultLogger.SetJack(lfn, maxsize)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Error(format, v...)
}

func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warn(format, v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Info(format, v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debug(format, v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(GenerateFmtStr(len(v)), v...)
}

func Warn(v ...interface{}) {
	DefaultLogger.Warn(GenerateFmtStr(len(v)), v...)
}

func Info(v ...interface{}) {
	DefaultLogger.Info(GenerateFmtStr(len(v)), v...)
}

func Debug(v ...interface{}) {
	DefaultLogger.Debug(GenerateFmtStr(len(v)), v...)
}

func LogLevel(logLevel string) string {
	if len(logLevel) == 0 {
		logLevel = "info"
	}
	updateLevel(logLevel)
	Warn("Set Log Level as", logLevel)
	return logLevel
}

func updateLevel(logLevel string) {
	switch strings.ToLower(logLevel) {
	case "debug":
		DefaultLogger.SetLevel(LevelDebug)
	case "info":
		DefaultLogger.SetLevel(LevelInformational)
	case "warn":
		DefaultLogger.SetLevel(LevelWarning)
	case "error":
		DefaultLogger.SetLevel(LevelError)
	default:
		DefaultLogger.SetLevel(LevelInformational)
	}
}

func GenerateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}

