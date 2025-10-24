package models

import "time"

type FriendRequest struct {
	// UserID   int `json:"user_id" binding:"required"`
	FriendID int `json:"friend_id" binding:"required"`
}

type FriendResponse struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	FriendID  int          `json:"friend_id"`
	Friend    UserResponse `json:"friend"`
	CreatedAt time.Time    `json:"created_at"`
}

type FriendshipStatus struct {
	IsFriend  bool `json:"is_friend"`
	IsPending bool `json:"is_pending"`
}
