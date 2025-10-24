package handler

import (
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CacheHandler struct {
	cacheService *service.CacheService
	postService  *service.PostService
}

func NewCacheHandler(cacheService *service.CacheService, postService *service.PostService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
		postService:  postService,
	}
}

// InvalidateCache godoc
// @Summary Инвалидировать кэш пользователя
// @Description Принудительно обновляет кэш для указанного пользователя
// @Tags Cache
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "ID пользователя (по умолчанию - текущий)"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cache/invalidate [post]
func (h *CacheHandler) InvalidateCache(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Если указан другой user_id (для админов)
	if targetUserIDStr := c.Query("user_id"); targetUserIDStr != "" {
		userID, err = strconv.Atoi(targetUserIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
	}

	if err := h.cacheService.RefreshCache(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache invalidated successfully",
		"user_id": userID,
	})
}

// GetCacheStats godoc
// @Summary Получить статистику кэша
// @Description Возвращает статистику и метрики кэша
// @Tags Cache
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /cache/stats [get]
func (h *CacheHandler) GetCacheStats(c *gin.Context) {
	stats, err := h.cacheService.GetCacheStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
