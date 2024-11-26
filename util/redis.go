package util

import (
	"context"
	"fmt"

	"magmar/config"

	"github.com/go-redis/redis/v8"
)

// RedisConnect ...
func RedisConnect(magmar *config.ViperConfig, zlog *Logger) (redisDB *redis.Client, err error) {
	host := fmt.Sprintf("%s:%d", magmar.GetString("redis.host"), magmar.GetInt("redis.port"))
	zlog.Infow("InitRedis", "redis_host", host)
	redisDB = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: "",
	})
	if _, err := redisDB.Ping(context.Background()).Result(); err != nil {
		return nil, err
	}
	return redisDB, nil
}
