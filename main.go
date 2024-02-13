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

    _ = parseAndWrite(filePathPdf, parseDocument, "outputs/docpdf.txt")
    _ = parseAndWrite(filePathPdf, readPdfByRows, "outputs/pdfrows.txt")
    _ = parseAndWrite(filePathDocx, parseDocument, "outputs/docx.txt")

	// completion, err := gpt("What is the meaning of life?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(completion)
}
