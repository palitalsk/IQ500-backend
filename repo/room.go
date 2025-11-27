package repo

import (
	"context"
	"fmt"
	"main/domain"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ROOM_COLLECTION = "room"
)

type RoomRepository struct {
	DB *mongo.Database
}

func NewRoomRepository(DB *mongo.Database) *RoomRepository {
	return &RoomRepository{
		DB,
	}
}

func (r *RoomRepository) GetRoomByUserId(id primitive.ObjectID) ([]domain.Room, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	option := options.Find()
	option.SetSort(bson.M{
		"update_at": -1,
	})
	defer cancel()

	var room []domain.Room

	filter := bson.M{
		"user_id": id,
	}

	obj, err := r.DB.Collection(ROOM_COLLECTION).Find(c, filter, option)
	if err != nil {
		fmt.Println("get room by user id err: ", err)
		return nil, err
	}

	err = obj.All(c, &room)
	if err != nil {
		fmt.Println("get room by user id err: ", err)
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) CreateRoom(payload domain.Room) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := r.DB.Collection(ROOM_COLLECTION).InsertOne(c, payload)

	if err != nil {
		fmt.Println("create room err: ", err)
		return err
	}
	return nil
}

func (r *RoomRepository) DeleteRoomByID(id primitive.ObjectID) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	_, err := r.DB.Collection(ROOM_COLLECTION).DeleteOne(c, filter)
	if err != nil {
		fmt.Println("delete room by id err: ", err)
		return err
	}
	return nil
}

func (r *RoomRepository) UpdateLastMessage(id primitive.ObjectID, lastMessage string) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	set := bson.M{
		"$set": bson.M{
			"last_message": lastMessage,
			"update_at":    time.Now(),
		},
	}

	_, err := r.DB.Collection(ROOM_COLLECTION).UpdateOne(c, filter, set)
	if err != nil {
		fmt.Println("update last message err: ", err)
		return err
	}
	return nil
}



