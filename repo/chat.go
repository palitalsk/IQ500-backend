package repo

import (
	"context"
	"main/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CHAT_COLLECTION = "chat"
)

type ChatRepository struct {
	DB *mongo.Database
}

func NewChatRepository(DB *mongo.Database) *ChatRepository {
	return &ChatRepository{
		DB,
	}
}

func (r *ChatRepository) GetChatByRoomID(id primitive.ObjectID) ([]domain.Chat, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	option := options.Find()
	option.SetSort(bson.M{
		"update_at": -1,
	})
	defer cancel()

	var chat []domain.Chat

	filter := bson.M{
		"room_id": id,
	}

	obj, err := r.DB.Collection(CHAT_COLLECTION).Find(c, filter, option)
	if err != nil {
		return nil, err
	}

	err = obj.All(c, &chat)
	if err != nil {
		return nil, err
	}

	return chat, nil
}

func (r *ChatRepository) CreateChat(payload domain.Chat) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := r.DB.Collection(CHAT_COLLECTION).InsertOne(c, payload)

	if err != nil {
		return err
	}
	return nil
}



