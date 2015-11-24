package mysql

import (
	"database/sql"
	"exception"
	_ "github.com/go-sql-driver/mysql"
)

var SQLExecuter = nil

type Executer struct {
	host string
	user string
	pwd  string
	name string
	port string
}

func New(host, user, pwd, name string) *Executer {
	return &Executer{host, user, pwd, name, port}
}

// func (exec *Executer) conn() {
// 	db, err := sql.Open("mysql", user+":"+pwd+"@tcp("+host+":"+port+")/"+name+"?charset=utf8")
// 	if err != nil {
// 		gg
// 	}
// }
