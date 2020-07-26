package mysql

import (
	"fmt"
	// _ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

var engin *gorose.Engin
var err error

// InitMysql 初始化
func InitMysql(host, port, user, pwd, db string) {
	fmt.Printf(user)
	fmt.Printf(db)

	DbConfig := gorose.Config{
		Driver: "mysql",
		Dsn: user + ":" + pwd + "@tcp(" + host + ":" + port +")/" + db + "?charset=utf8&parseTime=true",
		Prefix: "",
		SetMaxOpenConns: 300,
		SetMaxIdleConns: 10,
	}

	engin, err = gorose.Open(&DbConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// DB IOrm
func DB() gorose.IOrm {
	return engin.NewOrm()
}