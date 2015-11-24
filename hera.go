package hera

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
)

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, next)
}

type middleware struct {
	handler Handler
	next    *middleware
}

func (m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}

func Wrap(handler http.Handler) Handler {
	return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handler.ServeHTTP(rw, r)
		next(rw, r)
	})
}

type hera struct {
	middleware middleware
	handlers   []Handler
}

func New(handlers ...Handler) *hera {
	return &hera{
		handlers:   handlers,
		middleware: build(handlers)}
}

func classic() *hera {
	return New(NewRecovery(), Logger)
}

//func Run(confPath string) {
func Run(port string) {
	//InitEnv(confPath)
	startServ(port)
}

//加载环境变量和日志级别
func InitEnv(confPath string) {
	if _, err := os.Stat(confPath); err != nil {
		panic("hera conf file is not exist")
	}
	NewEnv(confPath)
	logLevel := ENV["LOGLEVEL"]
	if _, ok := LoggerLevel[logLevel]; !ok {
		panic("hera loglevel is wrong")
	}
	fmt.Println("prjname :" + ENV["PRJ_NAME"] + "  log level : " + logLevel)
	LogKit.Init(ENV["PRJ_NAME"], LoggerLevel[logLevel])
}

func startServ(port string) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	n := classic()
	n.Run(port)
}

func (n *hera) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	n.middleware.ServeHTTP(NewResponseWriter(rw), r)
}
func (n *hera) Use(handler Handler) {
	n.handlers = append(n.handlers, handler)
	n.middleware = build(n.handlers)
}

func (n *hera) UseFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	n.Use(HandlerFunc(handlerFunc))
}

func (n *hera) UseHandler(handler http.Handler) {
	n.Use(Wrap(handler))
}

func (n *hera) UseHandlerFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request)) {
	n.UseHandler(http.HandlerFunc(handlerFunc))
}

func (n *hera) Run(addr string) {
	n.UseHandler(NewRouter())
	Logger.Info(fmt.Sprintf("listening on %v", addr))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", addr), n); err != nil {
		Logger.Error(fmt.Sprintf("server start fail: %s", err))
		panic(fmt.Sprintf("server start fail: %s", err))
	}
}

func (n *hera) Handlers() []Handler {
	return n.handlers
}

func build(handlers []Handler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
		&middleware{},
	}
}
