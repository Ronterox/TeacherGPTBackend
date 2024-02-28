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
const GPTInputPrice = 0.0005 * 0.001
const GPTOutputPrice = 0.0015 * 0.001

type Question struct {
	Topic   string `json:"topic"`
	Content string `json:"content"`
	Chunk   string `json:"chunk"`
}

type QuestionSimple struct {
	Question
	Options []string `json:"options"`
	Answer  int      `json:"answer"`
}

type QuestionOpen struct {
	Question
	Answer  string `json:"answer"`
	Correct bool   `json:"correct"`
	Reason  string `json:"reason"`
}

func getJsonTemplate() (string, error) {
	list := []QuestionSimple{
		QuestionSimple{
			Question: Question{
				Topic:   "Computadoras",
				Content: "¿Qué es la memoria RAM?",
				Chunk:   "La memoria RAM es una memoria de acceso aleatorio que se utiliza para almacenar datos e instrucciones. Es una memoria volátil, lo que significa que los datos se pierden cuando se apaga la computadora. La memoria RAM es más rápida que la memoria de almacenamiento a largo plazo, como los discos duros y las unidades de estado sólido, pero también es más cara y tiene una capacidad de almacenamiento más limitada.",
			},
			Options: []string{"Memoria de solo lectura", "Memoria de acceso aleatorio", "Memoria de solo escritura", "Memoria de acceso secuencial"},
			Answer:  1,
		},
		QuestionSimple{
			Question: Question{
				Topic:   "...",
				Content: "...",
				Chunk:   "...",
			},
			Options: []string{"...", "...", "...", "..."},
			Answer:  0,
		},
	}
	return getTemplate(list)
}

func getJsonTemplateOpen() (string, error) {
	list := []QuestionOpen{
		QuestionOpen{
			Question: Question{
				Topic:   "Computadoras",
				Content: "¿Qué es la memoria RAM?",
				Chunk:   "La memoria RAM es una memoria de acceso aleatorio que se utiliza para almacenar datos e instrucciones. Es una memoria volátil, lo que significa que los datos se pierden cuando se apaga la computadora. La memoria RAM es más rápida que la memoria de almacenamiento a largo plazo, como los discos duros y las unidades de estado sólido, pero también es más cara y tiene una capacidad de almacenamiento más limitada.",
			},
		},
		QuestionOpen{
			Question: Question{
				Topic:   "...",
				Content: "...",
				Chunk:   "...",
			},
		},
	}
	return getTemplate(list)
}

func getTemplate[T QuestionSimple | QuestionOpen](questions []T) (string, error) {
	bytes, err := json.MarshalIndent(questions, "", "    ")
	if err != nil {
		return "", fmt.Errorf("json.MarshalIndent: %v", err)
	}
	return string(bytes), nil
}

func gpt(userData string, systemPrompts []gpt3.ChatCompletionRequestMessage) (string, error) {
	fmt.Printf("Getting %v API Key from .env file...\n", GPTModel)
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("Missing API KEY")
	}

	fmt.Printf("Generating %v completion...\n", GPTModel)
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	start := time.Now()
	resp, err := client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Messages: append(systemPrompts, gpt3.ChatCompletionRequestMessage{Role: "user", Content: userData}),
		Model:    GPTModel,
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

func gptQuestions(data Text, open bool) (string, error) {
	var jsonTemplate, prompt, filterPrompt string
	var err error

	if open {
		jsonTemplate, err = getJsonTemplateOpen()
		prompt = `Return a valid json array with test questions about the presented text.
        The scheme should follow the following example:\n%v`
		filterPrompt = `Make sure to write the questions in Spanish.
        If you aren't able to generate a question with the given text return an empty array.`
	} else {
		jsonTemplate, err = getJsonTemplate()
		prompt = `Return a valid json object with test questions and answers about the presented text. 
        The scheme should follow the following example:\n%v`
		filterPrompt = `Make sure to write the questions and answers in Spanish.
        If you aren't able to generate a question with the given text return an empty array.`
	}

	if err != nil {
		return "", fmt.Errorf("getJsonTemplate: %v", err)
	}

	systemPrompts := []gpt3.ChatCompletionRequestMessage{
		{Role: "system", Content: fmt.Sprintf(prompt, jsonTemplate)},
		{Role: "system", Content: filterPrompt}}

	return gpt(string(data), systemPrompts)
}

func gptMindMap(data Text) (string, error) {
	codeExample := `mindmap
    )My Mindmap(
        (Origins)
            [Long history]
            ::icon(fa fa-book)
            (Popularisation)
                [British popular psychology author Tony Buzan]
        (Research)
            [On effectivness<br/>and features]
            [On Automatic creation]
                (Uses)
                    [Creative techniques]
                    [Strategic planning]
                    [Argument mapping]`

	prompt := `Return just a valid Mermaid Mindmap of the presented text.
    This is an example of how a Mindmap should look like:\n%v`
	filterPrompt := `The syntax for creating Mindmaps relies on indentation for setting the levels in the hierarchy.
    Please use the following syntax )For the root(, (For Titles) and [For subtitles].`
	jailPrompt := `Don't ever use parenthesis inside of brackets. You can only use tabs and spaces for indentation. 
    There can only be one root, at idented level 0, and please return the content in Spanish.`

	systemPrompts := []gpt3.ChatCompletionRequestMessage{
		{Role: "system", Content: fmt.Sprintf(prompt, codeExample)},
		{Role: "system", Content: filterPrompt},
		{Role: "system", Content: jailPrompt}}

	return gpt(string(data), systemPrompts)
}
