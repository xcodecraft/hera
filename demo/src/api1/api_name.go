package api1

import (
	"fmt"
	hera  "github.com/xcodecraft/hera"
)

type Api_nameREST struct {
}

//curl 'localhost:8083/Api_name/Get?fd=123'
func (this *Api_nameREST) Get(c *hera.Context) error {
	params := c.Params
	var data string = ""
	for p_key, p_value := range params {
		if "1" == p_value {
			return c.Error("key has error", 1001, 400)
		} else {
			data += "key[" + p_key + "]=" + p_value + " " + hera.SERVER["HELLO"]
		}
	}

	hera.Logger.Info("have sucess vistited Hello::Get() interface")
	return c.Success("access-data:" + data)
}

//curl 'localhost:8083/Api_name/Set?fd=123'
func (this *Api_nameREST) Set(c *hera.Context) {
	fmt.Println("HelloREST::Set")
}

func init() {
	hera.NewRouter().AddRouter(&Api_nameREST{})
}
