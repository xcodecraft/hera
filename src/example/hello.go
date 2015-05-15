package example

import (
	"fmt"
	"hera"
)

type HelloREST struct {
}

//curl 'localhost:8083/Hello/Get?fd=123'
func (this *HelloREST) Get(c *hera.Context) error {
	params := c.Params
	for p_key, p_value := range params {
		if "1" == p_value {
			return c.Success("access-func")
		} else {
			return c.Error("param key:"+p_key+" value: "+p_value+"access-error", 1001, 400)
		}
	}
	return c.Success("access-func")
}

//curl 'localhost:8083/Hello/Set?fd=123'
func (this *HelloREST) Set(c *hera.Context) {
	fmt.Println("HelloREST::Set")
}

func init() {
	hera.NewRouter().AddRouter(&HelloREST{})
}
