package log

import (
	"io"
	l "log"
	"os"
)

type level uint8

const (
	ver level = 0
	deb level = 1
	inf level = 2
	war level = 3
	err level = 4
	fat level = 5
	off level = 6
)

var lvl = inf

// --- Initializer ---

// init changes the default output to stdout.
func init() {
	SetOutput(os.Stdout)
}

// --- Configuration functions ---

// SetOutput sets the output destination for the logger.
func SetOutput(w io.Writer) {
	l.SetOutput(w)
}

// SetLevel set the logging level for the logger. (Valid levels are: "verbose", "debug", "info",
// "warning", "error", "fatal", "off")
func SetLevel(l string) {
	switch l {
	case "verbose":
		lvl = ver
	case "debug":
		lvl = deb
	case "info":
		lvl = inf
	case "warning":
		lvl = war
	case "error":
		lvl = err
	case "fatal":
		lvl = fat
	case "off":
		lvl = off
	default:
		Fatalf("Invalid logging level '%s'!", l)
	}
}

// --- Logging functions ---

// Verb writes a verbose log message.
func Verb(msg string) {
	if lvl <= ver {
		l.Print("[VERB ]: " + msg)
	}
}

// Verbf writes a verbose log message and fills placeholders according to fmt.Printf().
func Verbf(msg string, v ...interface{}) {
	if lvl <= ver {
		l.Printf("[VERB ]: "+msg, v...)
	}
}

// Debug writes a debug log message.
func Debug(msg string) {
	if lvl <= deb {
		l.Print("[DEBUG]: " + msg)
	}
}

// Debugf writes a debug log message and fills placeholders according to fmt.Printf().
func Debugf(msg string, v ...interface{}) {
	if lvl <= deb {
		l.Printf("[DEBUG]: "+msg, v...)
	}
}

// Info writes a info log message.
func Info(msg string) {
	if lvl <= inf {
		l.Print("[INFO ]: " + msg)
	}
}

// Infof writes a info log message and fills placeholders according to fmt.Printf().
func Infof(msg string, v ...interface{}) {
	if lvl <= inf {
		l.Printf("[INFO ]: "+msg, v...)
	}
}

// Warn writes a warning log message.
func Warn(msg string) {
	if lvl <= war {
		l.Print("[WARN ]: " + msg)
	}
}

// Warnf writes a warning log message and fills placeholders according to fmt.Printf().
func Warnf(msg string, v ...interface{}) {
	if lvl <= war {
		l.Printf("[WARN ]: "+msg, v...)
	}
}

// Error writes a error log message.
func Error(msg string) {
	if lvl <= err {
		l.Print("[ERROR]: " + msg)
	}
}

// Errorf writes a error log message and fills placeholders according to fmt.Printf().
func Errorf(msg string, v ...interface{}) {
	if lvl <= err {
		l.Printf("[ERROR]: "+msg, v...)
	}
}

// Fatal writes a fatal log message followed by a call to os.Exit(1).
func Fatal(msg string) {
	if lvl <= fat {
		l.Fatal("[FATAL]: " + msg)
	}
}

// Fatalf writes a fatal log message and fills placeholders according to fmt.Printf() followed by a
// call to os.Exit(1).
func Fatalf(msg string, v ...interface{}) {
	if lvl <= fat {
		l.Fatalf("[FATAL]: "+msg, v...)
	}
}
