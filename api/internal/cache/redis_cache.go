package cache

import (
	"api/internal/config"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	wrapper *RedisWrapper
}

func NewRedisCache(cfg *config.Config) (*RedisCache, error) {
	wrapper, err := NewRedisWrapper(&cfg.Redis)
	if err != nil {
		return nil, err
	}

	return &RedisCache{
		wrapper: wrapper,
	}, nil
}

// Set сохраняет значение в кэш
func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.wrapper.Set(r.wrapper.ctx, key, data, expiration).Err()
}

// Get получает значение из кэша
func (r *RedisCache) Get(key string, dest interface{}) error {
	data, err := r.wrapper.Get(r.wrapper.ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

// Delete удаляет ключ из кэша
func (r *RedisCache) Delete(key string) error {
	return r.wrapper.Del(r.wrapper.ctx, key).Err()
}

// DeleteByPattern удаляет ключи по паттерну
func (r *RedisCache) DeleteByPattern(pattern string) error {
	if r.wrapper.clientType == "cluster" {
		log.Printf("Warning: DeleteByPattern in cluster mode may not work correctly for distributed keys")
	}

	keys, err := r.wrapper.Keys(r.wrapper.ctx, pattern).Result()
	if err != nil {
		return fmt.Errorf("failed to get keys by pattern: %w", err)
	}

	if len(keys) > 0 {
		return r.wrapper.Del(r.wrapper.ctx, keys...).Err()
	}

	return nil
}

// SafeDeleteByPattern безопасное удаление по паттерну для кластера
func (r *RedisCache) SafeDeleteByPattern(pattern string) error {
	if r.wrapper.clientType == "cluster" {
		return r.deleteByPatternCluster(pattern)
	}
	return r.DeleteByPattern(pattern)
}

// deleteByPatternCluster реализует безопасное удаление для кластера
func (r *RedisCache) deleteByPatternCluster(pattern string) error {
	var cursor uint64
	var allKeys []string

	for {
		var keys []string
		var err error
		keys, cursor, err = r.wrapper.Scan(r.wrapper.ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}
	}

	if len(allKeys) > 0 {
		batchSize := 100
		for i := 0; i < len(allKeys); i += batchSize {
			end := i + batchSize
			if end > len(allKeys) {
				end = len(allKeys)
			}

			batch := allKeys[i:end]
			if err := r.wrapper.Del(r.wrapper.ctx, batch...).Err(); err != nil {
				log.Printf("Failed to delete batch %d-%d: %v", i, end, err)
			}
		}
	}

	return nil
}

// Exists проверяет существование ключа
func (r *RedisCache) Exists(key string) bool {
	return r.wrapper.Exists(r.wrapper.ctx, key).Val() > 0
}

// Close закрывает соединение с Redis
func (r *RedisCache) Close() error {
	return r.wrapper.Close()
}

// HealthCheck проверяет состояние Redis
func (r *RedisCache) HealthCheck() error {
	return r.wrapper.Ping(r.wrapper.ctx).Err()
}

// GetStats возвращает статистику Redis
func (r *RedisCache) GetStats() map[string]interface{} {
	stats := r.wrapper.PoolStats()

	return map[string]interface{}{
		"client_type": r.wrapper.clientType,
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"timeouts":    stats.Timeouts,
		"total_conns": stats.TotalConns,
		"idle_conns":  stats.IdleConns,
		"stale_conns": stats.StaleConns,
	}
}

// GetClientType возвращает тип клиента
func (r *RedisCache) GetClientType() string {
	return r.wrapper.GetClientType()
}
