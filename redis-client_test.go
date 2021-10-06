package redisclient

import (
	"context"
	"testing"
)

func TestResolver(t *testing.T) {

	host := resolveHost("", "")
	if host != "redis:6379" {
		t.Fatal("host should be redis:6379", host)
	}

	host = resolveHost("", "80")
	if host != "redis:6379" {
		t.Fatal("host(,80) should be redis:6379", host)
	}

	host = resolveHost("redis", "")
	if host != "redis:6379" {
		t.Fatal("host(redis) should be redis:6379", host)
	}

	host = resolveHost("redis:6355", "")
	if host != "redis:6355" {
		t.Fatal("host(redis:6355) should be redis:6355")
	}

	host = resolveHost("redis:6355", "80")
	if host != "redis:6355" {
		t.Fatal("host(redis:6355,80) should be redis:6355")
	}

	host = resolveHost("redis", "6344")
	if host != "redis:6344" {
		t.Fatal("host(redis,6344) should be redis:6344")
	}

	host = resolveHost("redis://redis:6379", "6344")
	if host != "redis:6379" {
		t.Fatal("host(redis,6344) should be redis:6344")
	}

}

func TestClients(t *testing.T) {
	client := SharedRedisClient()

	client2 := SharedRedisClient()

	if client != client2 {
		t.Fatal("Shared should return same instance")
	}

	client3, err := NewRedisClientCtx(context.Background())
	if err != nil {
		t.Fatal("Should create instance with no issues")
	}

	if client3 == client {
		t.Fatal("Should return different instance")
	}

}
