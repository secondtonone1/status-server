package logging

import "fmt"

type BaseLogger interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Debug(args ...interface{})
	Warn(args ...interface{})
	Panic(args ...interface{})
	DPanic(args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	DPanicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	SetLevel(level string)
}

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	DPanicLevel
	PanicLevel
	FatalLevel
)

type ConsoleLogger struct {
	level int
}

func (cl *ConsoleLogger) Info(args ...interface{}) {
	if cl.level <= InfoLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Error(args ...interface{}) {
	if cl.level <= ErrorLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Debug(args ...interface{}) {
	if cl.level <= DebugLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Warn(args ...interface{}) {
	if cl.level <= WarnLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) DPanic(args ...interface{}) {
	if cl.level <= DPanicLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Panic(args ...interface{}) {
	if cl.level <= PanicLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Fatal(args ...interface{}) {
	if cl.level <= FatalLevel {
		fmt.Println(args)
	}
}

func (cl *ConsoleLogger) Infof(format string, args ...interface{}) {
	if cl.level <= InfoLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) Errorf(format string, args ...interface{}) {
	if cl.level <= ErrorLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) Debugf(format string, args ...interface{}) {
	if cl.level <= DebugLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) Warnf(format string, args ...interface{}) {
	if cl.level <= WarnLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) DPanicf(format string, args ...interface{}) {
	if cl.level <= DPanicLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) Panicf(format string, args ...interface{}) {
	if cl.level <= PanicLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) Fatalf(format string, args ...interface{}) {
	if cl.level <= FatalLevel {
		fmt.Printf(format, args...)
	}
}

func (cl *ConsoleLogger) SetLevel(level string) {
	switch level {
	case "debug":
		cl.level = 0
	case "info":
		cl.level = 1
	case "warn":
		cl.level = 2
	case "error":
		cl.level = 3
	default:
		cl.level = 1
	}
}
