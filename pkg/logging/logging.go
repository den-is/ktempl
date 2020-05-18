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

func ParseUserLvl(level string) logrus.Level {

	logrus_lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Log(logrus.WarnLevel, fmt.Sprintf("Was not able to parse user provided log level %q. Will use %q.", level, "error"))
		return logrus.ErrorLevel
	}
	return logrus_lvl

}

func LogWithFields(fields Fields, level string, args ...interface{}) {

	logrus_fields := logrus.Fields(fields)
	entry := logger.WithFields(logrus_fields)
	entry.Log(ParseUserLvl(level), args...)

}

func Log(level string, args ...interface{}) {

	logger.Log(ParseUserLvl(level), args...)

}
