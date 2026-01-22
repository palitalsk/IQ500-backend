package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Token     string             `json:"token" bson:"token"`
	Active    bool               `json:"active" bson:"active"`
	UpdateAt  time.Time          `json:"update_at" bson:"update_at"`
	CreateAt  time.Time          `json:"create_at" bson:"create_at"`
	LastLogin time.Time          `json:"last_login" bson:"last_login"`
}

type PayloadUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PayloadResetPassword struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
