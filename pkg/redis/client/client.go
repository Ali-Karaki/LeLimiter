package redis

import (
	"context"
	"rate-limiter/pkg/redis/lua"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis() (*redis.Client, error) {
	ctx := context.Background()
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := Client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	Client.FlushAll(context.Background()) 
	lua.LoadShaScript(Client)
	return Client, nil
}

func GetClient() *redis.Client {
	return Client
}
