package repository

import (
	"context"
	"dialog-service/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DialogRepository interface {
	SaveMessage(ctx context.Context, message *models.Message) error
	GetDialog(ctx context.Context, user1, user2 int) ([]models.Message, error)
}

type MongoDBRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBRepository(uri, dbName string) (*MongoDBRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	collection := db.Collection("messages")

	// Создаем индекс для быстрого поиска диалогов
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "from", Value: 1},
			{Key: "to", Value: 1},
			{Key: "created_at", Value: -1},
		},
	}
	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, err
	}

	return &MongoDBRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *MongoDBRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, message)
	return err
}

func (r *MongoDBRepository) GetDialog(ctx context.Context, user1, user2 int) ([]models.Message, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"from": user1, "to": user2},
			{"from": user2, "to": user1},
		},
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []models.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MongoDBRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return r.client.Disconnect(ctx)
}
