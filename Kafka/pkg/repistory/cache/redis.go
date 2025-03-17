package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Address  string
	Port     int
	Password string
	DB       int
}

type RedisResponse struct {
	Message string
	Success bool
}

func RedisConnect(cfg RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis bağlantısı başarısız: %v", err)
	}
	log.Println("✅ Redis bağlantısı başarılı!")
	return rdb
}
