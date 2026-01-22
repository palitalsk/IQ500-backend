package repo

import (
	"context"
	"errors"
	"main/domain"
	"main/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	USER_COLLECTION = "user_account"
)

type UserRepository struct {
	DB *mongo.Database
}

func NewUserRepository(DB *mongo.Database) *UserRepository {
	return &UserRepository{
		DB,
	}
}

func (r *UserRepository) GetUserByName(username string) (domain.User, error) {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user domain.User

	option := options.FindOne()
	option.SetSort(nil)

	filter := bson.M{
		"username": username,
	}

	err := r.DB.Collection(USER_COLLECTION).FindOne(c, filter, option).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (r *UserRepository) CreateUser(payload domain.User) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := r.DB.Collection(USER_COLLECTION).InsertOne(c, payload)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdatePasswordUser(id primitive.ObjectID, password string) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	option := options.Update()
	option.SetUpsert(true)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	set := bson.M{
		"$set": bson.M{
			"password":  utils.MD5Hash(password),
			"update_at": time.Now(),
		},
	}

	obj, err := r.DB.Collection(USER_COLLECTION).UpdateOne(c, filter, set, option)
	if err != nil {
		return err
	}

	if obj == nil {
		err := errors.New("obj is nil")
		return err
	}

	if obj.MatchedCount == 0 {
		err := errors.New("not found document")
		return err
	}

	return nil
}

func (r *UserRepository) UpdateToken(id primitive.ObjectID, token string) error {
	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	option := options.Update()
	option.SetUpsert(true)
	defer cancel()

	filter := bson.M{
		"_id": id,
	}

	set := bson.M{
		"$set": bson.M{
			"token":     token,
			"update_at": time.Now(),
		},
	}

	obj, err := r.DB.Collection(USER_COLLECTION).UpdateOne(c, filter, set, option)
	if err != nil {
		return err
	}

	if obj == nil {
		err := errors.New("obj is nil")
		return err
	}

	if obj.MatchedCount == 0 {
		err := errors.New("not found document")
		return err
	}

	return nil
}
