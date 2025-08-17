package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool
}

func NewRedisCache(cfg *CacheConfig) (*RedisCache, error) {
	if !cfg.Enabled {
		log.Println("Redis cache is disabled")
		return &RedisCache{ctx: context.Background()}, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()

	// Test the connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis cache...")
		return &RedisCache{ctx: ctx}, nil
	}

	log.Printf("Connected to Redis at %s:%s", cfg.Host, cfg.Port)
	return &RedisCache{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return nil // Cache disabled
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(r.ctx, key, jsonValue, expiration).Err()
}

func (r *RedisCache) Get(key string, dest interface{}) error {
	if r.client == nil {
		return fmt.Errorf("cache miss: Redis not available")
	}

	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: key not found")
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisCache) Delete(key string) error {
	if r.client == nil {
		return nil // Cache disabled
	}

	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) DeletePattern(pattern string) error {
	if r.client == nil {
		return nil // Cache disabled
	}

	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(r.ctx, keys...).Err()
	}

	return nil
}

func (r *RedisCache) Exists(key string) bool {
	if r.client == nil {
		return false
	}

	count, err := r.client.Exists(r.ctx, key).Result()
	return err == nil && count > 0
}

func (r *RedisCache) SetTTL(key string, expiration time.Duration) error {
	if r.client == nil {
		return nil
	}

	return r.client.Expire(r.ctx, key, expiration).Err()
}

func (r *RedisCache) GetTTL(key string) (time.Duration, error) {
	if r.client == nil {
		return 0, fmt.Errorf("Redis not available")
	}

	return r.client.TTL(r.ctx, key).Result()
}

func (r *RedisCache) Close() error {
	if r.client == nil {
		return nil
	}

	return r.client.Close()
}