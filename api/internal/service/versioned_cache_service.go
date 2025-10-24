package service

import (
	"api/internal/cache"
	"api/internal/models"
	"fmt"
)

type VersionedCacheService struct {
	cache   *cache.RedisCache
	version *CacheVersionService
}

func NewVersionedCacheService(redisCache *cache.RedisCache) *VersionedCacheService {
	return &VersionedCacheService{
		cache:   redisCache,
		version: NewCacheVersionService(redisCache),
	}
}

// GetFeedFromCache получает ленту из кэша с проверкой версии
func (v *VersionedCacheService) GetFeedFromCache(userID, page, pageSize int) (*models.FeedResponse, error) {
	baseKey := fmt.Sprintf("feed:user:%d:page:%d:size:%d", userID, page, pageSize)
	versionedKey, err := v.version.GenerateVersionedKey(baseKey, userID)
	if err != nil {
		return nil, err
	}

	var feed models.FeedResponse
	if err := v.cache.Get(versionedKey, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// SetFeedToCache сохраняет ленту в кэш с версией
func (v *VersionedCacheService) SetFeedToCache(userID, page, pageSize int, feed *models.FeedResponse) error {
	baseKey := fmt.Sprintf("feed:user:%d:page:%d:size:%d", userID, page, pageSize)
	versionedKey, err := v.version.GenerateVersionedKey(baseKey, userID)
	if err != nil {
		return err
	}

	return v.cache.Set(versionedKey, feed, FeedCacheTTL)
}

// InvalidateUserCache инвалидирует ВСЕ кэши пользователя через смену версии
func (v *VersionedCacheService) InvalidateUserCache(userID int) error {
	_, err := v.version.IncrementCacheVersion(userID)
	return err
}
