package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Text string

const TOKEN_LIMIT = 512

func (text Text) save(outPath string) error {
	log.Println("Saving file:", outPath)
	return os.WriteFile(outPath, []byte(text), 0644)
}

func (text Text) printPriceApprox(tokensValue float32) error {
	log.Println("Tokenizing text...")
	tokens, err := text.tokenize()
	if err != nil {
		return err
	}
	log.Println("Tokens:", tokens)
	log.Println("Expected price: ", float32(tokens)*tokensValue, "USD")
	return nil
}

func generateDir(outPath string) error {
	log.Println("Checking for directory:", outPath)
	if err := os.Mkdir(outPath, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func generateJsonAndSave(file Text, outPath string, open bool) (Text, error) {
	completion, err := gptQuestions(file, open)
	if err != nil {
		return "", err
	}
	completionText := Text(completion)
	return completionText, completionText.save(outPath)
}

func generateMermaidInkUrl(mermaid string) string {
	return "https://mermaid.ink/img/" + base64.URLEncoding.EncodeToString([]byte(mermaid))
}

func generateExam[T Question | QuestionOpen](file Text, fileName string, open bool) ([]byte, error) {
	log.Println("Generating exam from file:", fileName)

	if tokens, _ := file.tokenize(); tokens > TOKEN_LIMIT {
		CHUNKS := tokens / TOKEN_LIMIT
		var examResult []T
		for i := range CHUNKS {
			log.Printf("Generating exam from file: %s_%d... current %d of %d\n", fileName, i, i, CHUNKS)

			chunkFile := file[i*TOKEN_LIMIT : (i+1)*TOKEN_LIMIT]
			chunkName := fileName + "_" + strconv.Itoa(i)

			chunkPath := "outputs/" + chunkName + ".chunk"
			if _, err := os.Stat(chunkPath); err != nil {
				chunkFile.save(chunkPath)
			}

			bytes, err := generateExam[T](chunkFile, chunkName, open)
			if err != nil {
				return nil, err
			}

			var exam []T
			json.Unmarshal(bytes, &exam)

			examResult = append(examResult, exam...)
		}

		return json.MarshalIndent(examResult, "", "    ")
	}

	if err := file.printPriceApprox(GPTInputPrice); err != nil {
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

	completion, err := generateJsonAndSave(file, outPath, open)
	if err != nil {
		return nil, err
	}

	log.Println("Exam generated and saved to " + outPath)
	log.Println(completion)

	completion.printPriceApprox(GPTOutputPrice)
	return []byte(completion), nil
}
