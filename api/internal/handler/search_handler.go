package handler

import (
	"api/internal/models"
	"api/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	userService *service.UserService
}

func NewSearchHandler(userService *service.UserService) *SearchHandler {
	return &SearchHandler{userService: userService}
}

// SearchUsers godoc
// @Summary Поиск пользователей
// @Description Поиск анкет пользователей по имени и фамилии
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param first_name query string true "Часть имени для поиска" example("Конст")
// @Param last_name query string true "Часть фамилии для поиска" example("Оси")
// @Param page query int false "Номер страницы" default(1)
// @Param page_size query int false "Размер страницы" default(20)
// @Success 200 {object} models.UserSearchResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/search [get]
func (h *SearchHandler) SearchUsers(c *gin.Context) {
	var searchReq models.UserSearchRequest
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameters: " + err.Error()})
		return
	}

	// Получаем параметры пагинации
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Используем поиск с пагинацией
	result, err := h.userService.SearchUsersWithPaging(searchReq.FirstName, searchReq.LastName, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SearchUsersSimple godoc
// @Summary Простой поиск пользователей
// @Description Поиск анкет пользователей по имени и фамилии (без пагинации)
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param first_name query string true "Часть имени для поиска" example("Конст")
// @Param last_name query string true "Часть фамилии для поиска" example("Оси")
// @Success 200 {array} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/search/simple [get]
func (h *SearchHandler) SearchUsersSimple(c *gin.Context) {
	var searchReq models.UserSearchRequest
	if err := c.ShouldBindQuery(&searchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid search parameters: " + err.Error()})
		return
	}

	users, err := h.userService.SearchUsers(searchReq.FirstName, searchReq.LastName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
