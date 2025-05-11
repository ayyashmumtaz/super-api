package redisdb

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // Ganti kalau Redis pakai password
		DB:       0,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}
}

func IsTokenBlacklisted(token string) (bool, error) {
	val, err := Client.Get(Ctx, token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "blacklisted", nil
}

func BlacklistToken(token string, expiration time.Duration) error {
	return Client.Set(Ctx, token, "blacklisted", expiration).Err()
}
