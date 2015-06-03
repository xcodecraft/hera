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
	var data string = ""
	for p_key, p_value := range params {
		if "1" == p_value {
			return c.Error("key has error", 1001, 400)
		} else {
			data += "key[" + p_key + "]=" + p_value + " " + hera.SERVER["HELLO"]
		}
	}
	return c.Success("access-data:" + data)
}

//curl 'localhost:8083/Hello/Set?fd=123'
func (this *HelloREST) Set(c *hera.Context) {
	fmt.Println("HelloREST::Set")
}

func init() {
	hera.NewRouter().AddRouter(&HelloREST{})
}
