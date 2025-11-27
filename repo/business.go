package repo

import (
	"context"

	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BusinessRepository struct {
	collection *mongo.Collection
}

func NewBusinessRepository(db *MongoDBRepository) *BusinessRepository {
	return &BusinessRepository{
		collection: db.GetCollection("businesses"),
	}
}

// CreateBusiness
func (r *BusinessRepository) CreateBusiness(ctx context.Context, business *models.Business) error {
	business.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(ctx, business)
	return err
}

// GetBusinessByID
func (r *BusinessRepository) GetBusinessByID(ctx context.Context, businessID primitive.ObjectID) (*models.Business, error) {
	var business models.Business
	err := r.collection.FindOne(ctx, bson.M{"_id": businessID}).Decode(&business)
	if err != nil {
		return nil, err
	}
	return &business, nil
}

// BusinessExists >> เอาไว้ตรวจสอบว่ามี business code นี้อยู่ในระบบหรือยัง
func (r *BusinessRepository) BusinessExists(ctx context.Context, code string) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"code": code})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetBusinessByCode ดึงข้อมูล business โดยใช้ code
func (r *BusinessRepository) GetBusinessByCode(ctx context.Context, code string) (*models.Business, error) {
	var business models.Business
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&business)
	if err != nil {
		return nil, err
	}
	return &business, nil
}
