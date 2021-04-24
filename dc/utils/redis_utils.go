package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strings"
)

type RedisProxy struct {
	ns     string
	client *redis.Client
}

func NewRedisProxy(ns string) *RedisProxy {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	rp := RedisProxy{
		ns:     ns,
		client: client,
	}
	return &rp
}

func (r *RedisProxy) genKey(key string) string {
	return fmt.Sprintf("%v:%v", r.ns, key)
}

func (r *RedisProxy) SetString(key string, data string) error {
	key = r.genKey(key)
	ctx := context.TODO()
	err := r.client.Set(ctx, key, data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisProxy) GetString(key string) (string, error) {
	key = r.genKey(key)
	ctx := context.TODO()
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *RedisProxy) Set(key string, val interface{}) error {
	key = r.genKey(key)
	bBytes, err := json.Marshal(val)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	err = r.client.Set(ctx, key, string(bBytes), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisProxy) Get(key string, result interface{}) error {
	if !strings.Contains(key, r.ns) {
		key = r.genKey(key)
	}
	ctx := context.TODO()
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), result)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisProxy) GetKeys(wildcard string) []string {
	if !strings.Contains(wildcard, r.ns) {
		wildcard = fmt.Sprintf("%v:*%v", r.ns, wildcard)
	}
	ctx := context.TODO()
	ssc := r.client.Keys(ctx, wildcard)
	return ssc.Val()
}

func (r *RedisProxy) DelKeys(wildcard string) int {
	keys := r.GetKeys(wildcard)
	ctx := context.TODO()
	r.client.Del(ctx, keys...)
	return len(keys)
}

func (r *RedisProxy) SAdd(key string, members ...string) {
	key = r.genKey(key)
	r.client.SAdd(context.TODO(), key, members)
}

func (r *RedisProxy) SMembers(key string) []string {
	key = r.genKey(key)
	return r.client.SMembers(context.TODO(), key).Val()
}
