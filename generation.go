package main

import (
	"fmt"
	"log"
	"os"
)

type Text string

const TOKEN_LIMIT = 4096
const TOTAL_TOKEN_LIMIT = 4096 * 3
var tokenCount int

func (text Text) save(outPath string) error {
	log.Println("Saving file:", outPath)
	return os.WriteFile(outPath, []byte(text), 0644)
}

func (text Text) printPriceApprox() error {
	log.Println("Tokenizing text...")
	tokens, err := text.tokenize()
	if err != nil {
		return err
	}
	log.Println("Tokens:", tokens)
	log.Println("Expected price: ", float32(tokens)*GPTPrice, "USD")
	return nil
}

func generateDir(outPath string) error {
	log.Println("Creating directory:", outPath)
	if err := os.Mkdir(outPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func generateJsonAndSave(file Text, outPath string) (Text, error) {
	completion, err := gpt(string(file))
	if err != nil {
		return "", err
	}
	completionText := Text(completion)
	return completionText, completionText.save(outPath)
}

func generateExam(file Text, fileName string) ([]byte, error) {
	log.Println("Generating exam from file:", fileName)

    if tokenCount > TOTAL_TOKEN_LIMIT {
        return nil, fmt.Errorf("Token limit exceeded")
    }

    if tokens, _ := file.tokenize(); tokens > TOKEN_LIMIT {
        tokenCount += TOKEN_LIMIT
        return generateExam(file[:TOKEN_LIMIT], fileName)
    }

	if err := file.printPriceApprox(); err != nil {
		return nil, err
	}

	if err := generateDir("outputs"); err != nil {
		return nil, err
	}

	outPath := "outputs/" + fileName + ".json"
	if _, err := os.Stat(outPath); err == nil {
        log.Println("Using cached file:", outPath)
		return os.ReadFile(outPath)
	}

	completion, err := generateJsonAndSave(file, outPath)
	if err != nil {
		return nil, err
	}

	log.Println("Exam generated and saved to " + outPath)
	log.Println(completion)

	completion.printPriceApprox()
	return []byte(completion), nil
}
