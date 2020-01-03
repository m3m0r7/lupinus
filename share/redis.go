package share

import (
	"github.com/go-redis/redis/v7"
	"os"
)

type Redis struct {
	Client *redis.Client
}

func InitRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "",
		DB:       0,
	})

	return &Redis{
		Client: client,
	}
}
