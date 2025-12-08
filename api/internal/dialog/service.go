package dialog

import (
	"context"
)

type Service interface {
	SendMessage(ctx context.Context, fromUserID, toUserID int, text string) error
	GetDialog(ctx context.Context, userID1, userID2 int) ([]Message, error)
}

type serviceImpl struct {
	client *DialogClient
}

func NewService(client *DialogClient) Service {
	return &serviceImpl{
		client: client,
	}
}

func (s *serviceImpl) SendMessage(ctx context.Context, fromUserID, toUserID int, text string) error {
	_, err := s.client.SendMessage(ctx, fromUserID, toUserID, text)
	return err
}

func (s *serviceImpl) GetDialog(ctx context.Context, userID1, userID2 int) ([]Message, error) {
	dialog, err := s.client.GetDialog(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}
	return dialog.Messages, nil
}
