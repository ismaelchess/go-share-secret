package stores

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string, ctx context.Context) *RedisStore {
	dbr := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	return &RedisStore{
		client: dbr,
		ctx:    ctx,
	}
}

func (x *RedisStore) Save(key string, data string, expires time.Duration) error {
	err := x.client.Set(x.ctx, key, data, expires).Err()
	if err != nil {
		return err
	}
	return nil
}

func (x *RedisStore) Load(key string) (string, error) {
	value, err := x.client.Get(x.ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	if err != nil {
		return "", err
	}
	return value, nil
}

func (x *RedisStore) Delete(key string) error {
	err := x.client.Del(x.ctx, key).Err()
	if err != nil {
		return err
	}
	return nil
}
