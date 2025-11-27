package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessUser struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Username     string             `json:"username" bson:"username"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	BusinessCode string             `json:"business_code" bson:"business_code"`
	BusinessID   primitive.ObjectID `json:"business_id" bson:"business_id"`
	Role         string             `json:"role" bson:"role"` // owner, admin, member
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}


