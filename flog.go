package flog

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasttemplate"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var gormSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = regexp.MustCompile(`utils.utils\.go`).ReplaceAllString(file, "")
}

func fileWithLineNum() string {
	// the second caller usually from gorm internal, so set i start from 2
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// LogLevel log level
type LogLevel int
type H map[string]interface{}

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	ErrorLevel
	// Warn warn log level
	WarnLevel
	// Info info log level
	InfoLevel
	// Debug debug log level
	DebugLevel
	// Trace debug log level
	TraceLevel
)

// Config logger config
type Config struct {
	Json       bool
	Format     string
	TimeFormat string
	LogLevel   LogLevel
	FullPath   bool
}

type Interface interface {
	DumpJson()
	SetConfig(cfg Config)
	SetLevel(LogLevel)
	SetFormat(string)
	SetTimeFormat(string)
	SetFullPath(bool)

	With(fields Fields) *Entry
	Trace(a ...any)
	Tracef(format string, a ...any)
	Debug(a ...any)
	Debugf(format string, a ...any)
	Info(a ...any)
	Infof(format string, a ...any)
	Warn(a ...any)
	Warnf(format string, a ...any)
	Error(a ...any)
	Errorf(format string, a ...any)
	Fatal(a ...any)
	Fatalf(format string, a ...any)
	Panic(a ...any)
	Panicf(format string, a ...any)
}

type logger struct {
	mu MutexWrap

	entryPool sync.Pool
	Config
	traceStr, debugStr, infoStr, warnStr,
	errStr, fatalStr, panicStr string
	writeFields bool
}

var (
	defaultFormat     = "${time} ${level} ${path} ${msg}\t${fields}"
	defaultTimeFormat = "2006/01/02 15:04:05"
	std               = New(Config{
		Json:       false,
		Format:     defaultFormat,
		LogLevel:   TraceLevel,
		TimeFormat: defaultTimeFormat,
		FullPath:   false,
	})
)

func DumpJson() {
	std.DumpJson()
}
func SetLevel(level LogLevel) {
	std.SetLevel(level)
}

func SetFormat(format string) {
	std.SetFormat(format)
}

func SetTimeFormat(format string) {
	std.SetTimeFormat(format)
}

func SetFullPath(fullPath bool) {
	std.SetFullPath(fullPath)
}

func With(fields Fields) *Entry {
	return std.With(fields)
}
func Trace(a ...any) {
	std.Trace(a...)
}

func Tracef(format string, a ...any) {
	std.Tracef(format, a...)
}

func Debug(a ...any) {
	std.Debug(a...)
}

func Debugf(format string, a ...any) {
	std.Debugf(format, a...)
}

func Info(a ...any) {
	std.Info(a...)
}

func Infof(format string, a ...any) {
	std.Infof(format, a...)
}

func Warn(a ...any) {
	std.Warn(a...)
}

func Warnf(format string, a ...any) {
	std.Warnf(format, a...)
}

func Error(a ...any) {
	std.Error(a...)
}

func Errorf(format string, a ...any) {
	std.Errorf(format, a...)
}

func Fatal(a ...any) {
	std.Fatal(a...)
}

func Fatalf(format string, a ...any) {
	std.Fatalf(format, a...)
}

func Panic(a ...any) {
	std.Panic(a...)
}

func Panicf(format string, a ...any) {
	std.Panicf(format, a...)
}

var (
	traceStr = White + "[TRACE]" + Reset
	debugStr = White + "[DEBUG]" + Reset
	infoStr  = Cyan + "[INFO]" + Reset
	warnStr  = Yellow + "[WARN]" + Reset
	errStr   = Red + "[ERROR]" + Reset
	fatalStr = Red + "[FATAL]" + Reset
	panicStr = Red + "[PANIC]" + Reset

	levelStrMap = map[string]string{
		traceStr: "TRACE",
		debugStr: "DEBUG",
		infoStr:  "INFO",
		warnStr:  "WARN",
		errStr:   "ERROR",
		fatalStr: "FATAL",
		panicStr: "PANIC",
	}
)

func New(config Config) Interface {
	return &logger{
		Config:      config,
		traceStr:    traceStr,
		debugStr:    debugStr,
		infoStr:     infoStr,
		warnStr:     warnStr,
		errStr:      errStr,
		fatalStr:    fatalStr,
		panicStr:    panicStr,
		writeFields: strings.Contains(config.Format, "${fields}"),
	}
}

func Sprintf(format string, h map[string]any) string {
	t := fasttemplate.New(format, "${", "}")
	return t.ExecuteString(h)
}
func (l *logger) DumpJson() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Json = true
}

