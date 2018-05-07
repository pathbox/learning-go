package logger

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
	prefix string
	pid    int
}

func New(prefix string) *Logger {
	l := &Logger{logger: log.New(os.Stdout, "", log.LstdFlags)}
	l.prefix = fmt.Sprintf("[%s - %d]", prefix, os.Getpid())
	return l
}

func (l *Logger) Println(args ...interface{}) {
	l.logger.Println(append([]interface{}{l.prefix}, args...)...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Printf(l.prefix+" "+format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Println(args...)
	os.Exit(-1)
}
