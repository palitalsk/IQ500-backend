package handlers

import (
	"net/http"

	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService *services.ChatService
}

func NewChatHandler(chatService *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// HandleChat
func (h *ChatHandler) HandleChat(c *gin.Context) {
	businessCode := c.GetString("business_code") // ดึง business_code มาเช็ค
	if businessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business Code required"})
		return
	}

	var req models.ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	// Process chat request >> ส่งข้อความไปหา AI
	response, err := h.chatService.ProcessChat(businessCode, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
