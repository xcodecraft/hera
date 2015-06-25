package hera

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func CurlFunc(url string) (ret ReturnValue) {
	r, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("http.NewRequest: ", err.Error())
	}

	resp, err := http.DefaultClient.Do(r)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		fmt.Println("http.DefaultClient.Do: ", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode!=http.StatusOK: ", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		fmt.Println("ioutil.ReadAll: ", err.Error())
	}

	err = json.Unmarshal(data, &ret)
	if err != nil {
		fmt.Println("error:", err)
	}
	return ret
}
