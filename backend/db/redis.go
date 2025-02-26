package db

import "github.com/redis/go-redis/v9"

type RedisProvider struct {
	Client *redis.Client
}

var redisProvider *RedisProvider

func NewRedisProvider(uri, password string) *RedisProvider {
	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: password,
		DB:       0,
	})

	setRedisProvider(&RedisProvider{Client: client})
	return redisProvider
}

func setRedisProvider(provider *RedisProvider) {
	redisProvider = provider
}

func GetRedisProvider() *RedisProvider {
	return redisProvider
}

func (rp *RedisProvider) Close() error {
	return rp.Client.Close()
}
