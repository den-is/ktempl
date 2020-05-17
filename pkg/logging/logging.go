package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggingConfig struct {
	Level string
	File  string
}

type Fields map[string]interface{}

var logger = logrus.New()

func LoggerSetup(config *LoggingConfig) {

	if config.File != "" {
		f, errOpen := os.OpenFile(config.File, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if errOpen != nil {
			fmt.Println("was not able to open log file", config.File)
		}
		logger.SetOutput(f)
	}

	level, errLevel := logrus.ParseLevel(config.Level)
	if errLevel != nil {
		fmt.Println("was unable to set log level", config.Level)
	}
	logger.SetLevel(level)

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

}

func LogWithFields(fields Fields, level string, args ...interface{}) {

	logrus_fields := logrus.Fields(fields)

	entry := logger.WithFields(logrus_fields)
	logrus_lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Log(logrus.FatalLevel, fmt.Sprintf("Was not able to parse user provided log level %q", level))
	}
	entry.Log(logrus_lvl, args...)

}

func Log(level string, args ...interface{}) {

	logrus_lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Log(logrus.FatalLevel, fmt.Sprintf("Was not able to parse user provided log level %q", level))
	}
	logger.Log(logrus_lvl, args...)
}
