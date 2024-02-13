package main

import (
	"log"
	"os"
)

type Text string
const GPTModel = "gpt-3.5-turbo-0125"
const GPTPrice = 0.0005 * 0.001

func (t Text) save(outPath string) error {
    log.Println("Saving file:", outPath)
    return os.WriteFile(outPath, []byte(t), 0644)
}

func (text Text) printPriceApprox() {
    log.Println("Tokenizing text...")
    tokens, err := text.tokenize()
    if err != nil {
        log.Fatal("Error:", err)
    }
    log.Println("Tokens:", tokens)
    log.Println("Expected price: ", float32(tokens) * GPTPrice, "USD")
}

func parseAndWrite(filePath string, parser func(string) (Text, error), outPath string) (Text, error) {
	log.Println("Parsing file:", filePath)
	text, err := parser(filePath)
	if err != nil {
		return "", err
	}
    return text, text.save(outPath)
}

func generateJsonAndSave(file Text, outPath string) (Text, error) {
    completion, err := gpt(string(file))
    if err != nil {
        return "", err
    }
    completionText := Text(completion)
    return completionText, completionText.save(outPath)
}

func generateExam(file Text, fileName string) {
    log.Println("Generating exam from file:", fileName)
    file.printPriceApprox()

    jsonPath := "outputs/" + fileName + ".json"
    completion, err := generateJsonAndSave(file, jsonPath)
    if err != nil {
        log.Fatal("Error:", err)
    }

    log.Println("Exam generated and saved to " + jsonPath)
    log.Println(completion)

    completion.printPriceApprox()
}

func main() {
	filePathDocx := "samples/tema1.1_introducción_computadoras.docx"
	filePathPdf := "samples/tema_3.1._gestión_de_la_memoria_paginación_y_segmentación.pdf"

	pdf, err := parseAndWrite(filePathPdf, parseDocument, "outputs/docpdf.txt")
	docx, err1 := parseAndWrite(filePathDocx, parseDocument, "outputs/docx.txt")

	if err != nil || err1 != nil {
		log.Fatal("Error:", err, err1)
	}

    generateExam(pdf, "pdf")
    generateExam(docx, "docx")
}
