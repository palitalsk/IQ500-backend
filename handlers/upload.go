package handlers

import (
	"net/http"

	"main/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	documentService *services.DocumentService
}

func NewUploadHandler(documentService *services.DocumentService) *UploadHandler {
	return &UploadHandler{
		documentService: documentService,
	}
}

func (h *UploadHandler) HandleUpload(c *gin.Context) {
	var req struct {
		BusinessCode string `json:"business_code" binding:"required"`
		Title        string `json:"title" binding:"required"`
		Content      string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "business_code, title and content are required"})
		return
	}

	response, err := h.documentService.ProcessDocumentData(req.Title, req.Content, req.BusinessCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
