package hera

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Context struct {
	Params map[string]string
	Req    *http.Request       "request"
	Res    http.ResponseWriter "response"
}

type HeadStatus int

const (
	STATUS_OK                    HeadStatus = 200
	STATUS_MOVE                  HeadStatus = 302
	STATUS_REQUEST_ERROR         HeadStatus = 400
	STATUS_UNAUTHORIZED          HeadStatus = 401
	STATUS_NOT_FOUND             HeadStatus = 404
	STATUS_INTERNAL_SERVER_ERROR HeadStatus = 500
	STATUS_NOT_IMPLEMENTD        HeadStatus = 501
)

var headStatusMap = map[HeadStatus]string{
	STATUS_OK:                    "Ok",
	STATUS_MOVE:                  "Move temporarily",
	STATUS_REQUEST_ERROR:         "Request Error",
	STATUS_UNAUTHORIZED:          "Unauthorized",
	STATUS_NOT_FOUND:             "Not Found",
	STATUS_INTERNAL_SERVER_ERROR: "Internal Server Error",
	STATUS_NOT_IMPLEMENTD:        "Not Implementd",
}

type ReturnValue struct {
	Errno  int         `json:"errno"`
	Errmsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func NewReturnValue(errno int, errmsg string, data interface{}) *ReturnValue {
	return &ReturnValue{errno, errmsg, data}
}

func (rv *ReturnValue) Json() []byte {
	ret, _ := json.Marshal(*rv)
	return ret
}

func (c *Context) Success(data interface{}) error {
	rv := NewReturnValue(0, " ", data)
	c.SetHeader("ContentType", "application/json")
	c.Res.WriteHeader(int(STATUS_OK))
	c.Res.Write(rv.Json())
	return nil
}

func (c *Context) Error(errmsg string, errno int, statusCode HeadStatus) error {
	rv := NewReturnValue(errno, errmsg, nil)
	_, ok := headStatusMap[statusCode]
	if !ok {
		panic(errors.New("status code is invalid"))
	}
	c.SetHeader("ContentType", "application/json")
	c.Res.WriteHeader(int(statusCode))
	c.Res.Write(rv.Json())
	return nil
}

func (c *Context) SetHeader(key, value string) {
	c.Res.Header().Set(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}

func (c *Context) GetCookie(key string) string {
	cookie, err := c.Req.Cookie(key)
	if err != nil {
		return ""
	}
	return cookie.String()
}

func (c *Context) Redirect(url string) {
	c.Res.Header().Set("Location", url)
	c.Res.Write([]byte(""))
}
