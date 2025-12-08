package handler

import (
	"dialog-service/internal/models"
	"dialog-service/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DialogHandler struct {
	service service.DialogService
}

func NewDialogHandler(service service.DialogService) *DialogHandler {
	return &DialogHandler{service: service}
}

// GetUserIDFromContext - получение ID пользователя из заголовка X-User-Id
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userIDStr := c.GetHeader("X-User-Id")
	if userIDStr == "" {
		return 0, c.AbortWithError(http.StatusUnauthorized, gin.Error{})
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, c.AbortWithError(http.StatusBadRequest, gin.Error{})
	}

	return userID, nil
}

func (h *DialogHandler) SendMessage(c *gin.Context) {
	fromUserID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	toUserIDStr := c.Param("user_id")
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var request models.SendMessageRequest
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	message, err := h.service.SendMessage(c.Request.Context(), fromUserID, toUserID, request.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *DialogHandler) GetDialog(c *gin.Context) {
	fromUserID, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	toUserIDStr := c.Param("user_id")
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	messages, err := h.service.GetDialog(c.Request.Context(), fromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.DialogResponse{
		Messages: messages,
	}

	c.JSON(http.StatusOK, response)
}

func (h *DialogHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
