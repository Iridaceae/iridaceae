package log

import (
	"fmt"
	"net/http"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

const (
	FATAL = 0 // fatal only
	WARN  = 1 // warn + fatal
	INFO  = 2 // all
)

var (
	// shellMode determines what to print to. If false use logrus, else just print straight to /dev/tty
	shellMode = false
	success   = color.New(color.FgGreen).SprintFunc()
	info      = color.New(color.FgWhite).SprintFunc()
	warn      = color.New(color.FgYellow).SprintFunc()
	er        = color.New(color.FgRed).SprintFunc()
)

type Logger interface {
	Success(format string, args ...interface{})
	Prompt(p string)
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Fatal(err error)
	Name(name string) Logger
}

func SetLoggingLevel(l int) {
	switch l {
	case FATAL:
		logrus.SetLevel(logrus.FatalLevel)
	case WARN:
		logrus.SetLevel(logrus.WarnLevel)
	case INFO:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func Success(format string, args ...interface{}) {
	if shellMode {
		s := fmt.Sprintf(success(format), args...)
		fmt.Printf(s)
		return
	}
	logrus.Info(fmt.Sprintf(format, args...))
}

func Prompt(p string) {
	fmt.Print(info(p))
}

func Info(format string, args ...interface{}) {
	if shellMode {
		s := fmt.Sprintf(info(format), args...)
		fmt.Printf(s)
		return
	}
	logrus.Info(fmt.Sprintf(format, args...))
}

func WInfo(w http.ResponseWriter, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
	Info(format, args...)
}

func Warn(format string, args ...interface{}) {
	if shellMode {
		s := fmt.Sprintf(warn(format), args...)
		fmt.Printf(s)
		return
	}
	logrus.Warnf(format, args...)
}

func WWarn(w http.ResponseWriter, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
	Warn(format, args...)
}

func Fatal(err error) {
	if shellMode {
		s := fmt.Sprintf(er("fatal: %s"), err.Error())
		fmt.Println(s)
		panic(err)
	}
	logrus.Fatal(err)
}
