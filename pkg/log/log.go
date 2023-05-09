package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

const (
	LEVEL_DEBUG = "debug"
	LEVEL_INFO  = "info"
	LEVEL_WARN  = "warn"
	LEVEL_ERROR = "error"
	LEVEL_FATAL = "fatal"
)

var red = color.New(color.FgRed).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func writef(level, prefix string, code int, msg string, args ...interface{}) {
	var colorPrefix string
	switch level {
	case LEVEL_DEBUG:
		colorPrefix = magenta("●", " [", prefix, "]")
	case LEVEL_INFO:
		colorPrefix = cyan("ℹ", " [", prefix, "]")
	case LEVEL_WARN:
		colorPrefix = yellow("★", " [", prefix, "]")
	case LEVEL_ERROR:
		colorPrefix = red("⚠", " [", prefix, "/", code, "]")
	case LEVEL_FATAL:
		colorPrefix = red("✝︎", " [", prefix, "/", code, "]")
	}

	fmt.Println(colorPrefix, fmt.Sprintf(msg, args...))
}

func Debug(prefix, msg string) { writef(LEVEL_DEBUG, prefix, 0, msg) }
func Info(prefix, msg string)  { writef(LEVEL_INFO, prefix, 0, msg) }
func Warn(prefix, msg string)  { writef(LEVEL_WARN, prefix, 0, msg) }
func Error(prefix, msg string) { writef(LEVEL_ERROR, prefix, 0, msg) }
func Fatal(prefix, msg string) {
	writef(LEVEL_FATAL, prefix, 0, msg)
	os.Exit(1)
}

func Debugf(prefix, msg string, args ...interface{}) { writef(LEVEL_DEBUG, prefix, 0, msg, args...) }
func Infof(prefix, msg string, args ...interface{})  { writef(LEVEL_INFO, prefix, 0, msg, args...) }
func Warnf(prefix, msg string, args ...interface{})  { writef(LEVEL_WARN, prefix, 0, msg, args...) }
func Errorf(prefix, msg string, args ...interface{}) { writef(LEVEL_ERROR, prefix, 0, msg, args...) }
func Fatalf(prefix, msg string, args ...interface{}) {
	writef(LEVEL_FATAL, prefix, 0, msg, args...)
	os.Exit(1)
}
