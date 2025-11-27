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

// HandleUpload
func (h *UploadHandler) HandleUpload(c *gin.Context) {
	businessCode := c.GetString("business_code")
	if businessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business Code required"})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	// Process uploaded file
	response, err := h.documentService.ProcessUploadedFile(file, header.Filename, businessCode) // extract text, embed, upsert ไป pinecone
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
