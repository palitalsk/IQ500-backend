package port

import (
	"main/internal/core/domain"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepository interface {
	GetChatByRoomID(id primitive.ObjectID) ([]domain.Chat, error)
	CreateChat(payload domain.Chat) error
}

type ChatService interface {
	GetChatByRoomID(c *gin.Context, id primitive.ObjectID) ([]domain.Chat, error)
	Chat(c *gin.Context, payload domain.PayloadChat) (*domain.AIPredictionResult, error)
}
