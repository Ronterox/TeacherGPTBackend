package main

import (
	"context"
	"fmt"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

type Question struct {
    Topic string `json:"topic"`
    Content string `json:"content"`
    Options []string `json:"options"`
    Answer int `json:"answer"`
}

func gpt(message string) (string, error) {
    fmt.Println("Getting GPT3 API Key from .env file...")
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing API KEY")
	}

    fmt.Println("Generating GPT3 completion...")
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
            {Role: "system", Content: "Return a valid json object with test questions and answers about the presented text."}, 
            {Role: "user", Content: message}},
        Model: GPTModel,
	})
	if err != nil {
        return "", fmt.Errorf("ChatCompletion: %v", err)
	}

    usage := resp.Usage
    fmt.Printf("%v API Usage:\n", resp.Model)
    fmt.Printf("PromptTokens: %v\n", usage.PromptTokens)
    fmt.Printf("CompletionTokens: %v\n", usage.CompletionTokens)
    fmt.Printf("TotalTokens: %v\n", usage.TotalTokens)

	return resp.Choices[0].Message.Content, nil
}
