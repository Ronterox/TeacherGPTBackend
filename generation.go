package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type Text string

const TOKEN_LIMIT = 4096

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
	log.Println("Checking for directory:", outPath)
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

    if tokens, _ := file.tokenize(); tokens > TOKEN_LIMIT {
        CHUNKS := tokens / TOKEN_LIMIT
        var examResult []Question
        for i := range CHUNKS {
            log.Println("Generating exam from file:", fileName + "_" + strconv.Itoa(i) + "... current", i, "of", CHUNKS)

            chunkFile := file[i*TOKEN_LIMIT:(i+1)*TOKEN_LIMIT]
            chunkName := fileName + "_" + strconv.Itoa(i)

            chunkPath := "outputs/" + chunkName + ".chunk"
            if _, err := os.Stat(chunkPath); err != nil {
                chunkFile.save(chunkPath)
            }

            bytes, err := generateExam(chunkFile, chunkName)
            if err != nil {
                return nil, err
            }

            var exam []Question
            json.Unmarshal(bytes, &exam)

            examResult = append(examResult, exam...)
        }

        return json.MarshalIndent(examResult, "", "    ")
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
