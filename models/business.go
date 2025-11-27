package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Business struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Code         string             `json:"code" bson:"code"`
	Name         string             `json:"name" bson:"name"`
	BusinessType string             `json:"business_type" bson:"business_type"`
	Detail       string             `json:"detail" bson:"detail"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

type CreateBusinessRequest struct {
	Name         string `json:"name" binding:"required"`
	BusinessType string `json:"business_type" binding:"required"`
	Detail       string `json:"detail"`
}

type BusinessResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Code         string             `json:"code"`
	Name         string             `json:"name"`
	BusinessType string             `json:"business_type"`
	Detail       string             `json:"detail"`
	CreatedAt    time.Time          `json:"created_at"`
}

type BusinessListResponse struct {
	Businesses []BusinessResponse `json:"businesses"`
	Total      int                `json:"total"`
}
