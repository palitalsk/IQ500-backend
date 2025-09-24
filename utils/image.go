package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func IsBase64(s string) bool {
	if strings.HasPrefix(s, "data:image/") && strings.Contains(s, ";base64,") {
		return true
	}

	// decode as base64
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func IsURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func URLToBase64(url string) (string, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: status %d", resp.StatusCode)
	}

	// Read the response body
	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image data: %w", err)
	}

	// Convert to base64
	base64String := base64.StdEncoding.EncodeToString(imageData)
	return base64String, nil
}

func FileToBase64(fileData []byte) string {
	return base64.StdEncoding.EncodeToString(fileData)
}

func IsValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
		"image/bmp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

func ProcessImageInput(input string) (string, error) {
	if input == "" {
		return "", fmt.Errorf("image input is empty")
	}

	if IsBase64(input) {
		if strings.HasPrefix(input, "data:image/") && strings.Contains(input, ";base64,") {
			parts := strings.Split(input, ";base64,")
			if len(parts) == 2 {
				return parts[1], nil
			}
		}
		return input, nil
	}

	if IsURL(input) {
		return URLToBase64(input)
	}

	return "", fmt.Errorf("invalid image input: must be base64 or URL")
}

func ProcessImageFile(fileData []byte, contentType string) (string, error) {
	if len(fileData) == 0 {
		return "", fmt.Errorf("file is empty")
	}

	if !IsValidImageType(contentType) {
		return "", fmt.Errorf("invalid image type: %s. Supported types: jpeg, jpg, png, gif, webp, bmp", contentType)
	}

	// Convert to base64
	return FileToBase64(fileData), nil
}
