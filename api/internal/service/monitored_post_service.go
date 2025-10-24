package service

import (
	"api/internal/models"
	"api/internal/monitoring"
	"api/internal/repository"
	"log"
)

type MonitoredPostService struct {
	postRepo     *repository.PostRepository
	cacheService *CacheService
}

func NewMonitoredPostService(postRepo *repository.PostRepository, cacheService *CacheService) *MonitoredPostService {
	return &MonitoredPostService{
		postRepo:     postRepo,
		cacheService: cacheService,
	}
}

// GetFriendsPosts возвращает ленту постов друзей с кэшированием и метриками
func (m *MonitoredPostService) GetFriendsPosts(userID, page, pageSize int) (*models.FeedResponse, error) {
	// Пытаемся получить из кэша
	if cached, err := m.cacheService.GetFeedFromCache(userID, page, pageSize); err == nil {
		log.Printf("Cache HIT for feed: user=%d, page=%d", userID, page)
		monitoring.RecordCacheHit("feed")
		return cached, nil
	}

	log.Printf("Cache MISS for feed: user=%d, page=%d", userID, page)
	monitoring.RecordCacheMiss("feed")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	posts, total, err := m.postRepo.GetFriendsPosts(userID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	pages := (total + pageSize - 1) / pageSize

	feed := &models.FeedResponse{
		Posts: posts,
		Total: total,
		Page:  page,
		Pages: pages,
	}

	// Сохраняем в кэш асинхронно
	go func() {
		if err := m.cacheService.SetFeedToCache(userID, page, pageSize, feed); err != nil {
			log.Printf("Failed to cache feed for user %d: %v", userID, err)
		}
	}()

	return feed, nil
}

// CreatePost создает пост с инвалидацией кэша и метриками
func (m *MonitoredPostService) CreatePost(userID int, req *models.CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	err := m.postRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш и записываем метрики
	go func() {
		if err := m.cacheService.InvalidateUserFeedCache(userID); err != nil {
			log.Printf("Failed to invalidate feed cache for user %d: %v", userID, err)
		} else {
			monitoring.RecordCacheInvalidation("feed")
		}

		if err := m.cacheService.InvalidateUserPostsCache(userID); err != nil {
			log.Printf("Failed to invalidate posts cache for user %d: %v", userID, err)
		} else {
			monitoring.RecordCacheInvalidation("user_posts")
		}
	}()

	return post, nil
}
