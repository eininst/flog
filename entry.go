package flog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var sDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	sDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
}

func entryFileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, sDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

type Fields map[string]interface{}

type Entry struct {
	logger *logger
	format string
	data   Fields
}

func NewEntry(logger *logger) *Entry {
	f := logger.Format + "  ${fields}"
	if strings.Contains(logger.Format, "${fields}") {
		f = logger.Format
	}
	return &Entry{
		logger: logger,
		format: f,
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
		fmt.Println(e.entrySprintf(e.logger.traceStr, f, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Tracef(format string, a ...any) {
	if e.logger.LogLevel >= TraceLevel {
		fmt.Println(e.entrySprintf(e.logger.traceStr, format, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Debug(a ...any) {
	if e.logger.LogLevel >= DebugLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.entrySprintf(e.logger.debugStr, f, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Debugf(format string, a ...any) {
	if e.logger.LogLevel >= DebugLevel {
		fmt.Println(e.entrySprintf(e.logger.debugStr, format, e.fieldsFormat(White), a...))
	}
}

func (e *Entry) Info(a ...any) {
	if e.logger.LogLevel >= InfoLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.entrySprintf(e.logger.infoStr, f, e.fieldsFormat(Cyan), a...))
	}
}

func (e *Entry) Infof(format string, a ...any) {
	if e.logger.LogLevel >= InfoLevel {
		fmt.Println(e.entrySprintf(e.logger.infoStr, format, e.fieldsFormat(Cyan), a...))
	}
}

func (e *Entry) Warn(a ...any) {
	if e.logger.LogLevel >= WarnLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.entrySprintf(e.logger.warnStr, f, e.fieldsFormat(Yellow), a...))
	}
}

func (e *Entry) Warnf(format string, a ...any) {
	if e.logger.LogLevel >= WarnLevel {
		fmt.Println(e.entrySprintf(e.logger.warnStr, format, e.fieldsFormat(Yellow), a...))
	}
}

func (e *Entry) Error(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.entrySprintf(e.logger.errStr, f, e.fieldsFormat(Red), a...))
	}
}

func (e *Entry) Errorf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		fmt.Println(e.entrySprintf(e.logger.errStr, format, e.fieldsFormat(Red), a...))
	}
}

func (e *Entry) Fatal(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		fmt.Println(e.entrySprintf(e.logger.fatalStr, f, e.fieldsFormat(Red), a...))
		os.Exit(1)
	}
}

func (e *Entry) Fatalf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		fmt.Println(e.entrySprintf(e.logger.fatalStr, format, e.fieldsFormat(Red), a...))
		os.Exit(1)
	}
}

func (e *Entry) Panic(a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		f := e.logger.formatData(a...)
		r := e.entrySprintf(e.logger.panicStr, f, e.fieldsFormat(Red), a...)
		fmt.Println(r)
		panic(errors.New(r))
	}
}

func (e *Entry) Panicf(format string, a ...any) {
	if e.logger.LogLevel >= ErrorLevel {
		r := e.entrySprintf(e.logger.panicStr, format, e.fieldsFormat(Red), a...)
		fmt.Println(r)
		panic(errors.New(r))
	}
}

func (e *Entry) getPath() string {
	path := entryFileWithLineNum()
	if !e.logger.FullPath {
		arr := strings.Split(path, "/")
		path = arr[len(arr)-1]
	}
	return path
}

func (e *Entry) entrySprintf(levelStr string, format string, fields any, a ...any) string {
	path := e.getPath()
	msg := fmt.Sprintf(format, a...)

	if !e.logger.Json && e.logger.MsgMinLen > 0 {
		mlen := len(msg)
		if mlen < e.logger.MsgMinLen {
			for i := 0; i < e.logger.MsgMinLen-mlen; i++ {
				msg += " "
			}
		}
	}
	data := map[string]any{
		"level": levelStr,
		"time":  e.logger.t(),
		"path":  path,
		"msg":   msg,
		"pid":   os.Getpid(),
	}

	if e.logger.Json {
		for k, v := range fields.(Fields) {
			data[k] = v
		}
		s, ok := levelStrMap[data["level"].(string)]
		if ok {
			data["level"] = s
		}
		jsonStr, _ := json.Marshal(data)

		return string(jsonStr)
	} else {
		data["fields"] = fields

		return Sprintf(e.format, data)
	}
}

func (e *Entry) fieldsFormat(color string) any {
	if e.logger.Json {
		return e.data
	}
	r := ""
	var names []string
	for name := range e.data {
		names = append(names, name)
	}
	sort.Strings(names)
	for i, name := range names {
		if i > 0 {
			r += " "
		}
		r += fmt.Sprintf("%s%s%s=%v", color, name, Reset, e.data[name])
	}
	return r
}
