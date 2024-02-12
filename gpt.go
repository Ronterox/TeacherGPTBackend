package main

import (
	"context"
	"fmt"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func gpt(message string) (string, error) {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
            {Role: "system", Content: "You are a helpful assistant."}, 
            {Role: "user", Content: message}},
	})
	if err != nil {
        return "", fmt.Errorf("ChatCompletion: %v", err)
	}
	return resp.Choices[0].Message.Content, nil
}
