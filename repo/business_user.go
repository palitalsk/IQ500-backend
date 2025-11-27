package repo

import (
	"context"

	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BusinessUserRepository struct {
	collection *mongo.Collection
}

func NewBusinessUserRepository(db *MongoDBRepository) *BusinessUserRepository {
	return &BusinessUserRepository{
		collection: db.GetCollection("business_users"),
	}
}

// CreateBusinessUser
func (r *BusinessUserRepository) CreateBusinessUser(ctx context.Context, businessUser *models.BusinessUser) error {
	businessUser.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, businessUser)
	return err
}

// GetBusinessUsersByUserID 
func (r *BusinessUserRepository) GetBusinessUsersByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.BusinessUser, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var businessUsers []models.BusinessUser
	if err = cursor.All(ctx, &businessUsers); err != nil {
		return nil, err
	}
	return businessUsers, nil
}
