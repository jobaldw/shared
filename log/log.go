package log

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// Entry struct
var Entry *logrus.Entry

// Fields type
type Fields map[string]interface{}

// Configure logger
func Configure(app, level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.ErrorLevel
	}

	logrus.SetLevel(lvl)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// Base log field that will be in every log message
	Entry = logrus.WithField("application", app)
}

// File function
func File() *logrus.Entry {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return Entry
	}

	fileName := path.Base(file) + ":" + strconv.Itoa(line)

	return Entry.WithField("file", fileName)
}
