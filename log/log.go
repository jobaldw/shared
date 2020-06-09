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

// Fields map
type Fields map[string]interface{}

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

// Log message
func Log() *logrus.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}

	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]

	return logrus.WithField("file", filename).WithField("function", fn)
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
