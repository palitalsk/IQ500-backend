package handlers

import (
	"fmt"
	"main/domain"
	"main/port"
	"main/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoomChatHandler จัดการ endpoint chat ภายใน room เดิมของ IQ500
type RoomChatHandler struct {
	svc port.ChatService
}

func NewRoomChatHandler(svc port.ChatService) *RoomChatHandler {
	return &RoomChatHandler{
		svc: svc,
	}
}

func (h *RoomChatHandler) GetChatByRoomID(c *gin.Context) {
	fmt.Println("GetChatByRoomID")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	chat, err := h.svc.GetChatByRoomID(c, id)
	if err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", chat)
}

func (h *RoomChatHandler) Chat(c *gin.Context) {
	fmt.Println("Chat")

	var payload domain.PayloadChat
	if err := c.ShouldBind(&payload); err != nil {
		fmt.Println("error bind", err)
		utils.Response(c, http.StatusBadRequest, 1, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return
	}

	if err := h.svc.Chat(c, payload); err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", nil)
}


