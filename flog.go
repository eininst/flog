package flog

import (
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

func FileWithLineNum() string {
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
	Format     string
	TimeFormat string
	LogLevel   LogLevel
	FullPath   bool
}

type Interface interface {
	SetConfig(cfg Config)
	SetLevel(LogLevel)
	SetFormat(string)
	SetTimeFormat(string)
	SetFullPath(bool)

	Trace(string, ...interface{})
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
	Panic(string, ...interface{})
}

type logger struct {
	mu MutexWrap
	Config
	traceStr, debugStr, infoStr, warnStr,
	errStr, fatalStr, panicStr string
}

var (
	defaultFormat     = "${Time} ${Level} ${Path} ${Msg}"
	defaultTimeFormat = "2006/01/02 15:04:05"
	std               = New(Config{
		Format:     defaultFormat,
		LogLevel:   TraceLevel,
		TimeFormat: defaultTimeFormat,
		FullPath:   false,
	})
)

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

func Trace(msg string, data ...interface{}) {
	std.Trace(msg, data...)
}

func Debug(msg string, data ...interface{}) {
	std.Debug(msg, data...)
}

func Info(msg string, data ...interface{}) {
	std.Info(msg, data...)
}

func Warn(msg string, data ...interface{}) {
	std.Warn(msg, data...)
}

func Error(msg string, data ...interface{}) {
	std.Error(msg, data...)
}

func Fatal(msg string, data ...interface{}) {
	std.Fatal(msg, data...)
}

func Panic(msg string, data ...interface{}) {
	std.Panic(msg, data...)
}

func New(config Config) Interface {
	var (
		traceStr = White + "[TRACE]" + Reset
		debugStr = White + "[DEBUG]" + Reset
		infoStr  = Cyan + "[INFO]" + Reset
		warnStr  = Yellow + "[WARN]" + Reset
		errStr   = Red + "[ERROR]" + Reset
		fatalStr = Red + "[FATAL]" + Reset
		panicStr = Red + "[PANIC]" + Reset
	)

	return &logger{
		Config:   config,
		traceStr: traceStr,
		debugStr: debugStr,
		infoStr:  infoStr,
		warnStr:  warnStr,
		errStr:   errStr,
		fatalStr: fatalStr,
		panicStr: panicStr,
	}
}

func Sprintf(format string, h map[string]any) string {
	t := fasttemplate.New(format, "${", "}")
	return t.ExecuteString(h)
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

// Trace print messages
func (l *logger) Trace(msg string, data ...interface{}) {
	if l.LogLevel >= TraceLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.traceStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
	}
}

// Debug print messages
func (l *logger) Debug(msg string, data ...interface{}) {
	if l.LogLevel >= DebugLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.debugStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
	}
}

// Info print info
func (l *logger) Info(msg string, data ...interface{}) {
	if l.LogLevel >= InfoLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.infoStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
	}
}

// Warn print warn messages
func (l *logger) Warn(msg string, data ...interface{}) {
	if l.LogLevel >= WarnLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.warnStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
	}
}

// Error print error messages
func (l *logger) Error(msg string, data ...interface{}) {
	if l.LogLevel >= ErrorLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.errStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
	}
}

// Fatal print error messages
func (l *logger) Fatal(msg string, data ...interface{}) {
	if l.LogLevel >= ErrorLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.fatalStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
		os.Exit(1)
	}
}

func (l *logger) Panic(msg string, data ...interface{}) {
	if l.LogLevel >= ErrorLevel {
		path := l.getPath()
		fmt.Println(Sprintf(l.Format, map[string]any{
			"Level": l.panicStr,
			"Time":  l.t(),
			"Path":  path,
			"Msg":   fmt.Sprintf(msg, data...),
		}))
		panic(fmt.Sprintf(msg, data...))
	}
}

func (l *logger) t() string {
	return time.Now().Format(l.TimeFormat)
}

func (l *logger) getPath() string {
	path := FileWithLineNum()
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
