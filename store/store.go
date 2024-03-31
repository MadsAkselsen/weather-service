package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store struct {
	client *redis.Client
}

func NewClient(opts *redis.Options) *redis.Client {
	return redis.NewClient(opts)
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return s.client.Set(ctx, key, value, expiration).Err()
}

func (s *Store) Get(key string) (string, error) {
	ctx := context.Background()
	return s.client.Get(ctx, key).Result()
}
