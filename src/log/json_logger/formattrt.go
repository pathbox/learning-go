package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Formatter struct {
	json_formatter *logrus.JSONFormatter
}

func HereFormatter() *Formatter {
	return &Formatter{
		json_formatter: &logrus.JSONFormatter{
			// FieldMap: logrus.FieldMap{
			// 	logrus.FieldKeyTime:  "_time",
			// 	logrus.FieldKeyLevel: "_level",
			// 	//logrus.FieldKeyMsg:   "@msg",
			// },
			TimestampFormat: time.RFC3339Nano,
		},
	}
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.json_formatter.Format(entry)
}
