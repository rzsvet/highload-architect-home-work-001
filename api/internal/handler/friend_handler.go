package handler

import (
	"api/internal/models"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FriendHandler struct {
	friendService *service.FriendService
}

func NewFriendHandler(friendService *service.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

// AddFriend godoc
// @Summary Добавить друга
// @Description Добавляет пользователя в друзья
// @Tags Friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.FriendRequest true "Данные для добавления друга"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /friend/add [post]
func (h *FriendHandler) AddFriend(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.FriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.friendService.AddFriend(userID, req.FriendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

// DeleteFriend godoc
// @Summary Удалить друга
// @Description Удаляет пользователя из друзей
// @Tags Friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.FriendRequest true "Данные для удаления друга"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /friend/delete [post]
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.FriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	if err := h.friendService.DeleteFriend(userID, req.FriendID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend deleted successfully"})
}

// GetFriends godoc
// @Summary Получить список друзей
// @Description Возвращает список друзей пользователя
// @Tags Friends
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.FriendResponse
// @Failure 500 {object} map[string]string
// @Router /friends [get]
func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	friends, err := h.friendService.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friends)
}

// GetFriendshipStatus godoc
// @Summary Получить статус дружбы
// @Description Возвращает статус дружбы с указанным пользователем
// @Tags Friends
// @Produce json
// @Security BearerAuth
// @Param friend_id query int true "ID пользователя"
// @Success 200 {object} models.FriendshipStatus
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /friend/status [get]
func (h *FriendHandler) GetFriendshipStatus(c *gin.Context) {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	friendIDStr := c.Query("friend_id")
	friendID, err := strconv.Atoi(friendIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}

	status, err := h.friendService.GetFriendshipStatus(userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
