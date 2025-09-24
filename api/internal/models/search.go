package models

type UserSearchRequest struct {
	FirstName string `form:"first_name" binding:"required"`
	LastName  string `form:"last_name" binding:"required"`
}

type UserSearchResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total"`
}
