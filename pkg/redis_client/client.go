package redis_client

import (
	"context"
	"fmt"
	"os"

	redis "github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
)

type RedisClient struct {
	svc *redis.Client
}

type GetValue struct {
	Key string
}

type KVPair struct {
	Key   string
	Value string
}

type DeleteKeys struct {
	Keys []string
}

func InitClient() *RedisClient {
	var redisClient RedisClient
	redisClient.svc = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})

	//test redis connection
	pong, err := redisClient.svc.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to get response from redis server; [error: %v]", err))
	} else {
		log.Infof("Redis PING Response: %v", pong)
	}

	return &redisClient
}

func (redisClient *RedisClient) GetValue(ctx context.Context, key string) (value string, err error) {
	value, getErr := redisClient.svc.Get(ctx, key).Result()
	if getErr != nil {
		err = fmt.Errorf("key %v not found; %v", key, getErr)
	}

	return
}

func (redisClient *RedisClient) Store(ctx context.Context, kv KVPair) (err error) {
	if storeErr := redisClient.svc.Set(ctx, kv.Key, kv.Value, 0).Err(); storeErr != nil {
		err = fmt.Errorf("unable to store kv pair; [error: %v]", storeErr)
	}
	return
}

func (redisClient *RedisClient) UpdateValue(ctx context.Context, kv KVPair) (err error) {
	if _, getErr := redisClient.svc.Get(ctx, kv.Key).Result(); getErr != nil {
		err = fmt.Errorf("unable to update value; key %v not found; [error: %v]", kv.Key, getErr)
		return err
	}

	if storeErr := redisClient.svc.Set(ctx, kv.Key, kv.Value, 0).Err(); storeErr != nil {
		err = fmt.Errorf("unable to update value; failed to set vaule; [error: %v]", storeErr)
	}

	return
}

func (redisClient *RedisClient) Delete(ctx context.Context, keys DeleteKeys) (err error) {
	if dErr := redisClient.svc.Del(ctx, keys.Keys...).Err(); dErr != nil {
		err = fmt.Errorf("unable to delete keys %v; [error: %v]", keys.Keys, dErr)
	}

	return
}
