package services

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type RedisService struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *RedisService {
	return &RedisService{
		client: client,
	}
}
func (svc *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	return svc.client.Set(key, value, ttl*time.Second).Err()
}
func (svc *RedisService) Get(key string) (string, error) {
	fmt.Println("Getting key:", key)

	return svc.client.Get(key).Result()
}

func (svc *RedisService) GetKey(key string) (int, error) {
	str, err := svc.client.Get(key).Result()
	if err != nil {
		return 0, nil
	}

	return strconv.Atoi(str)
}
func (svc *RedisService) GetInt(key string) (int, error) {
	str, err := svc.client.Get(key).Result()
	if err != nil {
		return 0, nil
	}

	return strconv.Atoi(str)
}
func (svc *RedisService) SetStruct(key string, value interface{}, ttl time.Duration) error {
	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return svc.client.Set(key, string(serializedValue), ttl*time.Second).Err()
}

func (svc *RedisService) GetStruct(key string, outputStruct interface{}) error {
	value, err := svc.client.Get(key).Result()
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(value), &outputStruct); err != nil {
		return err
	}
	return nil
}
func (svc *RedisService) Delete(keys ...string) error {
	return svc.client.Del(keys...).Err()
}
