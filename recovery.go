package hera

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Recovery struct {
	PrintStack bool
	StackAll   bool
	StackSize  int
}

func NewRecovery() *Recovery {
	return &Recovery{
		//上线关掉此开关
		PrintStack: false,
		StackAll:   false,
		StackSize:  1024 * 80,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := make([]byte, rec.StackSize)

			stackStr := string(stack[:runtime.Stack(stack, rec.StackAll)])

			f := "hera  catch panic: %s , %s"

			if rec.PrintStack {
				fmt.Fprintf(rw, f, err, stackStr)
			}
			stackStr = strings.Replace(stackStr, "\n", "\t", -1)
			Logger.Error(fmt.Sprintf(f, err, stackStr))
		}
	}()

	next(rw, r)
}
