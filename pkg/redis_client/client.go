package redis_client

import (
	"context"
	"fmt"
	"os"

	redis "github.com/go-redis/redis/v9"
)

type RedisClient struct {
	Svc *redis.Client
}

func InitClient() (redisClient *RedisClient) {
	redisClient.Svc = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	return
}

func (redisClient *RedisClient) GetValue(ctx context.Context, key string) (value string, err error) {
	value, getErr := redisClient.Svc.Get(ctx, key).Result()
	if getErr != nil {
		err = fmt.Errorf("key %v not found; %v", key, getErr)
	}

	return
}

func (redisClient *RedisClient) Store(ctx context.Context, key string, value string) (err error) {
	if storeErr := redisClient.Svc.Set(ctx, key, value, 0).Err(); storeErr != nil {
		err = fmt.Errorf("unable to store kv pair; [error: %v]", storeErr)
	}
	return
}

func (redisClient *RedisClient) UpdateValue(ctx context.Context, key string, newValue string) (err error) {
	if _, getErr := redisClient.Svc.Get(ctx, key).Result(); getErr != nil {
		err = fmt.Errorf("unable to update value; key %v not found; [error: %v]", key, getErr)
		return err
	}

	if storeErr := redisClient.Svc.Set(ctx, key, newValue, 0).Err(); storeErr != nil {
		err = fmt.Errorf("unable to update value; failed to set vaule; [error: %v]", storeErr)
	}

	return
}

func (redisClient *RedisClient) Delete(ctx context.Context, keys ...string) (err error) {
	if dErr := redisClient.Svc.Del(ctx, keys...).Err(); dErr != nil {
		err = fmt.Errorf("unable to delete keys %v; [error: %v]", keys, dErr)
	}

	return
}
