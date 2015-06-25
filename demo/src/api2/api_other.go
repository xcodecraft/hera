package api2

import (
	"fmt"
	hera  "github.com/xcodecraft/hera"
)

type ApiotherREST struct {
}

//curl 'localhost:8083/Apiother/Method?fd=123'
func (this *ApiotherREST) Method(c *hera.Context) error {
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

//curl 'localhost:8083/Apiother/Set?fd=123'
func (this *ApiotherREST) Set(c *hera.Context) {
	fmt.Println("HelloREST::Set")
}

func init() {
	hera.NewRouter().AddRouter(&ApiotherREST{})
}