func (l *logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.LogLevel = level
}
func (l *logger) SetFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Format = format
}
func (l *logger) SetTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.TimeFormat = format
}
func (l *logger) SetFullPath(fullPath bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.FullPath = fullPath
}

func (l *logger) SetConfig(cfg Config) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Config = cfg
}

func (l *logger) With(fields Fields) *Entry {
	return NewEntry(l).With(fields)
}

// Trace print messages
func (l *logger) Trace(a ...any) {
	if l.LogLevel >= TraceLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.traceStr, f, nil, a...))
	}
}
func (l *logger) Tracef(format string, a ...any) {
	if l.LogLevel >= TraceLevel {
		fmt.Println(l.sprintf(l.traceStr, format, nil, a...))
	}
}

// Debug print messages
func (l *logger) Debug(a ...any) {
	if l.LogLevel >= DebugLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.debugStr, f, nil, a...))
	}
}
func (l *logger) Debugf(format string, a ...interface{}) {
	if l.LogLevel >= DebugLevel {
		fmt.Println(l.sprintf(l.debugStr, format, nil, a...))
	}
}

// Info print info
func (l *logger) Info(a ...any) {
	if l.LogLevel >= InfoLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.infoStr, f, nil, a...))
	}
}

func (l *logger) Infof(format string, a ...any) {
	if l.LogLevel >= InfoLevel {
		fmt.Println(l.sprintf(l.infoStr, format, nil, a...))
	}
}

// Warn print warn messages
func (l *logger) Warn(a ...any) {
	if l.LogLevel >= WarnLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.warnStr, f, nil, a...))
	}
}
func (l *logger) Warnf(format string, a ...any) {
	if l.LogLevel >= WarnLevel {
		fmt.Println(l.sprintf(l.warnStr, format, nil, a...))
	}
}

// Error print error messages
func (l *logger) Error(a ...any) {
	if l.LogLevel >= ErrorLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.errStr, f, nil, a...))
	}
}
func (l *logger) Errorf(format string, a ...any) {
	if l.LogLevel >= ErrorLevel {
		fmt.Println(l.sprintf(l.errStr, format, nil, a...))
	}
}

// Fatal print error messages
func (l *logger) Fatal(a ...any) {
	if l.LogLevel >= ErrorLevel {
		f := l.formatData(a...)
		fmt.Println(l.sprintf(l.fatalStr, f, nil, a...))
		os.Exit(1)
	}
}
func (l *logger) Fatalf(format string, a ...any) {
	if l.LogLevel >= ErrorLevel {
		fmt.Println(l.sprintf(l.fatalStr, format, nil, a...))
		os.Exit(1)
	}
}

func (l *logger) Panic(a ...any) {
	if l.LogLevel >= ErrorLevel {
		f := l.formatData(a...)
		r := l.sprintf(l.panicStr, f, nil, a...)
		fmt.Println(r)
		panic(r)
	}
}
func (l *logger) Panicf(format string, a ...any) {
	if l.LogLevel >= ErrorLevel {
		r := l.sprintf(l.panicStr, format, nil, a...)
		fmt.Println(r)
		panic(errors.New(r))
	}
}

func (l *logger) formatData(a ...any) string {
	f := ""
	for argNum := range a {
		if argNum > 0 {
			f += " "
		}
		f += "%v"
	}
	return f
}
func (l *logger) sprintf(levelStr string, format string, fields any, a ...any) string {
	path := l.getPath()
	data := map[string]any{
		"level": levelStr,
		"time":  l.t(),
		"path":  path,
		"msg":   fmt.Sprintf(format, a...),
	}
	if fields != nil && l.writeFields {
		if l.Json {
			for k, v := range fields.(Fields) {
				data[k] = v
			}
			delete(data, "fields")
		} else {
			data["fields"] = fields
		}
	}
	if l.Json {
		s, ok := levelStrMap[data["level"].(string)]
		if ok {
			data["level"] = s
		}
		jsonStr, _ := json.Marshal(data)

		return string(jsonStr)
	} else {
		return Sprintf(l.Format, data)
	}
}
func (l *logger) t() string {
	return time.Now().Format(l.TimeFormat)
}

func (l *logger) getPath() string {
	path := fileWithLineNum()
	if !l.FullPath {
		arr := strings.Split(path, "/")
		path = arr[len(arr)-1]
	}
	return path
}

type MutexWrap struct {
	lock     sync.Mutex
	disabled bool
}

func (mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}
