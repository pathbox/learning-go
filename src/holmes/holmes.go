package holmes

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sync/atomic"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota

	INFO

	WARN

	ERROR

	FATAL
)

var (
	started        int32
	loggerInstance Logger
	tagName        = map[LogLevel]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
)

func Start(decorators ...func(Logger) Logger) Logger {
	if atomic.CompareAndSwapInt32(&started, 0, 1) {
		loggerInstance = Logger{}
		for _, decorator := range decorators {
			loggerInstance = decorator(loggerInstance)
		}
		var logger *log.Logger

		var segment *logSegment
		if loggerInstance.logPath != "" {
			segment = newLogSegment(loggerInstance.unit, loggerInstance.logPath)
		}
		if segment != nil {
			logger = log.New(segment, "", log.LstdFlags)
		} else if loggerInstance.isStdout {
			logger = log.New(os.Stdout, "", log.LstdFlags)
		} else {
			logger = log.New(os.Stderr, "", log.LstdFlags)
		}
		loggerInstance.logger = logger
		return loggerInstance
	}
	panic("Start() already called")
}

func (l Logger) Stop() {
	if atomic.CompareAndSwapInt32(&l.stopped, 0, 1) {
		if l.printStack {
			traceInfo := make([]byte, 1<<16)
			n := runtime.Stack(traceInfo, true)
			l.logger.Printf("%s", traceInfo[:n])
			if l.isStdout {
				log.Printf("%s", traceInfo[:n])
			}
		}
		if l.segment != nil {
			l.segment.Close()
		}
		l.segment = nil
		l.logger = inl
		atomic.StoreInt32(&started, 0)
	}
}

type logSegment struct {
	unit         time.Duration
	logPath      string
	logFile      *os.File
	timeToCreate <-chan time.Time
}

func newLogSegment(unit time.Duration, logPath string) *logSegment {
	now := time.Now()
	if logPath != "" {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		name := getLogFileName(time.Now())
		logFile, err := os.OpenFile(path.Join(logPath, name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			if os.IsNotExist(err) {
				logFile, err = os.Create(path.Join(logPath, name))
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					return nil
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
				return nil
			}
		}
		next := now.Truncate(unit).Add(unit)
		var timeToCreate <-chan time.Time
		if unit == time.Hour || unit == time.Minute {
			timeToCreate = time.After(next.Sub(time.Now()))
		}
		return &logSegment{
			unit:         unit,
			logPath:      logPath,
			logFile:      logFile,
			timeToCreate: timeToCreate,
		}
	}
	return nil
}

func (ls *logSegment) Close() {
	ls.logFile.Close()
}

func (ls *logSegment) Write(p []byte) (n int, err error) {
	if ls.timeToCreate != nil && ls.logFile != os.Stdout && ls.logFile != os.Stderr {
		select {
		case current := <-ls.timeToCreate: // 日志文件切割时间到了，就新建一个日志文件，在这个新的日志文件里进行写操作
			ls.logFile.Close()
			ls.logFile = nil
			name := getLogFileName(current)
			ls.logFile, err = os.Create(path.Join(ls.logFile, name))
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				ls.logFile = os.Stderr
			} else {
				next := current.Truncate(ls.unit).Add(ls.unit)
				ls.timeToCreate = time.After(next.Sub(time.Now()))
			}
		default:
		}
	}
	return ls.logFile.Write(p) // 返回的时候进行一个写函数操作
}

func getLogFileName(t time.Time) string {
	proc := path.Base(os.Args[0])
	now := time.Now()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()
	pid := os.Getpid()
	return fmt.Sprintf("%s.%04d-%02d-%02d-%02d-%02d.%d.log",
		proc, year, month, day, hour, minute, pid)
}

// Logger is the logger type.
type Logger struct {
	logger     *log.Logger
	level      LogLevel
	segment    *logSegment
	stopped    int32
	logPath    string
	unit       time.Duration
	isStdout   bool
	printStack bool
}

