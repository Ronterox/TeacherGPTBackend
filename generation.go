package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math"
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

func generateMermaidInkUrl(mermaid string) string {
	return "https://mermaid.ink/img/" + base64.URLEncoding.EncodeToString([]byte(mermaid))
}

func generateExam[T QuestionInterface](file Text, fileName string, numberOfQuestions int) ([]byte, error) {
	if tokens, _ := file.tokenize(); tokens > TOKEN_LIMIT || numberOfQuestions > 1 {
        log.Println("The whole file has", tokens, "tokens")
		var examResult []T

        if tokens > TOKEN_LIMIT {
            CHUNKS := int(math.Ceil(float64(tokens) / TOKEN_LIMIT))
            for i := range CHUNKS {
                log.Printf("Generating exam from file: %s... current %d of %d\n", fileName, i, CHUNKS)

                chunkLimit := math.Min(float64((i+1)*TOKEN_LIMIT * 3), float64(len(file)))
                chunkText := file[i*TOKEN_LIMIT * 3 : int(chunkLimit)]
                chunkName := fileName + "_" + strconv.Itoa(i)

                if len(chunkText) < 32 {
                    continue
                }

                bytes, err := generateExam[T](chunkText, chunkName, 1)
                if err != nil {
                    return nil, err
                }

                var exam []T
                json.Unmarshal(bytes, &exam)

                examResult = append(examResult, exam...)
                numberOfQuestions--
            }
        }

        if numberOfQuestions > 0 {
            for i := range numberOfQuestions {
                log.Printf("Generating extra questions for file: %s_%d... got %d but %d are left\n", fileName, i, i, numberOfQuestions)
                chunkLimit := math.Min(float64((i+1)*TOKEN_LIMIT), float64(len(file)))
                chunk := file[i*TOKEN_LIMIT : int(chunkLimit)]

                if i < len(examResult) {
                    chunk = Text(examResult[i].GetQuestion().Chunk)
                }

                if len(chunk) < 32 {
                    continue
                }

                bytes, err := generateExam[T](chunk, fileName, 1)
                if err != nil {
                    return nil, err
                }

                var exam []T
                json.Unmarshal(bytes, &exam)
                examResult = append(examResult, exam...)
            }
        }

		return json.MarshalIndent(examResult, "", "    ")
	}

	if err := file.printPriceApprox(GPTInputPrice); err != nil {
		return nil, err
	}

	var t T
	completion, err := gptQuestions(file, reflect.TypeOf(t) == reflect.TypeOf(QuestionOpen{}))
	if err != nil {
		return nil, err
	}

	log.Println("Exam generated", completion)
	Text(completion).printPriceApprox(GPTOutputPrice)
	return []byte(completion), nil
}
