package services

import (
	"errors"
	"fmt"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoomChatService ให้บริการ chat ภายใน room เดิมของ IQ500
type RoomChatService struct {
	chatRepo port.ChatRepository
	roomRepo port.RoomRepository
}

func NewRoomChatService(
	chatRepo port.ChatRepository,
	roomRepo port.RoomRepository,
) *RoomChatService {
	return &RoomChatService{
		chatRepo,
		roomRepo,
	}
}

func (s *RoomChatService) GetChatByRoomID(c *gin.Context, id primitive.ObjectID) ([]domain.Chat, error) {
	chat, err := s.chatRepo.GetChatByRoomID(id)
	if err != nil {
		fmt.Println("error get chat", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, err
	}

	return chat, nil
}

func (s *RoomChatService) Chat(c *gin.Context, payload domain.PayloadChat) error {
	if payload.Message == "" && payload.Img == "" {
		fmt.Println("err message and img is empty")
		utils.Response(c, http.StatusBadRequest, 400, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", "err message and img is empty", nil)
		return errors.New("err message and img is empty")
	}

	messageType := ""
	message := ""
	if payload.Message != "" {
		messageType = "text"
		message = payload.Message
	} else {
		messageType = "image"
		message = payload.Img
	}

	roomID, _ := primitive.ObjectIDFromHex(payload.RoomID)

	chat := domain.Chat{
		ID:          primitive.NewObjectID(),
		RoomID:      roomID,
		Type:        "reply",
		MessageType: messageType,
		Message:     message,
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
	}

	err := s.chatRepo.CreateChat(chat)
	if err != nil {
		fmt.Println("error create chat", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return err
	}

	lastMessage := ""
	if chat.MessageType == "text" {
		lastMessage = chat.Message
	} else {
		lastMessage = "ส่งรูปภาพ"
	}
	err = s.roomRepo.UpdateLastMessage(roomID, lastMessage)
	if err != nil {
		fmt.Println("error update last message", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return err
	}

	return nil
}


