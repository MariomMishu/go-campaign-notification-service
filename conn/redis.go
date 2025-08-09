package conn

import (
	"ems/config"

	"github.com/labstack/gommon/log"

	"github.com/go-redis/redis"
)

var client *redis.Client

func ConnectRedis() {
	conf := config.Redis()

	log.Info("connecting to redis at ", conf.Host, ":", conf.Port, "...")

	client = redis.NewClient(&redis.Options{
		Addr:     conf.Host + ":" + conf.Port,
		Password: conf.Pass,
		DB:       conf.Db,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Info("failed to connect redis: ", err)
		panic(err)
	}

	log.Info("redis connection successful...")
}

func Redis() *redis.Client {
	return client
}
