package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"main/models"
)

type PineconeRepository struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

func NewPineconeRepository(baseURL, apiKey string) *PineconeRepository {
	return &PineconeRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// UpsertVectors >> upserts vectors to Pinecone
func (r *PineconeRepository) UpsertVectors(vectors []models.PineconeVector, namespace string) error {
	url := r.baseURL + "/vectors/upsert"
	reqBody := models.PineconeUpsertRequest{Vectors: vectors, Namespace: namespace}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("pinecone upsert error: %s", string(body))
	}
	return nil
}

// QueryVectors >> queries vectors from Pinecone
func (r *PineconeRepository) QueryVectors(vector []float64, topK int, namespace string) (*models.PineconeQueryResponse, error) {
	fmt.Printf("[Pinecone] Querying vectors (topK: %d, namespace: %s)\n", topK, namespace)
	startTime := time.Now()

	url := r.baseURL + "/query"
	reqBody := models.PineconeQueryRequest{
		Vector:          vector, //คำถามที่ี embed แล้ว
		TopK:            topK,
		IncludeMetadata: true,
		Namespace:       namespace, // business ID
	}
	b, _ := json.Marshal(reqBody)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		elapsed := time.Since(startTime)
		return nil, fmt.Errorf("failed to query Pinecone (after %v): %w", elapsed, err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(startTime)
	fmt.Printf("[Pinecone] Received response in %v (status: %d)\n", elapsed, resp.StatusCode)

	var pr models.PineconeQueryResponse
	if resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("pinecone query error: %s", string(body))
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	fmt.Printf("[Pinecone] Found %d matches\n", len(pr.Matches))
	return &pr, nil
}
