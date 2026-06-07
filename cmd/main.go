package main

import (
	"log"
	"os"
	"time"

	"github.com/gdns/cmd/server"
	"github.com/gdns/internal/cache"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load env err: %s", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_URL"),
		Password:     "",
		DB:           0,
		Protocol:     2,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		DialTimeout:  3 * time.Second,
	})
	r := cache.CreateNewRedisClient(redisClient)
	server.Serve(*r)
}
