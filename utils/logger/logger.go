package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

var (
	color_NORMAL  = "\033[0m"
	color_GRAY    = "\033[1;30m"
	color_RED     = "\033[1;31m"
	color_GREEN   = "\033[1;32m"
	color_YELLOW  = "\033[1;33m"
	color_BLUE    = "\033[1;34m"
	color_MAGENTA = "\033[1;35m"
	color_CYAN    = "\033[1;36m"
	color_WHITE   = "\033[1;37m"
)

type LogFormatter struct{}

func (slf *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var logColor = color_NORMAL

	// colorful log
	switch entry.Level {
	case logrus.DebugLevel:
		logColor = color_BLUE
	case logrus.InfoLevel:
		logColor = color_GREEN
	case logrus.WarnLevel:
		logColor = color_YELLOW
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		logColor = color_RED
	default:
	}
	timestamp := time.Now().Local().Format("06/01/02 15:04:05")

	// log format:
	// ----------------------------------------
	// [level] timestamp (file:line): field1=value1 field2=value2 ...
	// message
	// ----------------------------------------
	prefix := fmt.Sprintf("%s[%s] %s (%s:%d):%s ", logColor, entry.Level,
		timestamp,
		filepath.Base(entry.Caller.File),
		entry.Caller.Line,
		color_NORMAL,
	)

	fields := ""
	for k, v := range entry.Data {
		fields += fmt.Sprintf("  %s%s%s=%v", color_CYAN, k, color_NORMAL, v)
	}

	msg := fmt.Sprintf("%s%s\n%s\n", prefix, fields, entry.Message)
	return []byte(msg), nil
}

func init() {
	// open file location reporter
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&LogFormatter{})
}
