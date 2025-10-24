package service

import (
	"api/internal/models"
	"api/internal/repository"
)

type FriendService struct {
	friendRepo *repository.FriendRepository
	userRepo   *repository.UserRepository
}

func NewFriendService(friendRepo *repository.FriendRepository, userRepo *repository.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

// AddFriend добавляет друга
func (s *FriendService) AddFriend(userID, friendID int) error {
	return s.friendRepo.AddFriend(userID, friendID)
}

// DeleteFriend удаляет друга
func (s *FriendService) DeleteFriend(userID, friendID int) error {
	return s.friendRepo.DeleteFriend(userID, friendID)
}

// GetFriends возвращает список друзей
func (s *FriendService) GetFriends(userID int) ([]models.FriendResponse, error) {
	return s.friendRepo.GetFriends(userID)
}

// GetFriendshipStatus возвращает статус дружбы
func (s *FriendService) GetFriendshipStatus(userID, friendID int) (*models.FriendshipStatus, error) {
	return s.friendRepo.GetFriendshipStatus(userID, friendID)
}

// IsFriend проверяет, являются ли пользователи друзьями
func (s *FriendService) IsFriend(userID, friendID int) (bool, error) {
	return s.friendRepo.IsFriend(userID, friendID)
}
