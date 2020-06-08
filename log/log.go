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
type Fields map[string]interface{}

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

// Add fields to log
func Add(fields Fields) {
	var lFields logrus.Fields

	for k, v := range fields {
		lFields[k] = v
	}

	Entry.WithFields(lFields)
}

// Details of log
func Details() {
	Entry = logrus.WithFields(logrus.Fields{
		"line": line(),
		"file": file(),
	})
}

// Line number
func line() int {
	_, _, line, _ := runtime.Caller(2)
	return line
}

// File name
func file() string {
	_, filePath, _, _ := runtime.Caller(2)
	file := path.Base(filePath)

	return file
}
