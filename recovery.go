package hera

import (
	"fmt"
	"net/http"
	"runtime"
)

type Recovery struct {
	PrintStack bool
	StackAll   bool
	StackSize  int
}

func NewRecovery() *Recovery {
	return &Recovery{
		PrintStack: true,
		StackAll:   false,
		StackSize:  1024 * 8,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)
			stack = stack[:runtime.Stack(stack, rec.StackAll)]

			f := "PANIC: %s\n%s"
			Logger.Info(fmt.Sprintf(f, err, stack))

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stack)
			}
		}
	}()

	next(rw, r)
}
