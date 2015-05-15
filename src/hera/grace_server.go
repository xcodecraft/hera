package hera

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var defaultServer *GracefulServer

func SigHandler() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGTERM)
	for sig := range signals {
		if syscall.SIGTERM == sig {
			Close()
		} else if syscall.SIGHUP == sig {
			Restart()
		}

	}
}

func ListenAndServe(port int, handler http.Handler) (err error) {
	go SigHandler()

	if os.Getenv("_RESTART_") == "true" {
		defaultServer, err = NewServerFromFD(3, handler)
	} else {
		defaultServer, err = NewServer(port, handler)
	}
	if err != nil {
		return err
	}
	return defaultServer.Serve()
}

func Close() bool {
	return defaultServer.Close()
}

func Restart() bool {
	return defaultServer.Restart()
}
