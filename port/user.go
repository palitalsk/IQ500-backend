package port

import (
	"main/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	GetUserByName(email string) (domain.User, error)
	CreateUser(payload domain.User) error
	UpdatePasswordUser(id primitive.ObjectID, password string) error
	UpdateToken(id primitive.ObjectID, token string) error
}




