package logger

import (
	"github.com/sirupsen/logrus"
)

type LogHandler interface {
	DoLog(e *logrus.Entry)
}

type LogHook struct {
	Handler LogHandler
}

func (l *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (l *LogHook) Fire(e *logrus.Entry) error {
	go func() {
		l.Handler.DoLog(e)
	}()
	return nil
}
