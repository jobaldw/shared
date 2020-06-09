package log

import (
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	logrusStack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

// Entry struct
var Entry *logrus.Entry

// Logger struct
type Logger struct {
	Name  string
	Level string
	Stack bool
}

// Configure logger
func Configure(conf Logger) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		lvl = logrus.ErrorLevel
	}

	logrus.SetLevel(lvl)

	callerLevels := []logrus.Level{logrus.PanicLevel}
	stackLevels := logrus.AllLevels
	if !conf.Stack {
		stackLevels = []logrus.Level{logrus.PanicLevel}
	}

	logrus.AddHook(logrusStack.NewHook(callerLevels, stackLevels))

	Entry = logrus.WithFields(logrus.Fields{
		"application": conf.Name,
	})
}

// Details of log
func Details() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return Entry
	}

	fileName := path.Base(file) + ":" + strconv.Itoa(line)

	function := runtime.FuncForPC(pc).Name()
	funcName := function[strings.LastIndex(function, ".")+1:] + "()"

	Entry = logrus.WithField("file", fileName).WithField("function", funcName)
	return Entry
}
