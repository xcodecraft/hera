package hera

import (
// "fmt"
// "log"
// "net/http"
// "reflect"
// "runtime"
)

type XLogicException interface {
	GetCode() int
	GetMessage() string
	// Error() string
	Throw()
}
type XRunningtimeException interface {
	GetCode() int
	GetMessage() string
	// Error() string
	Throw()
}

// func (e *LogicException) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			if reflect.TypeOf(err) == "LogicException" {
// 				log.Fatal(err)
// 			}
// 		}
// 	}()
//
// 	next(rw, r)
// }
