package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"main/models"
)

type EmbedRepository struct {
	baseURL      string
	client       *http.Client
	geminiApiKey string
	generateURL  string
}

func NewEmbedRepository(baseURL, geminiApiKey string) *EmbedRepository {
	return &EmbedRepository{
		baseURL:      strings.TrimRight(baseURL, "/"),
		geminiApiKey: geminiApiKey,
		generateURL:  "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-exp:generateContent",
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (r *EmbedRepository) GetEmbeddings(texts []string) ([][]float64, error) {
	if r.baseURL == "" {
		return nil, fmt.Errorf("embed service URL is not configured. Please set EMBED_SERVICE_URL environment variable")
	}

	fmt.Printf("[Embed] Requesting embeddings for %d text(s) from %s/embed\n", len(texts), r.baseURL)
	startTime := time.Now()

	req := models.EmbedRequest{Texts: texts}
	b, _ := json.Marshal(req)
	url := r.baseURL + "/embed"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(httpReq)
	if err != nil {
		elapsed := time.Since(startTime)
		return nil, fmt.Errorf("failed to connect to embed service at %s (after %v): %w", url, elapsed, err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(startTime)
	fmt.Printf("[Embed] Received response in %v (status: %d)\n", elapsed, resp.StatusCode)

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("embed service status %d: %s", resp.StatusCode, string(body))
	}

	var er models.EmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return er.Embeddings, nil
}

func (r *EmbedRepository) GenerateAnswer(ctxText, query, gender string) (string, error) {
	if r.geminiApiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not configured. Please set it in your .env file")
	}

	fmt.Printf("[Generate] Calling Gemini API directly (query length: %d, context length: %d)\n", len(query), len(ctxText))
	startTime := time.Now()

	genderSuffix := "ค่ะ"
	notFoundText := "ขออภัย เราไม่พบข้อมูลที่เกี่ยวข้องค่ะ"
	if gender == "male" {
		genderSuffix = "ครับ"
		notFoundText = "ขออภัย เราไม่พบข้อมูลที่เกี่ยวข้องครับ"
	}

	prompt := fmt.Sprintf(`คุณเป็นผู้ช่วย AI ที่ใช้ภาษาสุภาพ อ่อนโยน
ตอบคำถามโดยใช้ภาษาไทยธรรมชาติ ไม่ต้องลงท้ายด้วย %s ทุกประโยค ใช้ให้พอดีตามจังหวะการพูด
ตอบคำถามตามบริบทที่กำหนดเท่านั้น หากไม่พบข้อมูลในบริบท ให้ตอบว่า '%s'

---
บริบท (Context):
%s

---
คำถาม: %s`, genderSuffix, notFoundText, ctxText, query)

	reqBody := struct {
		Contents []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"contents"`
		GenerationConfig struct {
			Temperature     float64 `json:"temperature"`
			MaxOutputTokens int     `json:"maxOutputTokens"`
		} `json:"generationConfig"`
	}{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{
					{Text: prompt},
				},
			},
		},
		GenerationConfig: struct {
			Temperature     float64 `json:"temperature"`
			MaxOutputTokens int     `json:"maxOutputTokens"`
		}{
			Temperature:     0.1,
			MaxOutputTokens: 300,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request failed: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", r.generateURL, r.geminiApiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		elapsed := time.Since(startTime)
		return "", fmt.Errorf("request failed (after %v): %w", elapsed, err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(startTime)
	fmt.Printf("[Generate] Received response in %v (status: %d)\n", elapsed, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	var genResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &genResp); err != nil {
		return "", fmt.Errorf("decode response failed: %w", err)
	}

	if len(genResp.Candidates) > 0 && len(genResp.Candidates[0].Content.Parts) > 0 {
		answer := strings.TrimSpace(genResp.Candidates[0].Content.Parts[0].Text)
		if answer != "" {
			return answer, nil
		}
	}

	return notFoundText, nil
}
