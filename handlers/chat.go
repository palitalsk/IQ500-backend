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

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var req models.ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	businessCode := req.Namespace
	if businessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_code (namespace) is required in request body"})
		return
	}

	response, err := h.chatService.ProcessChat(businessCode, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

