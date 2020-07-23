package goslog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewLogger(prefix string) logrus.FieldLogger {
	formatter := &logrus.TextFormatter{ForceColors: true}
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(formatter)
	logger.AddHook(&prefixHook{prefix: prefix})
	return logger
}

type prefixHook struct {
	prefix string
}

func (p *prefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (p *prefixHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] %s", p.prefix, entry.Message)
	return nil
}