func (l Logger) doPrintf(level LogLevel, format string, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if level >= l.level {
		funcName, fileName, lineNum := getRuntimeInfo()
		format = fmt.Sprintf("%5s [%s] (%s:%d) - %s", tagName[level], path.Base(funcName), path.Base(fileName), lineNum, format)
		l.logger.Printf(format, v...)
		if l.isStdout {
			log.Printf(format, v...)
		}
		if level == FATAL {
			os.Exit(1)
		}
	}
}
func (l Logger) doPrintln(level LogLevel, v ...interface{}) {
	if l.logger == nil {
		return
	}
	if level >= l.level {
		funcName, fileName, lineNum := getRuntimeInfo()
		prefix := fmt.Sprintf("%5s [%s] (%s:%d) - ", tagName[level], path.Base(funcName), path.Base(fileName), lineNum)
		value := fmt.Sprintf("%s%s", prefix, fmt.Sprintln(v...))
		l.logger.Print(value)
		if l.isStdout {
			log.Print(value)
		}
		if level == FATAL {
			os.Exit(1)
		}
	}
}

func getRuntimeInfo() (string, string, int) {
	pc, fn, ln, ok := runtime.Caller(3) // 3 steps up the stack frame
	if !ok {
		fn = "???"
		ln = 0
	}
	function := "???"
	caller := runtime.FuncForPC(pc)
	if caller != nil {
		function = caller.Name()
	}
	return function, fn, ln
}

// DebugLevel sets log level to debug.
func DebugLevel(l Logger) Logger {
	l.level = DEBUG
	return l
}

// InfoLevel sets log level to info.
func InfoLevel(l Logger) Logger {
	l.level = INFO
	return l
}

// WarnLevel sets log level to warn.
func WarnLevel(l Logger) Logger {
	l.level = WARN
	return l
}

// ErrorLevel sets log level to error.
func ErrorLevel(l Logger) Logger {
	l.level = ERROR
	return l
}

// FatalLevel sets log level to fatal.
func FatalLevel(l Logger) Logger {
	l.level = FATAL
	return l
}

// LogFilePath returns a function to set the log file path.
func LogFilePath(p string) func(Logger) Logger {
	return func(l Logger) Logger {
		l.logPath = p
		return l
	}
}

// EveryHour sets new log file created every hour.
func EveryHour(l Logger) Logger {
	l.unit = time.Hour
	return l
}

// EveryMinute sets new log file created every minute.
func EveryMinute(l Logger) Logger {
	l.unit = time.Minute
	return l
}

// AlsoStdout sets log also output to stdio.
func AlsoStdout(l Logger) Logger {
	l.isStdout = true
	return l
}

// PrintStack sets log output the stack trace info.
func PrintStack(l Logger) Logger {
	l.printStack = true
	return l
}

// Debugf prints formatted debug log.
func Debugf(format string, v ...interface{}) {
	loggerInstance.doPrintf(DEBUG, format, v...)
}

// Infof prints formatted info log.
func Infof(format string, v ...interface{}) {
	loggerInstance.doPrintf(INFO, format, v...)
}

// Warnf prints formatted warn log.
func Warnf(format string, v ...interface{}) {
	loggerInstance.doPrintf(WARN, format, v...)
}

// Errorf prints formatted error log.
func Errorf(format string, v ...interface{}) {
	loggerInstance.doPrintf(ERROR, format, v...)
}

// Fatalf prints formatted fatal log and exits.
func Fatalf(format string, v ...interface{}) {
	loggerInstance.doPrintf(FATAL, format, v...)
	os.Exit(1)
}

// Debugln prints debug log.
func Debugln(v ...interface{}) {
	loggerInstance.doPrintln(DEBUG, v...)
}

// Infoln prints info log.
func Infoln(v ...interface{}) {
	loggerInstance.doPrintln(INFO, v...)
}

// Warnln prints warn log.
func Warnln(v ...interface{}) {
	loggerInstance.doPrintln(WARN, v...)
}

// Errorln prints error log.
func Errorln(v ...interface{}) {
	loggerInstance.doPrintln(ERROR, v...)
}

// Fatalln prints fatal log and exits.
func Fatalln(v ...interface{}) {
	loggerInstance.doPrintln(FATAL, v...)
	os.Exit(1)
}
