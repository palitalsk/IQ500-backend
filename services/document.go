package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

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

// ProcessUploadedFile >> อัปโหลดไฟล์ และเก็บลง pinecone
func (s *DocumentService) ProcessUploadedFile(file io.Reader, filename, businessCode string) (*models.DocumentUploadResponse, error) {
	// Save to temp file
	tmp, err := os.CreateTemp("", "upload-*"+filepath.Ext(filename))
	if err != nil {
		return nil, fmt.Errorf("cannot create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	_, err = io.Copy(tmp, file)
	if err != nil {
		return nil, fmt.Errorf("cannot save uploaded file: %w", err)
	}
	tmp.Close()

	// Extract text
	var text string
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == ".pdf" {
		text, err = helpers.ReadPdfText(tmp.Name())
		if err != nil {
			return nil, fmt.Errorf("cannot read PDF: %w", err)
		}
	} else {
		// .txt or plain text
		b, err := os.ReadFile(tmp.Name())
		if err != nil {
			return nil, fmt.Errorf("cannot read file: %w", err)
		}
		text = string(b)
	}

	// Clean text
	text = helpers.CleanText(text)
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("no text extracted from file")
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
		id := fmt.Sprintf("%s-%s-%d", businessCode, filename, i)
		metadata := map[string]interface{}{
			"text":          chunks[i],
			"source":        filename,
			"business_code": businessCode,
			"chunk_index":   i,
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
		Filename:  filename,
		Chunks:    len(chunks),
		Status:    "completed",
		CreatedAt: time.Now(),
	}, nil
}
