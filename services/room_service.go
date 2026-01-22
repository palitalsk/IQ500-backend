package services

import (
	"fmt"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoomService ให้บริการจัดการ room เดิมของ IQ500
type RoomService struct {
	roomRepo port.RoomRepository
}

func NewRoomService(
	roomRepo port.RoomRepository,
) *RoomService {
	return &RoomService{
		roomRepo,
	}
}

func (s *RoomService) GetRoomByUserId(c *gin.Context, id primitive.ObjectID) ([]domain.Room, error) {
	room, err := s.roomRepo.GetRoomByUserId(id)
	if err != nil {
		fmt.Println("error get room", err)
		utils.Response(c, http.StatusInternalServerError, 500, "ไม่พบห้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, err
	}

	return room, nil
}

func (s *RoomService) CreateRoom(c *gin.Context, payload domain.PayloadRoom) (primitive.ObjectID, error) {
	userId, _ := primitive.ObjectIDFromHex(c.GetString("user_id"))

	id := primitive.NewObjectID()

	room := domain.Room{
		ID:          id,
		UserID:      userId,
		LastMessage: payload.Message,
		UpdateAt:    time.Now(),
		CreateAt:    time.Now(),
	}

	if err := s.roomRepo.CreateRoom(room); err != nil {
		fmt.Println("error create room", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return primitive.NilObjectID, err
	}

	return id, nil
}

func (s *RoomService) DeleteRoomByID(c *gin.Context, id primitive.ObjectID) error {
	if err := s.roomRepo.DeleteRoomByID(id); err != nil {
		fmt.Println("error delete room", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return err
	}

	return nil
}
