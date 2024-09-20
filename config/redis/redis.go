package redis

import (
	"context"
	"go_assignment/logger"
	"go_assignment/utils"

	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var Client *redis.Client

func ConnectToRedis() (err error) {
	ctx := context.Background()
	Client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString(utils.RedisAddress),
		Password: utils.EmptySpace,
		DB:       0,
	})
	pong, err := Client.Ping(ctx).Result()
	if err != nil {
		logger.E(err)
	}
	logger.I(pong)
	return
}
