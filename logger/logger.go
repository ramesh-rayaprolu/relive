package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	//ERROR - all error logs will have this string appended
	ERROR = "ERROR"
	//INFO - all info logs will have this string appended
	INFO = " INFO"
	//DEBUG - all debug logs will have this string appended
	DEBUG = "DEBUG"
)

//Logger - main logger structure
type Logger struct {
	debugEnabled bool
}

//NewLoggerObject - create a new logger object to be used for logging
func NewLoggerObject(debugFlag bool) (*Logger, error) {
	logObj := &Logger{
		debugEnabled: debugFlag,
	}
	return logObj, nil
}

//PrintInfo - Info statements
func (l *Logger) PrintInfo(format string, v ...interface{}) {

	logMsg := fmt.Sprintf(format, v...)
	logHdr := l.generateLogHeader(INFO)

	fmt.Println(logHdr, logMsg)
	return
}

//PrintError - Error statements
func (l *Logger) PrintError(format string, v ...interface{}) {

	logMsg := fmt.Sprintf(format, v...)
	logHdr := l.generateLogHeader(ERROR)

	fmt.Println(logHdr, logMsg)
	return
}

//PrintDebug - Debug statements
func (l *Logger) PrintDebug(format string, v ...interface{}) {
	if l.debugEnabled {
		logMsg := fmt.Sprintf(format, v...)
		logHdr := l.generateLogHeader(DEBUG)

		fmt.Println(logHdr, logMsg)
	}
	return
}

//generateLogHeader - generates a header string for logging
func (l *Logger) generateLogHeader(level string) string {
	pid := os.Getpid()
	pname := os.Args[0]
	_, fn, line, _ := runtime.Caller(2)
	set := strings.Split(fn, "/")
	filename := set[len(set)-1]
	timeStamp := time.Now().UTC().Format(time.StampMilli)
	header := fmt.Sprintf("[%s] %d %s %s:%d %s - ", timeStamp, pid, pname, filename, line, level)
	return header
}

//EnableDebug - enable debug logging if required
func (l *Logger) EnableDebug() {
	l.debugEnabled = true
}
