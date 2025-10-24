package cache

import (
	"api/internal/config"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisWrapper обертка для унификации интерфейса
type RedisWrapper struct {
	clientType string
	client     interface{} // Может быть *redis.Client или *redis.ClusterClient
	ctx        context.Context
}

func NewRedisWrapper(cfg *config.RedisConfig) (*RedisWrapper, error) {
	var client interface{}
	var clientType string

	ctx := context.Background()

	if cfg.ClusterMode {
		// Кластерная конфигурация
		if len(cfg.ClusterNodes) == 0 {
			return nil, fmt.Errorf("cluster nodes are required in cluster mode")
		}

		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        cfg.ClusterNodes,
			Password:     cfg.Password,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
		})

		client = clusterClient
		clientType = "cluster"
		log.Printf("Redis Cluster client initialized with nodes: %v", cfg.ClusterNodes)
	} else {
		// Standalone конфигурация
		addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
		standaloneClient := redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
		})

		client = standaloneClient
		clientType = "standalone"
		log.Printf("Redis Standalone client initialized: %s", addr)
	}

	// Проверяем соединение
	var pingErr error
	switch c := client.(type) {
	case *redis.Client:
		pingErr = c.Ping(ctx).Err()
	case *redis.ClusterClient:
		pingErr = c.Ping(ctx).Err()
	default:
		return nil, fmt.Errorf("unknown client type")
	}

	if pingErr != nil {
		return nil, fmt.Errorf("failed to connect to Redis (%s): %w", clientType, pingErr)
	}

	log.Printf("Redis %s connection established successfully", strings.ToUpper(clientType))
	return &RedisWrapper{
		clientType: clientType,
		client:     client,
		ctx:        ctx,
	}, nil
}

// Универсальные методы для работы с клиентом

func (w *RedisWrapper) Get(ctx context.Context, key string) *redis.StringCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Get(ctx, key)
	case *redis.ClusterClient:
		return c.Get(ctx, key)
	default:
		return nil
	}
}

func (w *RedisWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Set(ctx, key, value, expiration)
	case *redis.ClusterClient:
		return c.Set(ctx, key, value, expiration)
	default:
		return nil
	}
}

func (w *RedisWrapper) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Del(ctx, keys...)
	case *redis.ClusterClient:
		return c.Del(ctx, keys...)
	default:
		return nil
	}
}

func (w *RedisWrapper) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Exists(ctx, keys...)
	case *redis.ClusterClient:
		return c.Exists(ctx, keys...)
	default:
		return nil
	}
}

func (w *RedisWrapper) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Keys(ctx, pattern)
	case *redis.ClusterClient:
		return c.Keys(ctx, pattern)
	default:
		return nil
	}
}

func (w *RedisWrapper) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Scan(ctx, cursor, match, count)
	case *redis.ClusterClient:
		return c.Scan(ctx, cursor, match, count)
	default:
		return nil
	}
}

func (w *RedisWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Ping(ctx)
	case *redis.ClusterClient:
		return c.Ping(ctx)
	default:
		return nil
	}
}

func (w *RedisWrapper) PoolStats() *redis.PoolStats {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.PoolStats()
	case *redis.ClusterClient:
		return c.PoolStats()
	default:
		return nil
	}
}

func (w *RedisWrapper) Close() error {
	switch c := w.client.(type) {
	case *redis.Client:
		return c.Close()
	case *redis.ClusterClient:
		return c.Close()
	default:
		return fmt.Errorf("unknown client type")
	}
}

func (w *RedisWrapper) GetClientType() string {
	return w.clientType
}

func (w *RedisWrapper) GetContext() context.Context {
	return w.ctx
}
