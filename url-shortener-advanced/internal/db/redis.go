package db

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	Rdb *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR") // e.g. "localhost:6379" or "redis:6379"
	if addr == "" {
		addr = "localhost:6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := Rdb.Ping(Ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
}
