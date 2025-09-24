package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/internal/core/domain"
	"main/internal/core/port"
	"main/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AIService struct {
	baseURL string
	timeout time.Duration
}

func NewAIService(baseURL string, timeout int) port.AIService {
	return &AIService{
		baseURL: baseURL,
		timeout: time.Duration(timeout) * time.Second,
	}
}

func (s *AIService) PredictSlip(c *gin.Context, imageInput string) (*domain.AIPredictionResult, error) {
	// request payload
	payload := map[string]string{
		"image": imageInput,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marshal request", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{
		Timeout: s.timeout,
	}

	// Make request to AI service
	resp, err := client.Post(s.baseURL+"/predict-slip", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error calling AI service", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, fmt.Errorf("failed to call AI service: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("AI service error", resp.StatusCode, string(body))
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", fmt.Sprintf("AI service returned status %d", resp.StatusCode), nil)
		return nil, fmt.Errorf("AI service returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result domain.AIPredictionResult
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("error unmarshal response", err)
		utils.Response(c, http.StatusInternalServerError, 500, "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง", err.Error(), nil)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Convert to domain result
	domainResult := &domain.AIPredictionResult{
		IsSlip:    result.IsSlip,
		OCRResult: result.OCRResult,
	}

	return domainResult, nil
}
