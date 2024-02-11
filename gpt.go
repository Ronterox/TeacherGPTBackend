package main

import (
	"context"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

func gpt(message string) string {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo,
		Messages: []gpt3.ChatCompletionRequestMessage{
            {Role: "system", Content: "You are a helpful assistant."}, 
            {Role: "user", Content: message}},
	})
	if err != nil {
		log.Fatalln(err)
	}
	return resp.Choices[0].Message.Content
}
