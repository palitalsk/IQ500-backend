package service

import (
	"errors"
	"fmt"
	"main/internal/core/domain"
	"main/internal/core/port"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatService struct {
	chatRepo  port.ChatRepository
	roomRepo  port.RoomRepository
	aiService port.AIService
}

func NewChatService(
	chatRepo port.ChatRepository,
	roomRepo port.RoomRepository,
	aiService port.AIService,
) *ChatService {
	return &ChatService{
		chatRepo,
		roomRepo,
		aiService,
	}
}

func (s *ChatService) GetChatByRoomID(c *gin.Context, id primitive.ObjectID) ([]domain.Chat, error) {
	chat, err := s.chatRepo.GetChatByRoomID(id)
	if err != nil {
		fmt.Println("error get chat", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, err
	}

	return chat, nil
}

func (s *ChatService) Chat(c *gin.Context, payload domain.PayloadChat) (*domain.AIPredictionResult, error) {
	if payload.Message == "" && payload.Img == "" {
		utils.Response(c, http.StatusBadRequest, 400, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", "err message and img is empty", nil)
		return nil, errors.New("err message and img is empty")
	}

	messageType := ""
	var aiResult *domain.AIPredictionResult
	var image string
	var message string

	if payload.Img != "" {
		messageType = "image"
		image = payload.Img

		if s.aiService != nil {
			var err error
			aiResult, err = s.aiService.PredictSlip(c, payload.Img)
			if err != nil {
				fmt.Println("error calling AI service:", err)
				aiResult = nil
			}
		}
	}

	if payload.Message != "" {
		if messageType == "" {
			messageType = "text"
		}
		message = payload.Message
	}

	roomID, _ := primitive.ObjectIDFromHex(payload.RoomID)

	chat := domain.Chat{
		ID:          primitive.NewObjectID(),
		RoomID:      roomID,
		Type:        "reply",
		MessageType: messageType,
		Message:     message,
		Image:       image,
		AIResult:    aiResult,
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
	}

	if err := s.chatRepo.CreateChat(chat); err != nil {
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, err
	}

	// อัปเดตข้อความล่าสุดเป็นข้อความหรือ placeholder
	lastMessage := "ข้อความใหม่"
	if chat.Message != "" {
		lastMessage = chat.Message
	}
	if err := s.roomRepo.UpdateLastMessage(roomID, lastMessage); err != nil {
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, err
	}

	return aiResult, nil
}
