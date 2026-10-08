package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yourusername/loganalyzer/config"
)

var log = logrus.New()

// Init configures the global logger.
func Init(cfg *config.Config) {
	// Set log level based on environment (default to info)
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if levelStr == "" {
		levelStr = "info"
	}
	if lvl, err := logrus.ParseLevel(levelStr); err == nil {
		log.SetLevel(lvl)
	}

	// Use JSON formatter for production, Text for dev
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	log.SetOutput(os.Stdout)
}

// Debug logs a debug message.
func Debug(msg string, fields logrus.Fields) {
	log.WithFields(fields).Debug(msg)
}

// Info logs an info message.
func Info(msg string, fields logrus.Fields) {
	log.WithFields(fields).Info(msg)
}

// Warn logs a warning.
func Warn(msg string, fields logrus.Fields) {
	log.WithFields(fields).Warn(msg)
}

// Error logs an error.
func Error(msg string, fields logrus.Fields) {
	log.WithFields(fields).Error(msg)
}

// Fatal logs a fatal error and exits.
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Helper wrappers that match Logrus signatures.
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}