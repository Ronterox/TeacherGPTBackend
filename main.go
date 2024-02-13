package main

import (
	"log"
	"os"
)

func parseAndWrite(filePath string, parser func(string) (string, error), outPath string) error {
	log.Println("Parsing file:", filePath)
	text, err := parser(filePath)
	if err != nil {
		return err
	}

	log.Println("Writing file:", outPath)
	return os.WriteFile(outPath, []byte(text), 0644)
}

func main() {
	filePathDocx := "samples/tema1.1_introducción_computadoras.docx"
	filePathPdf := "samples/tema_3.1._gestión_de_la_memoria_paginación_y_segmentación.pdf"

	err := parseAndWrite(filePathPdf, parseDocument, "outputs/docpdf.txt")
	err1 := parseAndWrite(filePathPdf, readPdfByRows, "outputs/pdfrows.txt")
	err2 := parseAndWrite(filePathDocx, parseDocument, "outputs/docx.txt")

	if err != nil || err1 != nil || err2 != nil {
		log.Fatal("Error:", err, err1, err2)
	}

    text, err := os.ReadFile("outputs/docx.txt")
	completion, err := gpt(string(text))
	if err != nil {
		log.Fatal(err)
	}
	log.Println(completion)
    os.WriteFile("outputs/gpt.json", []byte(completion), 0644)
}
