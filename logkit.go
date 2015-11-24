package hera

import (
	"fmt"
	"log/syslog"
	"net/http"
	"time"
)

//使用方式
//logger =  hera.LogKit.Logger("目录")
//logger.debug('something message')
//logger.info('something message')

var LogKit *XLogKit = &XLogKit{} //用于管理XLogger
var Logger *XLogger = nil        //hera框架内部日志，存放在_hera目录下

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var LoggerLevel = map[string]int{
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
}

type XLogKit struct {
	prjName  string
	logLevel int
	loggers  map[string]*XLogger
}
type XLogger struct {
	logName   string
	logWriter *syslog.Writer
}

func (this *XLogKit) Init(prjName string, logLevel int) {
	this.prjName = prjName
	this.logLevel = logLevel
	this.loggers = make(map[string]*XLogger)
	if Logger == nil {
		Logger = newLogger("_hera")
	}
	this.loggers["_hera"] = Logger
	this.loggers["_all"] = newLogger("_all")     //存放所有的日志
	this.loggers["_speed"] = newLogger("_speed") //访问时间
}

func newLogger(logName string) *XLogger {
	logWriter := getWriter(LogKit.prjName, logName)
	logger := &XLogger{logName, logWriter}
	return logger
}
func getWriter(prjName, logName string) *syslog.Writer {
	writer, _ := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL6, prjName+"/"+logName)
	return writer
}

func (this *XLogKit) Logger(logName string) *XLogger {
	if logName == "" {
		panic("XLogger log name missing")
	}
	logger, ok := this.loggers[logName]
	if !ok {
		logger = newLogger(logName)
		this.loggers[logName] = logger
	}
	return logger
}

func (this *XLogger) Debug(str string) {
	if LogKit.logLevel <= LevelDebug {
		logger := LogKit.Logger("_all")
		logger.logWriter.Info(" [debug] " + str)
		this.logWriter.Info(" [debug] " + str)
	}
}
func (this *XLogger) Info(str string) {
	if LogKit.logLevel <= LevelInfo {
		logger := LogKit.Logger("_all")
		logger.logWriter.Info(" [info] " + str)
		this.logWriter.Info(" [info] " + str)
	}
}
func (this *XLogger) Warn(str string) {
	if LogKit.logLevel <= LevelWarn {
		logger := LogKit.Logger("_all")
		logger.logWriter.Info(" [warn] " + str)
		this.logWriter.Info(" [warn] " + str)
	}
}
func (this *XLogger) Error(str string) {
	if LogKit.logLevel <= LevelError {
		logger := LogKit.Logger("_all")
		logger.logWriter.Info(" [error] " + str)
		this.logWriter.Info(" [error] " + str)
	}
}

func (this *XLogger) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	next(rw, r)

	res := rw.(ResponseWriter)
	logger := LogKit.Logger("_speed")
	logger.Info(fmt.Sprintf("rest %s [%s]  %v %s usetime: %v ", r.Method, r.URL.Path, res.Status(), http.StatusText(res.Status()), time.Since(start)))

}
