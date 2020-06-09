package log

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Entry struct
var Entry *logrus.Entry

// Fields type
type Fields map[string]string

// Logger struct
type Logger struct {
	App   string
	Level string
}

// Configure logger
func Configure(conf Logger) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		level = logrus.ErrorLevel
	}

	logrus.SetLevel(level)

	// Base log field that will be in every log message
	Entry = logrus.WithFields(logrus.Fields{
		"application": conf.App,
	})
}

// Details of log
func Details() *logrus.Entry {
	return details()
}

// Add addition fields to details of log
func Add(fields Fields) *logrus.Entry {
	logrusFields := make(logrus.Fields)

	for k, v := range fields {
		logrusFields[k] = v
	}

	return details().WithFields(logrus.Fields(logrusFields))
}

func details() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return Entry
	}

	fileName := path.Base(file) + ":" + strconv.Itoa(line)
	function := runtime.FuncForPC(pc).Name()
	funcName := function[strings.LastIndex(function, ".")+1:] + "()"

	return Entry.WithField("file", fileName).WithField("function", funcName)
}
