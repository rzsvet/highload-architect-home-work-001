package service

import (
	"context"
	"dialog-service/internal/models"
	"dialog-service/internal/repository"
	"errors"
)

type DialogService interface {
	SendMessage(ctx context.Context, fromUserID int, toUserID int, text string) (*models.Message, error)
	GetDialog(ctx context.Context, userID1, userID2 int) ([]models.Message, error)
}

type dialogService struct {
	repo repository.DialogRepository
}

func NewDialogService(repo repository.DialogRepository) DialogService {
	return &dialogService{repo: repo}
}

func (s *dialogService) SendMessage(ctx context.Context, fromUserID, toUserID int, text string) (*models.Message, error) {
	if fromUserID == toUserID {
		return nil, errors.New("cannot send message to yourself")
	}

	if text == "" {
		return nil, errors.New("message text cannot be empty")
	}

	message := &models.Message{
		From: fromUserID,
		To:   toUserID,
		Text: text,
	}

	err := s.repo.SaveMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *dialogService) GetDialog(ctx context.Context, userID1, userID2 int) ([]models.Message, error) {
	if userID1 == userID2 {
		return nil, errors.New("user IDs must be different")
	}

	messages, err := s.repo.GetDialog(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
