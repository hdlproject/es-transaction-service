package helper

import (
	"errors"
	"fmt"
	"runtime"
)

func getCurrentFunctionInfo() (string, int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function, frame.Line
}

func WrapError(originalErr error) error {
	functionName, lineNumber := getCurrentFunctionInfo()
	return fmt.Errorf("(%s:%d):%w", functionName, lineNumber, originalErr)
}

func GetOriginalError(wrappedErr error) error {
	for errors.Unwrap(wrappedErr) != nil {
		wrappedErr = errors.Unwrap(wrappedErr)
	}

	return wrappedErr
}

func WrapLog(originalLog string) string {
	functionName, lineNumber := getCurrentFunctionInfo()
	return fmt.Sprintf("(%s:%d):%s", functionName, lineNumber, originalLog)
}
