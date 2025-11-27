package handlers

import (
	"net/http"

	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BusinessHandler struct {
	businessService *services.BusinessService
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// CreateBusiness
func (h *BusinessHandler) CreateBusiness(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req models.CreateBusinessRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()
	business, err := h.businessService.CreateBusiness(ctx, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, business)
}

// GetUserBusinesses ดึงจาก user_id
func (h *BusinessHandler) GetUserBusinesses(c *gin.Context) {
	userIDStr := c.GetString("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx := c.Request.Context()
	businesses, err := h.businessService.GetUserBusinesses(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, businesses)
}

// GetBusiness ดึงจาก business_code
func (h *BusinessHandler) GetBusiness(c *gin.Context) {
	businessCode := c.Param("business_code")
	if businessCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Business Code required"})
		return
	}

	ctx := c.Request.Context()
	business, err := h.businessService.GetBusinessByCode(ctx, businessCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, business)
}
