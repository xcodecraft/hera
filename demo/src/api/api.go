package api

import hera "github.com/xcodecraft/hera"

type ApiREST struct {
}

func (this *ApiREST) GetHello(c *hera.Context) error {
	return c.Success("Hello")
}

func (this *ApiREST) GetWorld(c *hera.Context) error {
	return c.Error("GetWorld has error", 1001, 400)
}

func init() {
	hera.NewRouter().AddRouter(&ApiREST{})
}
