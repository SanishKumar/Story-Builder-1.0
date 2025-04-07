package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	apiURL = "https://api-inference.huggingface.co/models/your-model-name"
)

// AIResponse represents the expected structure of the response.
type AIResponse struct {
	GeneratedText string `json:"generated_text"`
}

// GenerateText calls the Hugging Face API to generate text.
func GenerateText(prompt string) (string, error) {
	apiKey := os.Getenv("HF_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("HF_API_KEY environment variable not set")
	}

	payload := map[string]string{"inputs": prompt}
	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result []AIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result) > 0 {
		return result[0].GeneratedText, nil
	}

	return "", fmt.Errorf("no generated text received")
}
