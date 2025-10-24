package service

import (
	"api/internal/models"
	"api/internal/repository"
	"log"
)

type PostService struct {
	postRepo     *repository.PostRepository
	cacheService *CacheService
}

func NewPostService(postRepo *repository.PostRepository, cacheService *CacheService) *PostService {
	return &PostService{
		postRepo:     postRepo,
		cacheService: cacheService,
	}
}

// CreatePost создает новый пост с инвалидацией кэша
func (s *PostService) CreatePost(userID int, req *models.CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	err := s.postRepo.CreatePost(post)
	if err != nil {
		return nil, err
	}

	// Инвалидируем кэш ленты друзей и постов пользователя
	go func() {
		if err := s.cacheService.InvalidateUserFeedCache(userID); err != nil {
			log.Printf("Failed to invalidate feed cache for user %d: %v", userID, err)
		}
		if err := s.cacheService.InvalidateUserPostsCache(userID); err != nil {
			log.Printf("Failed to invalidate posts cache for user %d: %v", userID, err)
		}
	}()

	return post, nil
}

// GetPost возвращает пост по ID (без кэширования, так как редко запрашиваются по одному)
func (s *PostService) GetPost(postID int) (*models.PostResponse, error) {
	return s.postRepo.GetPost(postID)
}

// UpdatePost обновляет пост с инвалидацией кэша
func (s *PostService) UpdatePost(postID, userID int, req *models.UpdatePostRequest) error {
	err := s.postRepo.UpdatePost(postID, userID, req)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	go func() {
		if err := s.cacheService.InvalidateUserFeedCache(userID); err != nil {
			log.Printf("Failed to invalidate feed cache for user %d: %v", userID, err)
		}
		if err := s.cacheService.InvalidateUserPostsCache(userID); err != nil {
			log.Printf("Failed to invalidate posts cache for user %d: %v", userID, err)
		}
	}()

	return nil
}

// DeletePost удаляет пост с инвалидацией кэша
func (s *PostService) DeletePost(postID, userID int) error {
	err := s.postRepo.DeletePost(postID, userID)
	if err != nil {
		return err
	}

	// Инвалидируем кэш
	go func() {
		if err := s.cacheService.InvalidateUserFeedCache(userID); err != nil {
			log.Printf("Failed to invalidate feed cache for user %d: %v", userID, err)
		}
		if err := s.cacheService.InvalidateUserPostsCache(userID); err != nil {
			log.Printf("Failed to invalidate posts cache for user %d: %v", userID, err)
		}
	}()

	return nil
}

// GetUserPosts возвращает посты пользователя с кэшированием
func (s *PostService) GetUserPosts(userID, page, pageSize int) (*models.FeedResponse, error) {
	// Пытаемся получить из кэша
	if cached, err := s.cacheService.GetUserPostsFromCache(userID, page, pageSize); err == nil {
		log.Printf("Cache HIT for user posts: user=%d, page=%d", userID, page)
		return cached, nil
	}

	log.Printf("Cache MISS for user posts: user=%d, page=%d", userID, page)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	posts, total, err := s.postRepo.GetUserPosts(userID, pageSize, offset)
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
		if err := s.cacheService.SetUserPostsToCache(userID, page, pageSize, feed); err != nil {
			log.Printf("Failed to cache user posts for user %d: %v", userID, err)
		}
	}()

	return feed, nil
}

// GetFriendsPosts возвращает ленту постов друзей с кэшированием
func (s *PostService) GetFriendsPosts(userID, page, pageSize int) (*models.FeedResponse, error) {
	// Пытаемся получить из кэша
	if cached, err := s.cacheService.GetFeedFromCache(userID, page, pageSize); err == nil {
		log.Printf("Cache HIT for feed: user=%d, page=%d", userID, page)
		return cached, nil
	}

	log.Printf("Cache MISS for feed: user=%d, page=%d", userID, page)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	posts, total, err := s.postRepo.GetFriendsPosts(userID, pageSize, offset)
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
		if err := s.cacheService.SetFeedToCache(userID, page, pageSize, feed); err != nil {
			log.Printf("Failed to cache feed for user %d: %v", userID, err)
		}
	}()

	return feed, nil
}
