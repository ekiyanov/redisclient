package redisclient

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

func redisClientFromEnv() *redis.Client {
	host, password, db := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DB")

	if host == "" {
		host = "redis:6379"
	}

	dbInt := 0
	if db != "" {
		dbInt, _ = strconv.Atoi(db)
	}

	log.Println("Creating Redis Client to", host, "db", dbInt)

	return redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbInt,
	})

}

func NewRedisClient() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return NewRedisClientCtx(ctx)
}

func NewRedisClientCtx(ctx context.Context) *redis.Client {
	c := make(chan *redis.Client, 1)
	go func() {
		for {
			cli := redisClientFromEnv()
			if cli != nil {
				err := cli.Ping(ctx).Err()
				if err == nil {
					c <- cli
					return
				} else {
					log.Println("unable to reach redis. Retry in 500ms", err)
					time.Sleep(500 * time.Millisecond)
				}
			} else {
				log.Println("redis client is not exist. probably bad env vars")
				time.Sleep(2500 * time.Millisecond)
			}
		}
	}()

	select {
	case <-ctx.Done():
		return nil
	case cc := <-c:
		return cc
		break
	}

	log.Println("failed to get redis")
	return nil
}
