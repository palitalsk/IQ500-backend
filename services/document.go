package services

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"main/helpers"
	"main/models"
	"main/repo"
)

type DocumentService struct {
	embedRepo    *repo.EmbedRepository
	pineconeRepo *repo.PineconeRepository
}

func NewDocumentService(embedRepo *repo.EmbedRepository, pineconeRepo *repo.PineconeRepository) *DocumentService {
	return &DocumentService{
		embedRepo:    embedRepo,
		pineconeRepo: pineconeRepo,
	}
}

// ProcessDocumentData >> รับ title และ content แล้วเก็บลง pinecone
func (s *DocumentService) ProcessDocumentData(title, content, businessCode string) (*models.DocumentUploadResponse, error) {
	// รวม title และ content
	var text string
	if title != "" {
		text += fmt.Sprintf("ชื่อเรื่อง: %s\n\n", title)
	}
	if content != "" {
		text += content
	}

	// Clean text
	text = helpers.CleanText(text)
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("title and content cannot be empty")
	}

	// Split into chunks
	chunks := helpers.SplitTextIntoChunks(text, 150) // ~150 words per chunk

	// Get embeddings
	embs, err := s.embedRepo.GetEmbeddings(chunks)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	// สร้าง vectors
	var vectors []models.PineconeVector

	for i, v := range embs {
		// ใช้ ObjectID แทนชื่อเพื่อให้เป็น unique
		id := primitive.NewObjectID().Hex()
		metadata := map[string]interface{}{
			"text":          chunks[i],
			"source":        title, // เก็บ title เดิมใน metadata
			"business_code": businessCode,
			"chunk_index":   i,
		}
		// เพิ่ม title ใน metadata
		if title != "" {
			metadata["title"] = title
		}
		vec := models.PineconeVector{
			ID:       id,
			Values:   v,
			Metadata: metadata,
		}
		vectors = append(vectors, vec)
	}

	// ส่งข้อมูลเข้า Pinecone โดยใช้ business_code เป็น namespace
	err = s.pineconeRepo.UpsertVectors(vectors, businessCode)
	if err != nil {
		return nil, fmt.Errorf("pinecone upsert failed: %w", err)
	}

	documentID := fmt.Sprintf("doc_%d", time.Now().Unix())

	return &models.DocumentUploadResponse{
		ID:        documentID,
		Title:     title,
		Chunks:    len(chunks),
		Status:    "completed",
		CreatedAt: time.Now(),
	}, nil
}
