package handler

import (
	"api/internal/dialog"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DialogHandler struct {
	dialogService dialog.Service
}

func NewDialogHandler(dialogService dialog.Service) *DialogHandler {
	return &DialogHandler{
		dialogService: dialogService,
	}
}

type SendMessageRequest struct {
	To   int    `json:"to" binding:"required"`
	Text string `json:"text" binding:"required"`
}

func (h *DialogHandler) SendMessage(c *gin.Context) {
	// Получаем ID пользователя из контекста (зависит от вашей реализации аутентификации)
	fromUserID, err := getUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.dialogService.SendMessage(c.Request.Context(), fromUserID, req.To, req.Text); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message sent"})
}

func (h *DialogHandler) GetDialog(c *gin.Context) {
	fromUserID, err := getUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	toUserIDStr := c.Param("user_id")
	toUserID, err := strconv.Atoi(toUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	messages, err := h.dialogService.GetDialog(c.Request.Context(), fromUserID, toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func getUserIdFromContext(c *gin.Context) (int, error) {
	// Пример получения ID пользователя из JWT токена
	// Замените на вашу реализацию аутентификации
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("no authorization header")
	}

	// Здесь должна быть логика парсинга JWT и получения user_id
	// Для примера, просто возвращаем тестовое значение
	return 1, nil
}
