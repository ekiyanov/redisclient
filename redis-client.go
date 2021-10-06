package redisclient

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

func resolveHostFromEnv() string {
	host, port := os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")
	return resolveHost(host, port)
}

func resolveHost(host, port string) string {

	if host == "" {
		host = "redis:6379"
	}

	pair := strings.Split(host, ":")
	switch len(pair) {
	case 1:
		if port == "" {
			port = "6379"
		} else {
			_, err := strconv.ParseInt(port, 10, 32)
			if err != nil {
				log.Println("unable to parse port. Fallback to default 6379")
				port = "6379"
			}
		}
		host = host + ":" + port
	case 2:
	default:
		log.Println("Unable to parse REDIS_HOST. Have more than 1 semicolon. Fallback to redis:6379")

		return "redis:6379"
	}

	return host
}

func redisClientFromEnv() *redis.Client {
	password, db := os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_DB")

	host := resolveHostFromEnv()

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

var shared *redis.Client
var sharedMu sync.Mutex

func SharedRedisClient() *redis.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return SharedRedisClientCtx(ctx)
}

func SharedRedisClientCtx(ctx context.Context) *redis.Client {
	sharedMu.Lock()
	defer sharedMu.Unlock()

	if shared == nil {
		var err error
		shared, err = NewRedisClientCtx(ctx)
		if err != nil {
			log.Println("Unable to create redis client", err)
		}
	}

	return shared
}

func NewRedisClient() (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return NewRedisClientCtx(ctx)
}

func NewRedisClientCtx(ctx context.Context) (*redis.Client, error) {
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
		return nil, errors.New("Context timeout")
	case cc := <-c:
		return cc, nil
		break
	}

	log.Println("failed to get redis")
	return nil, errors.New("Unexpected error")
}
