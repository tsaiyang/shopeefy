package di

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

const redisConfigKey = "redis"

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr string `yaml:"addr"`
	}

	var cfg Config
	if err := viper.UnmarshalKey(redisConfigKey, &cfg); err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
	if client == nil {
		panic("init redis client failed")
	}

	return client
}
