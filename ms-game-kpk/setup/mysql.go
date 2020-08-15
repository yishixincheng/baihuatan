package setup

import (
	conf "baihuatan/pkg/config"
	"baihuatan/pkg/mysql"
)

// InitMysql - 链接mysql
func InitMysql() {
	mysql.InitMysql(conf.MysqlConfig.Host,
		conf.MysqlConfig.Port,
		conf.MysqlConfig.User,
		conf.MysqlConfig.Pwd,
		conf.MysqlConfig.Db,
	)
}




