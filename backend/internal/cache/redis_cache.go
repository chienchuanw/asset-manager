package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis 快取服務
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisCache 建立新的 Redis 快取服務
func NewRedisCache(addr string, password string, db int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()

	// 測試連線
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

// Get 取得快取資料
func (c *RedisCache) Get(key string, dest interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss: key not found")
	}
	if err != nil {
		return fmt.Errorf("failed to get cache: %w", err)
	}

	// 反序列化 JSON
	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

// Set 設定快取資料
func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	// 序列化為 JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	// 儲存到 Redis
	if err := c.client.Set(c.ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete 刪除快取資料
func (c *RedisCache) Delete(key string) error {
	if err := c.client.Del(c.ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// Exists 檢查快取是否存在
func (c *RedisCache) Exists(key string) (bool, error) {
	result, err := c.client.Exists(c.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	return result > 0, nil
}

// Close 關閉 Redis 連線
func (c *RedisCache) Close() error {
	return c.client.Close()
}

