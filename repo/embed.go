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

type EmbedRepository struct {
	baseURL string
}

func NewEmbedRepository(baseURL string) *EmbedRepository {
	return &EmbedRepository{
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

// GetEmbeddings >> ส่งข้อความไป embed-service เพื่อให้สร้าง embed
func (r *EmbedRepository) GetEmbeddings(texts []string) ([][]float64, error) {
	req := models.EmbedRequest{Texts: texts}
	b, _ := json.Marshal(req)
	url := r.baseURL + "/embed"

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("embed service status %d: %s", resp.StatusCode, string(body))
	}

	var er models.EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, err
	}
	return er.Embeddings, nil
}

// GenerateAnswer >> ส่ง context กับ query ไป embed-service เพื่อให้ AI สร้างคำตอบ
func (r *EmbedRepository) GenerateAnswer(context, query, gender string) (string, error) {
	reqBody := models.GenerateRequest{
		Context: context,
		Query:   query,
		Gender:  gender,
	}
	b, _ := json.Marshal(reqBody)

	url := r.baseURL + "/generate"
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("AI service error: %s", string(body))
	}

	var gr models.GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return "", err
	}

	return gr.Answer, nil
}
