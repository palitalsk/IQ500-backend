package services

import (
	"fmt"
	"strings"

	"main/models"
	"main/repo"
)

type ChatService struct {
	embedRepo    *repo.EmbedRepository
	pineconeRepo *repo.PineconeRepository
}

func NewChatService(embedRepo *repo.EmbedRepository, pineconeRepo *repo.PineconeRepository) *ChatService {
	return &ChatService{
		embedRepo:    embedRepo,
		pineconeRepo: pineconeRepo,
	}
}

// ProcessChat
func (s *ChatService) ProcessChat(businessCode string, req *models.ChatRequest) (*models.ChatResponse, error) {
	// เอาไว้ระบุจำนวนคำตอบที่ใกล้เคียง
	if req.TopK == 0 {
		req.TopK = 3
	}

	// ทำ embeddings สำหรับคำถามทีไได้รับมาด้วย
	embs, err := s.embedRepo.GetEmbeddings([]string{req.Query})
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	queryVector := embs[0]

	// หาจาก Pinecone โดยใช้ business_code เป็น namespace
	matches, err := s.pineconeRepo.QueryVectors(queryVector, req.TopK, businessCode)
	if err != nil {
		return nil, fmt.Errorf("pinecone query failed: %w", err)
	}

	// รวมข้อความที่ได้จาก Pinecone เป็น context (context เป็นข้อมูลที่เกี่ยวข้องกับคำถามที่จะส่งไปให้ AI)
	var contextTexts []string
	for _, match := range matches.Matches {
		if t, ok := match.Metadata["text"].(string); ok {
			contextTexts = append(contextTexts, t)
		}
	}
	context := strings.Join(contextTexts, "\n")

	// ส่ง context กับ query ไปให้ AI สร้างคำตอบ
	gender := req.Gender
	if gender == "" {
		gender = "female" // default เป็น female
	}
	answer, err := s.embedRepo.GenerateAnswer(context, req.Query, gender)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	return &models.ChatResponse{
		Query:     req.Query, // คำถาม
		Namespace: businessCode,
		Results:   answer,
	}, nil
}
