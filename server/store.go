package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	Save(ctx context.Context, key string, data string, expires time.Duration) error
	Load(ctx context.Context, key string) (string, error)
}

// implementation with maps

type MapStore struct {
	m map[string]string
}

func NewMapStore() *MapStore {
	return &MapStore{
		m: make(map[string]string),
	}
}

func (x *MapStore) Save(ctx context.Context, key string, data string, expires time.Duration) error {
	time.AfterFunc(expires, func() {
		if _, ok := x.m[key]; ok {
			delete(x.m, key)
		}
	})
	x.m[key] = data
	return nil
}

func (x *MapStore) Load(ctx context.Context, key string) (string, error) {
	val, ok := x.m[key]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return val, nil
}

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(host, port string) *RedisStore {
	//	una pinche conexion
	//	vamos a considerar que ya creaste la conexion y tienes un cliente, se lo asignas a tu objeto
	return &RedisStore{
		client: nil,
	}
}

func (x *RedisStore) Save(ctx context.Context, key string, data string, expires time.Duration) error {
	cmd := x.client.Set(ctx, key, data, expires)
	return cmd.Err()
}

func (x *RedisStore) Load(ctx context.Context, key string) (string, error) {
	cmd := x.client.Get(ctx, key)
	//data,err:=cmd.Bytes()
	//if err!=nil{
	//	log.Print("redis error: ", err.Error())
	//	return nil, fmt.Errorf("not found")
	//}
	//return string(data),nil
	data, err := cmd.Result()
	if err != nil {
		log.Print("redis error: ", err.Error())
		return "", fmt.Errorf("not found")
	}
	return data, nil
}
