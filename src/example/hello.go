package example

import (
	"fmt"
	"hera"
)

type HelloREST struct {
}

//curl 'localhost:8083/Hello/Get?fd=123'
func (this *HelloREST) Get(c *hera.Context) {
	params := c.Params
	for p_key, p_value := range params {
		c.Json("param key:" + p_key + " value: " + p_value)
	}
	c.Json("   access ok")
}

//curl 'localhost:8083/Hello/Set?fd=123'
func (this *HelloREST) Set(c *hera.Context) {
	fmt.Println("HelloREST::Set")
}

func init() {
	hera.NewRouter().AddRouter(&HelloREST{})
}
