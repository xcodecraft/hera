package hera

import (
	"fmt"
	// "log"
	"net/http"
	"reflect"
	"strings"
)

var r *Router

type Router struct {
	autoRouter map[string]reflect.Type
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//对于现实LogicException的类做统一的处理
	defer func(w http.ResponseWriter) {
		if err := recover(); err != nil {
			e, ok := err.(XLogicException)
			if ok {
				rv := NewReturnValue(e.GetCode(), e.GetMessage(), nil)
				w.Header().Set("ContentType", "application/json")
				w.WriteHeader(200)
				w.Write(rv.Json())
				return
			}
			panic(err)
		}
	}(w)
	//获取url并判断该url在注册map中
	r.ParseForm()
	structMethod := r.URL.Path
	structMethod = strings.ToLower(structMethod)
	structType, ok := this.autoRouter[structMethod]

	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("NO REST "))
		return
	}
	//获取反射的方法
	methods := strings.Split(structMethod, "/")
	methodName := ucfirst(methods[len(methods)-1])
	method := reflect.New(structType).MethodByName(methodName)
	//构造反射的参数
	params := make(map[string]string)
	for formKey, formValue := range r.Form {
		formKey = strings.ToLower(formKey)
		Value := strings.ToLower(formValue[0])
		params[formKey] = Htmlspecialchars(Value)
	}
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(&Context{Params: params, Req: r, Res: w})
	//call反射的方法
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

	//将类名Api_name中的‘_’转换为url中的'/'
	arrayStructName := strings.Split(structName, "_")
	structPrefix := strings.Join(arrayStructName, "/")

	if structName == structPrefix {
		return
	}
	for i := 0; i < reflectType.NumMethod(); i++ {
		//这里解决一个bug，当methodName首字母为小写是，也会
		//被映射到，导致出异常
		methodName := reflectType.Method(i).Name
		if !isUcFirst(methodName) {
			continue
		}
		key := "/" + structPrefix + "/" + methodName
		key = strings.ToLower(key)
		this.autoRouter[key] = structType
		fmt.Println("register router : " + key)
	}
}

//将首字母转换未大写
func ucfirst(s string) string {
	r := []rune(s)
	r[0] = r[0] - 32
	return string(r)
}

//判断首字母是否大学
func isUcFirst(s string) bool {
	r := []rune(s)
	if r[0] >= 65 && r[0] <= 90 {
		return true
	}
	return false
}

//单例
func NewRouter() *Router {
	if r == nil {
		r = &Router{
			autoRouter: make(map[string]reflect.Type),
		}
	}
	return r
}
