package setup

import (
	"github.com/go-redis/redis"
	conf "baihuatan/pkg/config"
	"log"
)

// InitRedis 初始化Redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:       conf.Redis.Host,
		Password:   conf.Redis.Password, 
		DB:         conf.Redis.Db,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. Error : %v", err)
	}
	conf.Redis.RedisConn = client
	//			data, err := conn.BRPop(time.Second, conf.Redis.Proxy2layerQueueName).Result()
}