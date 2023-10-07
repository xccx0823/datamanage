package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

var logger = logrus.New()

func Init() {
	log := logrus.New()
	log.SetFormatter(&formatter{})
	logger = log
}

type formatter struct{}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	_, file, line, ok := runtime.Caller(8)
	if !ok {
		file = "???"
		line = 0
	}
	fileName := filepath.Base(file)
	levelColor := getColorByLevel(entry.Level)

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf(
		"[%s] | "+levelColor+"%s"+"\x1b[0m"+" | %s:%d | %s\n",
		timestamp,
		strings.ToUpper(entry.Level.String()),
		fileName,
		line,
		entry.Message,
	)
	b.WriteString(newLog)
	return b.Bytes(), nil
}

func getColorByLevel(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "\x1b[34m"
	case logrus.InfoLevel:
		return "\x1b[32m"
	case logrus.WarnLevel:
		return "\x1b[33m"
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return "\x1b[31m"
	default:
		return "\x1b[0m"
	}
}

func GetLogger() *logrus.Logger {
	return logger
}
