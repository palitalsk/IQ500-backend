package services

import (
	"fmt"
	"strings"
	"time"

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
	fmt.Printf("[ChatService] Processing chat request (business_code: %s, query: %s)\n", businessCode, req.Query)
	startTime := time.Now()

	// เอาไว้ระบุจำนวนคำตอบที่ใกล้เคียง
	if req.TopK == 0 {
		req.TopK = 3
	}

	// ทำ embeddings สำหรับคำถามทีไได้รับมาด้วย
	fmt.Printf("[ChatService] Step 1/3: Getting embeddings...\n")
	embs, err := s.embedRepo.GetEmbeddings([]string{req.Query})
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}
	queryVector := embs[0]

	// หาจาก Pinecone โดยใช้ business_code เป็น namespace
	fmt.Printf("[ChatService] Step 2/3: Querying Pinecone...\n")
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
	fmt.Printf("[ChatService] Found %d context chunks (total length: %d)\n", len(contextTexts), len(context))

	// ส่ง context กับ query ไปให้ AI สร้างคำตอบ
	gender := req.Gender
	if gender == "" {
		gender = "female" // default เป็น female
	}
	fmt.Printf("[ChatService] Step 3/3: Generating answer with Gemini...\n")
	answer, err := s.embedRepo.GenerateAnswer(context, req.Query, gender)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("[ChatService] Completed in %v\n", elapsed)

	return &models.ChatResponse{
		Results: answer,
	}, nil
}
