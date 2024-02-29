package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"reflect"
	"strconv"
)

type Text string

const TOKEN_LIMIT = 512

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

func generateJson(file Text, open bool) (Text, error) {
	completion, err := gptQuestions(file, open)
	if err != nil {
		return "", err
	}
	return Text(completion), nil
}

func generateMermaidInkUrl(mermaid string) string {
	return "https://mermaid.ink/img/" + base64.URLEncoding.EncodeToString([]byte(mermaid))
}

func generateExam[T QuestionSimple | QuestionOpen](file Text, fileName string) ([]byte, error) {
	log.Println("Generating exam from file:", fileName)

	if tokens, _ := file.tokenize(); tokens > TOKEN_LIMIT {
		CHUNKS := tokens / TOKEN_LIMIT
		var examResult []T
		for i := range CHUNKS {
			log.Printf("Generating exam from file: %s_%d... current %d of %d\n", fileName, i, i, CHUNKS)

			chunkText := file[i*TOKEN_LIMIT : (i+1)*TOKEN_LIMIT]
			chunkName := fileName + "_" + strconv.Itoa(i)

			bytes, err := generateExam[T](chunkText, chunkName)
			if err != nil {
				return nil, err
			}

			var exam []T
			json.Unmarshal(bytes, &exam)

			// for _, question := range exam {
			// }

			examResult = append(examResult, exam...)
		}

		return json.MarshalIndent(examResult, "", "    ")
	}

	if err := file.printPriceApprox(GPTInputPrice); err != nil {
		return nil, err
	}

	var t T
	completion, err := generateJson(file, reflect.TypeOf(t) == reflect.TypeOf(QuestionOpen{}))
	if err != nil {
		return nil, err
	}

	log.Println("Exam generated", completion)

	completion.printPriceApprox(GPTOutputPrice)
	return []byte(completion), nil
}
