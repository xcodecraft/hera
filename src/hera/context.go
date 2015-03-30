package hera

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Params map[string]string
	Req    *http.Request       "request"
	Res    http.ResponseWriter "response"
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

func (c *Context) Json(data interface{}) error {
	c.Res.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(data)
	if err == nil {
		c.Res.Write(j)
	}
	return err
}

func (c *Context) Redirect(url string) {
	c.Res.Header().Set("Location", url)
	c.Res.Write([]byte(""))
}
