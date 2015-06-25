package api1

import (
	"testing"

	hera "github.com/xcodecraft/hera"
)

func Test_Get(t *testing.T) {
	url := "http://localhost:8083/Api_name/Get?fd=1"
	ret := hera.CurlFunc(url)
	if ret.Errno != 0 {
		t.Error("get error")
	}
}
