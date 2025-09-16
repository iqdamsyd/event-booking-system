package redis

import (
	"event-booking-system/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	redisUrl := config.GetConfig().RedisUrl
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
