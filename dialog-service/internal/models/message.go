package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	From      int                `json:"from" bson:"from" binding:"required"`
	To        int                `json:"to" bson:"to" binding:"required"`
	Text      string             `json:"text" bson:"text" binding:"required"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
}

type SendMessageRequest struct {
	Text string `json:"text" binding:"required"`
}

type DialogResponse struct {
	Messages []Message `json:"messages"`
}
