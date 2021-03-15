package main

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Save(key string, data string, expires time.Duration) error
	Load(key string) (string, error)
	Delete(key string) error
}

type MapSyncStore struct {
	StoreData sync.Map
}

func (x *MapSyncStore) Save(key string, data string, expires time.Duration) error {
	time.AfterFunc(expires, func() {
		x.StoreData.Delete(key)
	})
	x.StoreData.Store(key, data)
	return nil
}

func (x *MapSyncStore) Load(key string) (string, error) {
	value, ok := x.StoreData.Load(key)
	if !ok {
		return "", nil
	}
	return value.(string), nil
}

func (x *MapSyncStore) Delete(key string) error {
	x.StoreData.Delete(key)
	return nil
}

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
		return "", nil
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
