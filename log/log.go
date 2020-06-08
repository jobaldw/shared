package log

import (
	"path"
	"runtime"
)

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
