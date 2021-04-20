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
	// ShellMode determines what to print to. If false use logrus, else just print straight to /dev/tty
	ShellMode = false
	success   = color.New(color.FgGreen).SprintFunc()
	info      = color.New(color.FgWhite).SprintFunc()
	warn      = color.New(color.FgYellow).SprintFunc()
	er        = color.New(color.FgRed).SprintFunc()
)

type Logger interface {
	Success(format string, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Fatal(err error)
	Name(name string)
}

// Logging defines a wrapper around sirupsen/logrus with defined name
type Logging struct {
	name     string
	log      *logrus.Logger
	LogLevel int
}

// CreateLogger will create a new *Logging that wraps around logrus.Logger with a given name
func CreateLogger(name string) *Logging {
	// NOTE: for future logrus customization
	logrusLogger := logrus.New()
	return &Logging{name: name, log: logrusLogger, LogLevel: 1}
}

// SetLoggingLevel defines level for logrus to log
func (l *Logging) SetLoggingLevel(lvl int) {
	switch lvl {
	case FATAL:
		l.log.SetLevel(logrus.FatalLevel)
		l.LogLevel = FATAL
	case WARN:
		l.log.SetLevel(logrus.WarnLevel)
		l.LogLevel = WARN
	case INFO:
		l.log.SetLevel(logrus.InfoLevel)
		l.LogLevel = INFO
	}
}

// Success is a wrapper around Info with color
func (l *Logging) Success(format string, args ...interface{}) {
	if ShellMode {
		s := fmt.Sprintf(success(format), args...)
		fmt.Printf(s)
		return
	}
	l.log.Info(fmt.Sprintf(format, args...))
}

// Name returns name of given logger
func (l *Logging) Name(n string) {
	l.name = n
}

// Info is a logrus.Info wrapper
func (l *Logging) Info(format string, args ...interface{}) {
	if ShellMode {
		s := fmt.Sprintf(info(format), args...)
		fmt.Printf(s)
		return
	}
	l.log.Info(fmt.Sprintf(format, args...))
}

// WInfo is logrus.Info wrapper for future API calls
func (l *Logging) WInfo(w http.ResponseWriter, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
	l.Info(format, args...)
}

// Warn is logrus.Warn wraps with color
func (l *Logging) Warn(format string, args ...interface{}) {
	if ShellMode {
		s := fmt.Sprintf(warn(format), args...)
		fmt.Printf(s)
		return
	}
	l.log.Warnf(format, args...)
}

// WWarn is Warn but for API calls
func (l *Logging) WWarn(w http.ResponseWriter, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
	l.Warn(format, args...)
}

// Fatal is a wrapper for logrus.Fatal
func (l *Logging) Fatal(err error) {
	if ShellMode {
		s := fmt.Sprintf(er("fatal: %s"), err.Error())
		fmt.Println(s)
		panic(err)
	}
	l.log.Fatal(err)
}
