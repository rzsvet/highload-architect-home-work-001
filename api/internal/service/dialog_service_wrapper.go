package service

import (
	"api/internal/client"
	"context"
)

type DialogService interface {
	SendMessage(ctx context.Context, fromUserID, toUserID int, text string) error
	GetDialog(ctx context.Context, userID1, userID2 int) ([]client.Message, error)
}

type dialogServiceImpl struct {
	dialogClient *client.DialogClient
}

func NewDialogService(dialogClient *client.DialogClient) DialogService {
	return &dialogServiceImpl{
		dialogClient: dialogClient,
	}
}

func (s *dialogServiceImpl) SendMessage(ctx context.Context, fromUserID, toUserID int, text string) error {
	_, err := s.dialogClient.SendMessage(ctx, fromUserID, toUserID, text)
	return err
}

func (s *dialogServiceImpl) GetDialog(ctx context.Context, userID1, userID2 int) ([]client.Message, error) {
	dialog, err := s.dialogClient.GetDialog(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return dialog.Messages, nil
}
