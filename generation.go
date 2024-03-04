package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math"
	"reflect"
)

type Text string

const TOKEN_LIMIT = 2048

func (text Text) printPriceApprox(tokensValue float32) {
	log.Println("Tokenizing text...")
	tokens := text.tokenize()
	log.Println("Tokens:", tokens)
	log.Println("Expected price: ", float32(tokens)*tokensValue, "USD")
}

func generateMermaidInkUrl(mermaid string) string {
	return "https://mermaid.ink/img/" + base64.URLEncoding.EncodeToString([]byte(mermaid))
}

func getChunk(file Text, i int) Text {
	chunkLimit := math.Min(float64((i+1)*TOKEN_LIMIT*3), float64(len(file)))
	chunkStart := math.Min(float64(i*TOKEN_LIMIT*3), chunkLimit)
	return file[int(chunkStart):int(chunkLimit)]
}

func appendExamBytes[T IQuestion](examResult []T, bytes []byte) ([]T, error) {
	var exam []T
	err := json.Unmarshal(bytes, &exam)
	if err != nil {
		return nil, err
	}
	return append(examResult, exam...), nil
}

func generateExtraQuestions[T IQuestion](file Text, fileName string, numberOfQuestions int, examResult []T) ([]byte, error) {
	for i := range numberOfQuestions {
		log.Printf("Generating extra questions for file: %s_%d... got %d but %d are left\n", fileName, i, i, numberOfQuestions)

		var chunk Text
		var getPrompt func(string) string
		var template string

		if reflect.TypeOf(examResult[i]) == reflect.TypeOf(QuestionOpen{}) {
			template = getQuestionOpenTemplate()
			getPrompt = getPromptQuestionsOpen
		} else {
			template = getQuestionTemplate()
			getPrompt = getPromptQuestions
		}

		if i < len(examResult)-1 {
			template = getTemplate([]IQuestion{examResult[i]})
			chunk = Text(examResult[i].GetQuestion().Chunk)
		} else {
			chunk = getChunk(file, i)
		}

		if len(chunk) < 32 {
			continue
		}

		bytes, err := generateExam[T](chunk, fileName, 1, getPrompt(template))
		if err != nil {
			return nil, err
		}

		examResult, err = appendExamBytes(examResult, bytes)
	}

	return json.MarshalIndent(examResult, "", "    ")
}

func generateExamsList[T IQuestion](file Text, fileName string, numberOfQuestions int, tokens int, prompt string) ([]byte, error) {
	log.Println("The whole file has", tokens, "tokens")
	var examResult []T

	if tokens > TOKEN_LIMIT {
		CHUNKS := int(math.Ceil(float64(tokens) / TOKEN_LIMIT))
		for i := range CHUNKS {
			log.Printf("Generating exam from file: %s... current %d of %d\n", fileName, i, CHUNKS)

			chunkText := getChunk(file, i)
			if len(chunkText) < 32 {
				continue
			}

			bytes, err := generateExam[T](chunkText, fileName, 1, prompt)
			if err != nil {
				return nil, err
			}

			examResult, err = appendExamBytes(examResult, bytes)
			numberOfQuestions--
		}
	}

	if numberOfQuestions > 0 {
		return generateExtraQuestions(file, fileName, numberOfQuestions, examResult)
	}

	return json.MarshalIndent(examResult, "", "    ")
}

// Recursive function to generate exams
func generateExam[T IQuestion](file Text, fileName string, numberOfQuestions int, prompt string) ([]byte, error) {
	if tokens := file.tokenize(); tokens > TOKEN_LIMIT || numberOfQuestions > 1 {
		// This should only be called the first time
		return generateExamsList[T](file, fileName, numberOfQuestions, tokens, prompt)
	}

	file.printPriceApprox(GPTInputPrice)

	filterPrompt := `Make sure to write the questions in Spanish.
    Make sure to not repeat the same question as the example one.
    If you aren't able to generate a question with the given text return an empty array.`

	completion, err := gptQuestions(file, prompt, filterPrompt)
	if err != nil {
		return nil, err
	}

	log.Println("Exam generated", completion)
	Text(completion).printPriceApprox(GPTOutputPrice)
	return []byte(completion), nil
}
