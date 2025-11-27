package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"main/models"
)

type PineconeRepository struct {
	baseURL string
	apiKey  string
}

func NewPineconeRepository(baseURL, apiKey string) *PineconeRepository {
	return &PineconeRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
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
	url := r.baseURL + "/query"
	reqBody := models.PineconeQueryRequest{
		Vector:          vector, //คำถามที่ี embed แล้ว
		TopK:            topK,
		IncludeMetadata: true,
		Namespace:       namespace, // business ID
	}
	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Api-Key", r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pr models.PineconeQueryResponse
	if resp.StatusCode >= 300 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("pinecone query error: %s", string(body))
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	return &pr, nil
}
