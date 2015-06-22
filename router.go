package hera

import (
	"net/http"
	"reflect"
	"strings"
)

var r *Router

type Router struct {
	autoRouter map[string]reflect.Type
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	structMethod := r.URL.Path
	structType, ok := this.autoRouter[structMethod]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("NO REST "))
		return
	}
	methods := strings.Split(structMethod, "/")
	methodName := methods[2]
	params := make(map[string]string)
	for formKey, formValue := range r.Form {
		params[formKey] = formValue[0]
	}
	method := reflect.New(structType).MethodByName(methodName)
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(&Context{Params: params, Req: r, Res: w})
	method.Call(in)
	return
}

func (this *Router) Start(addr string) {
	http.ListenAndServe(addr, this)
}

func (this *Router) AddRouter(i interface{}) {
	reflectVal := reflect.ValueOf(i)
	reflectType := reflectVal.Type()
	structType := reflect.Indirect(reflectVal).Type()
	structName := structType.Name()
	structPrefix := strings.TrimSuffix(structName, "REST")

	if structName == structPrefix {
		return
	}
	for i := 0; i < reflectType.NumMethod(); i++ {
		methodName := reflectType.Method(i).Name
		key := "/" + structPrefix + "/" + methodName
		this.autoRouter[key] = structType
	}

}

func NewRouter() *Router {
	if r == nil {
		r = &Router{
			autoRouter: make(map[string]reflect.Type),
		}
	}
	return r
}
