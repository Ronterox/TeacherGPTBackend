package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

const GPTModel = "gpt-3.5-turbo-0125"
const GPTPrice = 0.0005 * 0.001

type Question struct {
	Topic   string   `json:"topic"`
	Content string   `json:"content"`
	Options []string `json:"options"`
	Answer  int      `json:"answer"`
}

func getJsonTemplate() (string, error) {
	// Make this an array
	bytes, err := json.MarshalIndent([]Question{
		Question{
			Topic:   "Computadoras",
			Content: "¿Qué es la memoria RAM?",
			Options: []string{"Memoria de solo lectura", "Memoria de acceso aleatorio", "Memoria de solo escritura", "Memoria de acceso secuencial"},
			Answer:  2,
		},
		Question{
			Topic:   "...",
			Content: "...",
			Options: []string{"...", "...", "...", "..."},
			Answer:  0,
		},
	}, "", "    ")

	if err != nil {
		return "", fmt.Errorf("MarshalIndent: %v", err)
	}

	return string(bytes), nil
}

func gpt(message string) (string, error) {
	fmt.Printf("Getting %v API Key from .env file...\n", GPTModel)
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing API KEY")
	}

	fmt.Printf("Generating %v completion...\n", GPTModel)
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	jsonTemplate, err := getJsonTemplate()
	if err != nil {
		return "", fmt.Errorf("getJsonTemplate: %v", err)
	}

	prompt := `Return a valid json object with test questions and answers about the presented text. 
    The scheme should follow the following example:\n%v`
    filterPrompt := `Make sure to write the questions and answers in Spanish.
    If you aren't able to generate a question with the given text return an empty array.`

	start := time.Now()
	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{Role: "system", Content: fmt.Sprintf(prompt, jsonTemplate)},
			{Role: "system", Content: filterPrompt},
			{Role: "user", Content: message}},
		Model: GPTModel,
	})
	fmt.Printf("%v API call took: %v\n", GPTModel, time.Since(start))
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
