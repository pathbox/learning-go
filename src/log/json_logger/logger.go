package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

type LogWriter struct {
	Writers []io.Writer
}

func (lw *LogWriter) Write(content []byte) (n int, err error) {
	n = 0
	err = nil
	for _, writer := range lw.Writers {
		n, err = writer.Write(content)
		if err != nil {
			return
		}
	}
	return
}

type Fields logrus.Fields

type PreLogger interface {
	PreLog() *Fields
}

type Entry struct {
	LogrusEntry *logrus.Entry
	Pre         PreLogger
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(HereFormatter())
}

func preLog() *Entry {
	pc, fileName, lineNum, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	entry := logrus.WithFields(logrus.Fields{
		"_loc":  fmt.Sprintf("%s:%d", path.Base(fileName), lineNum),
		"_func": funcName,
	})

	return &Entry{LogrusEntry: entry, Pre: nil}
}

func preLogFatal() *Entry {
	pc, fileName, lineNum, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	entry := logrus.WithFields(logrus.Fields{
		"_loc":   fmt.Sprintf("%s:%d", path.Base(fileName), lineNum),
		"_func":  funcName,
		"_trace": string(debug.Stack()),
	})

	return &Entry{LogrusEntry: entry, Pre: nil}
}

func preLogWithFields(fields *Fields) *Entry {
	newEntry := logrus.WithFields(logrus.Fields(*fields))
	return &Entry{LogrusEntry: newEntry, Pre: nil}
}

func SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
}

func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

func AddHook(hook logrus.Hook) {
	logrus.AddHook(hook)
}

func Debug(args ...interface{}) {
	preLog().LogrusEntry.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	preLog().LogrusEntry.Debugf(format, args...)
}

func Info(args ...interface{}) {
	preLog().LogrusEntry.Info(args...)
}

func Infof(format string, args ...interface{}) {
	preLog().LogrusEntry.Infof(format, args...)
}

func Warn(args ...interface{}) {
	preLog().LogrusEntry.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	preLog().LogrusEntry.Warnf(format, args...)
}

func Error(args ...interface{}) {
	preLog().LogrusEntry.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	preLog().LogrusEntry.Errorf(format, args...)
}
func Fatal(args ...interface{}) {
	entry := preLogFatal().LogrusEntry
	entry.Time = time.Now()
	entry.Level = logrus.FatalLevel
	entry.Message = fmt.Sprint(args...)
	entry.Logger.Hooks.Fire(logrus.FatalLevel, entry)
	serialized, _ := entry.Logger.Formatter.Format(entry)
	entry.Logger.Out.Write(serialized)

}

func Fatalf(format string, args ...interface{}) {
	entry := preLogFatal().LogrusEntry
	entry.Time = time.Now()
	entry.Level = logrus.FatalLevel
	entry.Message = fmt.Sprintf(format, args...)
	entry.Logger.Hooks.Fire(logrus.FatalLevel, entry)
	serialized, _ := entry.Logger.Formatter.Format(entry)
	entry.Logger.Out.Write(serialized)
}

func Panic(args ...interface{}) {
	preLogFatal().LogrusEntry.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	preLogFatal().LogrusEntry.Panicf(format, args...)
}

// 用此函数返回的是一个新对象指针
func WithFields(fields *Fields) *Entry {
	return preLogWithFields(fields)
}

func (e Entry) WithFields(fields *Fields) *Entry {
	e.LogrusEntry = e.LogrusEntry.WithFields(logrus.Fields(*fields))
	return &e
}

func (e *Entry) Prelog() {
	pc, fileName, lineNum, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	e.LogrusEntry = e.LogrusEntry.WithFields(logrus.Fields{
		"_loc":  fmt.Sprintf("%s:%d", path.Base(fileName), lineNum),
		"_func": funcName,
	})
}

func (e *Entry) PrelogFatal() {
	pc, fileName, lineNum, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	e.LogrusEntry = e.LogrusEntry.WithFields(logrus.Fields{
		"_loc":   fmt.Sprintf("%s:%d", path.Base(fileName), lineNum),
		"_func":  funcName,
		"_trace": string(debug.Stack()),
	})
}

func (e *Entry) Debug(args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Debug(args...)
}

func (e *Entry) Debugf(format string, args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Debugf(format, args...)
}

func (e *Entry) Info(args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Info(args...)
}

func (e *Entry) Infof(format string, args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Infof(format, args...)
}

func (e *Entry) Warn(args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Warn(args...)
}

func (e *Entry) Warnf(format string, args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Warnf(format, args...)
}

func (e *Entry) Error(args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Error(args...)
}

func (e *Entry) Errorf(format string, args ...interface{}) {
	e.Prelog()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Errorf(format, args...)
}

func (e *Entry) Fatal(args ...interface{}) {
	e.PrelogFatal()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Time = time.Now()
	entry.Level = logrus.FatalLevel
	entry.Message = fmt.Sprint(args...)
	entry.Logger.Hooks.Fire(logrus.FatalLevel, entry)
	serialized, _ := entry.Logger.Formatter.Format(entry)
	entry.Logger.Out.Write(serialized)
}

func (e *Entry) Fatalf(format string, args ...interface{}) {
	e.PrelogFatal()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Time = time.Now()
	entry.Level = logrus.FatalLevel
	entry.Message = fmt.Sprintf(format, args...)
	entry.Logger.Hooks.Fire(logrus.FatalLevel, entry)
	serialized, _ := entry.Logger.Formatter.Format(entry)
	entry.Logger.Out.Write(serialized)
}

func (e *Entry) Panic(args ...interface{}) {
	e.PrelogFatal()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Panic(args...)
}

func (e *Entry) Panicf(format string, args ...interface{}) {
	e.PrelogFatal()
	entry := e.LogrusEntry
	if e.Pre != nil {
		entry = entry.WithFields(logrus.Fields(*e.Pre.PreLog()))
	}
	entry.Panicf(format, args...)
}
