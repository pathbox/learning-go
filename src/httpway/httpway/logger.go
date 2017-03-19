package httpway

import (
	"fmt"
)

type Logger interface {
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Debug(v ...interface{})
}

const default_depth = 5

type internalLogger struct {
	l      Logger
	id     uint64
	prefix string
}

func (il *internalLogger) Warning(v ...interface{}) {
	if il.l == nil {
		return
	}
	il.log(il.l.Warning, v...)
}

func (il *internalLogger) Error(v ...interface{}) {
	if il.l == nil {
		return
	}

	il.log(il.l.Error, v...)
}

func (il *internalLogger) Info(v ...interface{}) {
	if il.l == nil {
		return
	}
	il.log(il.l.Info, v...)
}

func (il *internalLogger) Debug(v ...interface{}) {
	if il.l == nil {
		return
	}
	il.log(il.l.Debug, v...)
}

// 这里很有意思， 你可以看成是一个 binding。 block f，和 值 v。 对值v进行一些操作后， 放入blcok f。这里block 就是一个函数，而且这个函数是当成参数一样传过来，是可变的。这样 log方法就可以根据不同的f参数而执行不同的操作
func (il *internalLogger) log(f func(v ...interface{}), v ...interface{}) {
	if len(v) > 1 {
		v[0] = fmt.Sprintf("[%x] %s", il.id, v[0].(string))
	} else {
		v[0] = fmt.Sprintf("[%x] %v", il.id, v[0])
	}

	if il.prefix != "" {
		v[0] = fmt.Sprintf("[%s] %s", il.prefix, v[0].(string))
	}

	f(v...)
}

type internalServerLoggerWriter struct {
	l Logger
}

func (islw *internalServerLoggerWriter) Write(b []byte) (n int, err error) {
	islw.l.Error("Go Server: %s", string(b))
	return len(b), nil
}
