package handler

import (
	"api/internal/models"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost godoc
// @Summary Создать пост
// @Description Создает новый пост пользователя
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePostRequest true "Данные поста"
// @Success 201 {object} models.Post
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/create [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	post, err := h.postService.CreatePost(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// GetPost godoc
// @Summary Получить пост
// @Description Возвращает пост по ID
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Success 200 {object} models.PostResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/get/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.postService.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// UpdatePost godoc
// @Summary Обновить пост
// @Description Обновляет существующий пост
// @Tags Posts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Param request body models.UpdatePostRequest true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/update/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req models.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.postService.UpdatePost(postID, userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
}

// DeletePost godoc
// @Summary Удалить пост
// @Description Удаляет пост пользователя
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID поста"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /post/delete/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.postService.DeletePost(postID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetUserPosts godoc
// @Summary Получить посты пользователя
// @Description Возвращает посты указанного пользователя
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "ID пользователя (по умолчанию - текущий)"
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(20)
// @Success 200 {object} models.FeedResponse
// @Failure 500 {object} map[string]string
// @Router /posts [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	var targetUserID int
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Если указан user_id, используем его, иначе используем текущего пользователя
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		targetUserID, err = strconv.Atoi(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	} else {
		targetUserID = userID
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	feed, err := h.postService.GetUserPosts(targetUserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feed)
}

// GetFeed godoc
// @Summary Лента постов друзей
// @Description Возвращает ленту постов друзей пользователя
// @Tags Posts
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(20)
// @Success 200 {object} models.FeedResponse
// @Failure 500 {object} map[string]string
// @Router /post/feed [get]
func (h *PostHandler) GetFeed(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	feed, err := h.postService.GetFriendsPosts(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feed)
}
