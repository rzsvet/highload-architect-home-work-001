package service

import (
	"api/internal/cache"
	"api/internal/models"
	"fmt"
	"time"
)

type CacheService struct {
	cache *cache.RedisCache
}

func NewCacheService(redisCache *cache.RedisCache) *CacheService {
	return &CacheService{
		cache: redisCache,
	}
}

// Константы TTL
const (
	FeedCacheTTL      = 5 * time.Minute
	UserPostsCacheTTL = 10 * time.Minute
	UserCacheTTL      = 30 * time.Minute
	SearchCacheTTL    = 2 * time.Minute
)

// GenerateFeedCacheKey генерирует ключ для кэша ленты
func (s *CacheService) GenerateFeedCacheKey(userID, page, pageSize int) string {
	return fmt.Sprintf("feed:user:%d:page:%d:size:%d", userID, page, pageSize)
}

// GenerateUserPostsCacheKey генерирует ключ для кэша постов пользователя
func (s *CacheService) GenerateUserPostsCacheKey(userID, page, pageSize int) string {
	return fmt.Sprintf("posts:user:%d:page:%d:size:%d", userID, page, pageSize)
}

// GetFeedFromCache получает ленту из кэша
func (s *CacheService) GetFeedFromCache(userID, page, pageSize int) (*models.FeedResponse, error) {
	key := s.GenerateFeedCacheKey(userID, page, pageSize)

	var feed models.FeedResponse
	if err := s.cache.Get(key, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// SetFeedToCache сохраняет ленту в кэш
func (s *CacheService) SetFeedToCache(userID, page, pageSize int, feed *models.FeedResponse) error {
	key := s.GenerateFeedCacheKey(userID, page, pageSize)
	return s.cache.Set(key, feed, FeedCacheTTL)
}

// GetUserPostsFromCache получает посты пользователя из кэша
func (s *CacheService) GetUserPostsFromCache(userID, page, pageSize int) (*models.FeedResponse, error) {
	key := s.GenerateUserPostsCacheKey(userID, page, pageSize)

	var feed models.FeedResponse
	if err := s.cache.Get(key, &feed); err != nil {
		return nil, err
	}

	return &feed, nil
}

// SetUserPostsToCache сохраняет посты пользователя в кэш
func (s *CacheService) SetUserPostsToCache(userID, page, pageSize int, feed *models.FeedResponse) error {
	key := s.GenerateUserPostsCacheKey(userID, page, pageSize)
	return s.cache.Set(key, feed, UserPostsCacheTTL)
}

// InvalidateUserFeedCache инвалидирует кэш ленты пользователя
func (s *CacheService) InvalidateUserFeedCache(userID int) error {
	pattern := fmt.Sprintf("feed:user:%d:*", userID)

	// Используем безопасное удаление для кластера
	if s.cache.GetClientType() == "cluster" {
		return s.cache.SafeDeleteByPattern(pattern)
	}
	return s.cache.DeleteByPattern(pattern)
}

// InvalidateUserPostsCache инвалидирует кэш постов пользователя
func (s *CacheService) InvalidateUserPostsCache(userID int) error {
	pattern := fmt.Sprintf("posts:user:%d:*", userID)

	if s.cache.GetClientType() == "cluster" {
		return s.cache.SafeDeleteByPattern(pattern)
	}
	return s.cache.DeleteByPattern(pattern)
}

// InvalidateUserCache инвалидирует кэш данных пользователя
func (s *CacheService) InvalidateUserCache(userID int) error {
	key := fmt.Sprintf("user:%d", userID)
	return s.cache.Delete(key)
}

// InvalidateSearchCache инвалидирует кэш поиска
func (s *CacheService) InvalidateSearchCache(firstName, lastName string) error {
	pattern := fmt.Sprintf("search:%s:%s:*", firstName, lastName)

	if s.cache.GetClientType() == "cluster" {
		return s.cache.SafeDeleteByPattern(pattern)
	}
	return s.cache.DeleteByPattern(pattern)
}

// GetCacheStats возвращает статистику кэша
func (s *CacheService) GetCacheStats() (map[string]interface{}, error) {
	stats := s.cache.GetStats()
	stats["status"] = "active"
	return stats, nil
}

// HealthCheck проверяет состояние кэша
func (s *CacheService) HealthCheck() error {
	return s.cache.HealthCheck()
}

// RefreshCache принудительно обновляет кэш пользователя
func (s *CacheService) RefreshCache(userID int) error {
	errors := []error{}

	if err := s.InvalidateUserFeedCache(userID); err != nil {
		errors = append(errors, fmt.Errorf("feed cache: %w", err))
	}

	if err := s.InvalidateUserPostsCache(userID); err != nil {
		errors = append(errors, fmt.Errorf("posts cache: %w", err))
	}

	if err := s.InvalidateUserCache(userID); err != nil {
		errors = append(errors, fmt.Errorf("user cache: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("cache refresh errors: %v", errors)
	}

	return nil
}
