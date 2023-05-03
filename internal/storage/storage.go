package storage

import (
	"context"
	"github.com/Volkacid/mediaBalancer/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	Client *redis.Client
}

func NewRedisStorage(ctx context.Context) (*RedisStorage, error) {
	conf := config.GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddress,
		Password: conf.RedisPassword,
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &RedisStorage{Client: rdb}, nil
}
