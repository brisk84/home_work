package logger

import (
	"io/fs"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	path  string
	level string
}

func New(path string, level string) *Logger {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)
	var logOut *os.File
	switch path {
	case "stderr":
		logOut = os.Stderr
	case "stdout":
		logOut = os.Stdout
	default:
		flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
		fileMode := 755
		logOut, err = os.OpenFile(path, flag, fs.FileMode(fileMode))
		if err != nil {
			logOut = os.Stdout
		}
	}
	logrus.SetOutput(logOut)
	return &Logger{path: path, level: level}
}

func (l Logger) Info(msg string) {
	logrus.Info(msg)
}

func (l Logger) Error(msg string) {
	logrus.Error(msg)
}
