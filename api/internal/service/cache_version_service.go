package service

import (
	"api/internal/cache"
	"fmt"
	"strconv"
)

type CacheVersionService struct {
	cache *cache.RedisCache
}

func NewCacheVersionService(redisCache *cache.RedisCache) *CacheVersionService {
	return &CacheVersionService{
		cache: redisCache,
	}
}

// GetCacheVersion получает текущую версию кэша для пользователя
func (c *CacheVersionService) GetCacheVersion(userID int) (int, error) {
	key := fmt.Sprintf("cache_version:user:%d", userID)

	// Используем существующие методы RedisCache вместо прямого доступа
	var versionStr string
	if err := c.cache.Get(key, &versionStr); err != nil {
		if err.Error() == "cache miss" {
			// Если версии нет, создаем начальную
			return c.SetCacheVersion(userID, 1)
		}
		return 0, err
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, err
	}

	return version, nil
}

// SetCacheVersion устанавливает новую версию кэша
func (c *CacheVersionService) SetCacheVersion(userID, version int) (int, error) {
	key := fmt.Sprintf("cache_version:user:%d", userID)
	err := c.cache.Set(key, strconv.Itoa(version), 0) // Бессрочное хранение
	return version, err
}

// IncrementCacheVersion увеличивает версию кэша
func (c *CacheVersionService) IncrementCacheVersion(userID int) (int, error) {
	// Получаем текущую версию
	currentVersion, err := c.GetCacheVersion(userID)
	if err != nil {
		// Если ошибка, создаем новую версию
		return c.SetCacheVersion(userID, 1)
	}

	// Увеличиваем версию
	newVersion := currentVersion + 1
	return c.SetCacheVersion(userID, newVersion)
}

// GenerateVersionedKey генерирует ключ с версией
func (c *CacheVersionService) GenerateVersionedKey(baseKey string, userID int) (string, error) {
	version, err := c.GetCacheVersion(userID)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:v%d", baseKey, version), nil
}
