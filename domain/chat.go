package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	RoomID      primitive.ObjectID `json:"room_id" bson:"room_id"`
	Message     string             `json:"message" bson:"message"`
	Type        string             `json:"type" bson:"type"`                 // chat, reply
	MessageType string             `json:"message_type" bson:"message_type"` // text, image
	UpdateAt    time.Time          `json:"update_at" bson:"update_at"`
	CreateAt    time.Time          `json:"create_at" bson:"create_at"`
}

type PayloadChat struct {
	RoomID  string `json:"room_id" binding:"required"`
	Message string `json:"message"`
	Img     string `json:"img"`
}




