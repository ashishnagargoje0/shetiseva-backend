
package ai

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/go-resty/resty/v2"
)

// GetAIReply calls OpenAI's API with the user's message and returns the AI's reply
func GetAIReply(prompt string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", errors.New("missing OPENAI_API_KEY environment variable")
	}

	client := resty.New()

	requestBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a helpful assistant for Indian farmers. Give useful, local, and practical advice about crops, agriculture, weather, and rural farming.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		return "", err
	}

	if resp.StatusCode() != 200 {
		return "", errors.New("OpenAI API error: " + resp.String())
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", errors.New("no AI response returned")
	}

	return result.Choices[0].Message.Content, nil
}
