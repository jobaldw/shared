package log

import (
	"os"
	"path"
	"runtime"

	logrusStack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

// Entry struct
var Entry *logrus.Entry

// Fields map
var Fields logrus.Fields

// Logger struct
type Logger struct {
	Name  string
	Level logrus.Level
	Stack bool
}

// Configure logger
func Configure(conf Logger) {
	callerLevels := []logrus.Level{logrus.PanicLevel}
	stackLevels := logrus.AllLevels

	if !conf.Stack {
		stackLevels = []logrus.Level{logrus.PanicLevel}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(conf.Level)
	logrus.AddHook(logrusStack.NewHook(callerLevels, stackLevels))

	Entry = logrus.WithFields(logrus.Fields{
		"application": conf.Name,
	})
}

// Line number
func Line() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

// File name
func File() string {
	_, filePath, _, _ := runtime.Caller(1)
	file := path.Base(filePath)

	return file
}
