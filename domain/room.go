package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
	LastMessage string             `json:"last_message" bson:"last_message"`
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	CreateAt    time.Time          `json:"create_at" bson:"create_at"`
}

type PayloadRoom struct {
	Message string `json:"message"`
}

type ResponseCreateRoom struct {
	ID primitive.ObjectID `json:"id"`
}
