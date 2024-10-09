package cache

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "gin-books-api/configs"
    "github.com/go-redis/redis/v8"
)

// CacheExpiration defines the default expiration time for cache entries
const CacheExpiration = 10 * time.Second

// GetCachedData retrieves data from Redis cache.
// Returns true if data is successfully retrieved and unmarshaled into dest.
func GetCachedData(ctx context.Context, key string, dest interface{}) bool {
    cacheData, err := config.RedisClient.Get(ctx, key).Result()
    if err == redis.Nil {
        // Cache miss
        return false
    } else if err != nil {
        log.Printf("Redis GET error for key %s: %v", key, err)
        return false
    }

    if err := json.Unmarshal([]byte(cacheData), dest); err != nil {
        log.Printf("JSON Unmarshal error for key %s: %v", key, err)
        return false
    }
    return true
}

// SetCachedData stores data in Redis cache.
// Returns an error if the operation fails.
func SetCachedData(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
    jsonData, err := json.Marshal(data)
    if err != nil {
        return err
    }
    return config.RedisClient.Set(ctx, key, jsonData, expiration).Err()
}