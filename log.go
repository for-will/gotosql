package db

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	tagExecSql = Green.Add("Exec Sql: ")
	tagError   = Red.Add("Error: ")
)

func LogSql(format string, a ...interface{}) {
	fmt.Printf(fileLine(1) + tagExecSql)
	fmt.Printf(format+"\n", a...)
}

func LogError(format string, a ...interface{}) {
	fmt.Printf(fileLine(0) + tagError)
	fmt.Printf(format+"\n", a...)
}

func getCallerFrame(skip int) (frame runtime.Frame, ok bool) {
	const skipOffset = 2 // skip getCallerFrame and Callers

	pc := make([]uintptr, 1)
	numFrames := runtime.Callers(skip+skipOffset, pc)
	if numFrames < 1 {
		return
	}

	frame, _ = runtime.CallersFrames(pc).Next()
	return frame, frame.PC != 0
}

func fileLine(skip int) string {
	frame, _ := getCallerFrame(2 + skip)
	file := frame.File
	line := frame.Line
	if file != "" {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	}
	var sb = new(strings.Builder)
	fmt.Fprintf(sb, "%s:%d: ", file, line)
	return sb.String()
}
