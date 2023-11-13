package helpers

import (
	"path/filepath"
	"runtime"
	"strings"
)

func GetFunctionName() string {
	pc := make([]uintptr, 8)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	splitted := strings.Split(frame.Func.Name(), ".")
	return splitted[len(splitted)-1]
}

func GetCallerDirectory() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	return filepath.Dir(filename)
}
