package flog

import (
	"errors"
	"fmt"
	"os"
)

type Fields map[string]interface{}

type Entry struct {
	logger *logger
	data   Fields
}

func NewEntry(logger *logger) *Entry {
	return &Entry{
		logger: logger,
		data:   make(Fields),
	}
}

func (e *Entry) With(fields Fields) *Entry {
	for k, v := range fields {
		e.data[k] = v
	}
	return e
}

func (e *Entry) Trace(a ...any) {
	if e.logger.LogLevel >= TraceLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.traceStr, f, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Tracef(format string, a ...any) {
	if e.logger.LogLevel >= TraceLevel {
		fmt.Println(e.logger.sprintf(e.logger.traceStr, format, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Debug(a ...any) {
	if e.logger.LogLevel >= DebugLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.debugStr, f, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Debugf(format string, a ...any) {
	if e.logger.LogLevel >= DebugLevel {
		fmt.Println(e.logger.sprintf(e.logger.debugStr, format, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Info(a ...any) {
	if e.logger.LogLevel >= InfoLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.infoStr, f, e.fieldsFormat(Cyan), a...))
	}
}

func (e *Entry) Infof(format string, a ...any) {
	if e.logger.LogLevel >= InfoLevel {
		fmt.Println(e.logger.sprintf(e.logger.infoStr, format, e.fieldsFormat(Cyan), a...))
	}
}

func (e *Entry) Warn(a ...any) {
	if e.logger.LogLevel >= WarnLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.warnStr, f, e.fieldsFormat(Yellow), a...))
	}
}

func (e *Entry) Warnf(format string, a ...any) {
	if e.logger.LogLevel >= WarnLevel {
		fmt.Println(e.logger.sprintf(e.logger.warnStr, format, e.fieldsFormat(Yellow), a...))
	}
}

func (e *Entry) Error(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.errStr, f, e.fieldsFormat(Red), a...))
	}
}

func (e *Entry) Errorf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		fmt.Println(e.logger.sprintf(e.logger.errStr, format, e.fieldsFormat(Red), a...))
	}
}

func (e *Entry) Fatal(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.logger.sprintf(e.logger.fatalStr, f, e.fieldsFormat(Red), a...))
		os.Exit(1)
	}
}

func (e *Entry) Fatalf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		fmt.Println(e.logger.sprintf(e.logger.fatalStr, format, e.fieldsFormat(Red), a...))
		os.Exit(1)
	}
}

func (e *Entry) Panic(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		r := e.logger.sprintf(e.logger.panicStr, f, e.fieldsFormat(Red), a...)
		fmt.Println(r)
		panic(errors.New(r))
	}
}

func (e *Entry) Panicf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		r := e.logger.sprintf(e.logger.panicStr, format, e.fieldsFormat(Red), a...)
		fmt.Println(r)
		panic(errors.New(r))
	}
}

func (e *Entry) fieldsFormat(color string) any {
	if e.logger.Json {
		return e.data
	}
	r := ""
	i := 0

	for k, v := range e.data {
		if i > 0 {
			r += " "
		}
		r += fmt.Sprintf("%s%s%s=%v", color, k, Reset, v)
		i += 1
	}
	return r
}
