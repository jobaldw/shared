package log

import (
	"os"
	"path"
	"runtime"

	logrusStack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

//Entry is a logger struct that can utilize the logrus package
var Entry *logrus.Entry

//Logger used to configure app logging
type Logger struct {
	App       string
	Level     logrus.Level
	ShowStack bool
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)

	Entry = logrus.WithFields(logrus.Fields{
		"application": "-",
	})
}

//Configure customizing our implementation of a logger
func Configure(logger Logger) {
	callerLevels := []logrus.Level{logrus.PanicLevel}
	stackLevels := logrus.AllLevels

	if !logger.ShowStack {
		stackLevels = []logrus.Level{logrus.PanicLevel}
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.Level)
	logrus.AddHook(logrusStack.NewHook(callerLevels, stackLevels))

	Entry = logrus.WithFields(logrus.Fields{
		"application": logger.App,
	})
}

// Trace prints out a Trace level log message.
func Trace(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Trace(msg...)
}

// Debug prints out a Debug level log message.
func Debug(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Debug(msg...)
}

// Info prints out a Info level log message.
func Info(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Info(msg...)
}

// Warn prints out a Warn level log message.
func Warn(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Warn(msg...)
}

// Error prints out a Error level log message.
func Error(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Error(msg...)
}

// Fatal prints out a Fatal level log message.
func Fatal(msg ...interface{}) {
	_, filePath, lineNum, _ := runtime.Caller(1)
	fileName := path.Base(filePath)

	Entry.WithFields(logrus.Fields{
		"file": fileName,
		"line": lineNum,
	}).Fatal(msg...)
}
