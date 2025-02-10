package repository

import (
	"context"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type RedisRepositoryInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository() RedisRepositoryInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("Erro ao conectar ao Redis: %v", err))
	}

	return &RedisRepository{client: client}
}

func (r *RedisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	resp := r.client.Set(ctx, key, value, expiration)
	return resp.Err()
}

func (r *RedisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	return result == 1, err
}

func (r *RedisRepository) Close() error {
	return r.client.Close()
}
