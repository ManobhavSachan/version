package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

// SetLevel sets the logging level
func SetLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// GetLogger returns the logger instance
func GetLogger() *logrus.Logger {
	return log
}

// Info logs info level message
func Info(msg string, fields map[string]interface{}) {
	if fields != nil {
		log.WithFields(fields).Info(msg)
	} else {
		log.Info(msg)
	}
}

// Error logs error level message
func Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err
	log.WithFields(fields).Error(msg)
}

// Debug logs debug level message
func Debug(msg string, fields map[string]interface{}) {
	if fields != nil {
		log.WithFields(fields).Debug(msg)
	} else {
		log.Debug(msg)
	}
}
