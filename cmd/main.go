package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gdns/cmd/server"
	"github.com/gdns/internal/blocklist"
	"github.com/gdns/internal/cache"
	"github.com/redis/go-redis/v9"
)

func main() {
	sqliteClient, err := blocklist.CreateNewDbConn(os.Getenv("SQLITE_URL"))
	if err != nil {
		log.Fatalf("failed to connect to sqlite prev: %s", err)
	}

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err = sqliteClient.Migrate(ctxWithTimeout)
	if err != nil {
		log.Fatalf("failed to create db prev: %s", err)
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
	srv := server.NewServer(*r, *sqliteClient)
	if err := srv.Start(); err != nil {
		log.Fatalf("failed to start server: %s", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	srv.Stop()
}
