package models

import "time"

type Document struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"business_id"`
	Filename   string    `json:"filename"`
	FileType   string    `json:"file_type"`
	FileSize   int64     `json:"file_size"`
	Chunks     int       `json:"chunks"`
	Status     string    `json:"status"` // processing, completed, failed
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentUploadResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Chunks    int       `json:"chunks"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
