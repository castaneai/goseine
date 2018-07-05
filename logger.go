package goseine

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type LogFormatter struct {
	category string
}

func NewLogger(category string) *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &LogFormatter{category: category}
	return logger
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s][%s][%s]%s",
		entry.Time.Format("2006/01/02 15:04:05.00000"), entry.Level, f.category, entry.Message)), nil
}
