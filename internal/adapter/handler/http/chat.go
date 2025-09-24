package http

import (
	"fmt"
	"main/internal/core/domain"
	"main/internal/core/port"
	"main/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatHandler struct {
	svc port.ChatService
}

func NewChatHandler(svc port.ChatService) *ChatHandler {
	return &ChatHandler{
		svc,
	}
}

func (h *ChatHandler) GetChatByRoomID(c *gin.Context) {
	fmt.Println("GetChatByRoomID")

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	chat, err := h.svc.GetChatByRoomID(c, id)
	if err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", chat)
}

func (h *ChatHandler) Chat(c *gin.Context) {
	fmt.Println("Chat")

	var payload domain.PayloadChat

	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// รับข้อมูลจาก form data
		roomID := c.PostForm("room_id")
		message := c.PostForm("message")

		if roomID == "" {
			utils.Response(c, http.StatusBadRequest, 1, "room_id is required", "room_id is required", nil)
			return
		}

		var imgData string

		file, err := c.FormFile("image")
		if err == nil {
			// มีไฟล์อัปโหลด
			src, err := file.Open()
			if err != nil {
				fmt.Println("error opening file", err)
				utils.Response(c, http.StatusBadRequest, 1, "ไม่สามารถเปิดไฟล์ได้", err.Error(), nil)
				return
			}
			defer src.Close()

			fileData := make([]byte, file.Size)
			_, err = src.Read(fileData)
			if err != nil {
				fmt.Println("error reading file", err)
				utils.Response(c, http.StatusBadRequest, 1, "ไม่สามารถอ่านไฟล์ได้", err.Error(), nil)
				return
			}

			base64Image, err := utils.ProcessImageFile(fileData, file.Header.Get("Content-Type"))
			if err != nil {
				fmt.Println("error processing file", err)
				utils.Response(c, http.StatusBadRequest, 1, "ไม่สามารถประมวลผลไฟล์ได้", err.Error(), nil)
				return
			}
			imgData = base64Image
		} else {
			// ไม่มีไฟล์ แต่อาจมี URL
			imgData = c.PostForm("img")
		}

		payload = domain.PayloadChat{
			RoomID:  roomID,
			Message: message,
			Img:     imgData,
		}
	} else {
		// รับข้อมูลจาก JSON
		if err := c.ShouldBind(&payload); err != nil {
			fmt.Println("error bind", err)
			utils.Response(c, http.StatusBadRequest, 1, "ข้อมูลไม่ถูกต้อง กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
			return
		}
	}

	if err := h.svc.Chat(c, payload); err != nil {
		return
	}

	utils.Response(c, http.StatusOK, 200, "success", "ok", nil)
}
