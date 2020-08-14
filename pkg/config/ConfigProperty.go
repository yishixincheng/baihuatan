package conf

import (
	"github.com/etcd-io/etcd/clientv3"
	"github.com/go-redis/redis"
)

var (
	Redis        RedisConf
	Etcd         EtcdConf
	MysqlConfig  MysqlConf
	TraceConfig  TraceConf
)
// EtcdConf 配置
type EtcdConf struct {
	EtcdConn              *clientv3.Client   //链接
	Host                  string             //主机
}

// TraceConf 配置
type TraceConf struct {
	Host  string
	Port  string
	Url   string
}

// MysqlConf 配置
type MysqlConf struct {
	Host  string
	Port  string
	User  string
	Pwd   string
	Db    string
}


// RedisConf 配置
type RedisConf struct {
	RedisConn             *redis.Client //链接
	Host                  string
	Password              string
	Db                    int
}

// AccessLimitConf 访问控制
type AccessLimitConf struct {
	IPSecAccessLimit   int   //IP每秒钟访问限制
	UserSecAccessLimit int   //用户每秒钟访问限制
	IPMinAccessLimit   int   //IP每分钟访问限制
	UserMinAccessLimit int   //用户每分钟访问限制
}