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
            {Role: "system", Content: "Return a valid json object with questions and answers about the presented text."}, 
            {Role: "user", Content: message}},
        Model: "gpt-3.5-turbo-0125",
	})
	if err != nil {
        return "", fmt.Errorf("ChatCompletion: %v", err)
	}
	return resp.Choices[0].Message.Content, nil
}
