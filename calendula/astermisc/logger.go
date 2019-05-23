package astermisc

import (
	"fmt"
)

// Logger is a logger interface compatible with both stdlib and some
// 3rd party loggers.
type Logger interface {
	Output(int, string) error
}

// Ilogger represents the internal logging api we use.
type Ilogger interface {
	Logger
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type debug interface {
	Debug() bool

	// Debugf print a formatted debug line.
	Debugf(format string, v ...interface{})
	// Debugln print a debug line.
	Debugln(v ...interface{})
}

// internalLog implements the additional methods used by our internal logging.
type internalLog struct {
	Logger
}

// Println replicates the behaviour of the standard logger.
func (t internalLog) Println(v ...interface{}) {
	t.Output(2, fmt.Sprintln(v...))
}

// Printf replicates the behaviour of the standard logger.
func (t internalLog) Printf(format string, v ...interface{}) {
	t.Output(2, fmt.Sprintf(format, v...))
}

// Print replicates the behaviour of the standard logger.
func (t internalLog) Print(v ...interface{}) {
	t.Output(2, fmt.Sprint(v...))
}

type discard struct{}

func (t discard) Debug() bool {
	return false
}

// Debugf print a formatted debug line.
func (t discard) Debugf(format string, v ...interface{}) {}

// Debugln print a debug line.
func (t discard) Debugln(v ...interface{}) {}
