package hera

import (
	"fmt"
	"log/syslog"
	"net/http"
	"time"
)

var Logger *XLogger = nil

// Log levels to control the logging output.
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

type XLogger struct {
	logName   string
	logLevel  int
	logWriter *syslog.Writer
}

func NewLogger(logName string, logLevel int) *XLogger {
	if Logger == nil {
		Logger = &XLogger{}
		Logger.Init(logName, LevelDebug)
	}
	return Logger
}

func (this *XLogger) Init(logName string, logLevel int) {
	this.logName = logName
	this.logLevel = logLevel
	this.logWriter = getWriter(logName)
}

func getWriter(logName string) *syslog.Writer {
	writer, _ := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, logName)
	return writer
}

func (this *XLogger) Logger() *syslog.Writer {
	if this.logName == "" {
		panic("XLogger log name missing")
	}
	if this.logWriter == nil {
		panic("XLogger log writer missing")
	}
	return this.logWriter
}

func (this *XLogger) Debug(str string) {
	if this.logLevel <= LevelDebug {
		this.Logger().Info(" [debug] " + str)
	}
}
func (this *XLogger) Info(str string) {
	if this.logLevel <= LevelInfo {
		this.Logger().Info(" [info] " + str)
	}
}
func (this *XLogger) Warn(str string) {
	if this.logLevel <= LevelWarn {
		this.Logger().Info(" [warn] " + str)
	}
}
func (this *XLogger) Error(str string) {
	if this.logLevel <= LevelError {
		this.Logger().Info(" [error] " + str)
	}
}

func (this *XLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	this.Info(fmt.Sprintf("Started %s %s", r.Method, r.URL.Path))

	next(rw, r)

	res := rw.(ResponseWriter)
	this.Info(fmt.Sprintf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start)))
}
