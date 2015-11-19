package hera

//这个是一个一般得工具类，提供框架一般得工具服务

import (
	"strings"
)

//做一般得sql过滤喝xss攻击
func Htmlspecialchars(str string) string {

	rst := strings.Replace(str, "&amp;", "&", -1)
	rst = strings.Replace(str, "&quot;", "\"", -1)
	rst = strings.Replace(str, "&lt;", "<", -1)
	rst = strings.Replace(str, "&gt;", ">", -1)

	rst = strings.Replace(str, "&", "&amp;", -1)
	rst = strings.Replace(str, "\"", "&quot;", -1)
	rst = strings.Replace(str, "<", "&lt;", -1)
	rst = strings.Replace(str, ">", "&gt;", -1)

	rst = strings.Replace(str, "'", "\\'", -1)
	rst = strings.Replace(str, "\\", "\\\\", -1)
	rst = strings.Replace(str, "NUL", "\\NUL", -1)
	rst = strings.Replace(str, "NULL", "\\NULL", -1)

	return rst
}
