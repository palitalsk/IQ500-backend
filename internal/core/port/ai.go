package port

import (
	"main/internal/core/domain"

	"github.com/gin-gonic/gin"
)

type AIService interface {
	PredictSlip(c *gin.Context, imageInput string) (*domain.AIPredictionResult, error)
}
