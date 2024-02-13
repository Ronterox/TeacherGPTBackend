package main

import (
	"log"
	"testing"
)

func TestGeneration(t *testing.T) {
	filePathDocx := "samples/tema1.1_introducción_computadoras.docx"
	filePathPdf := "samples/tema_3.1._gestión_de_la_memoria_paginación_y_segmentación.pdf"

	pdf, err := parseAndWrite(filePathPdf, parseDocument, "outputs/docpdf.txt")
	docx, err1 := parseAndWrite(filePathDocx, parseDocument, "outputs/docx.txt")

	if err != nil || err1 != nil {
		log.Fatal("Error:", err, err1)
	}

    _, err = generateExam(pdf, "pdf")
    _, err1 = generateExam(docx, "docx")

	if err != nil || err1 != nil {
		log.Fatal("Error:", err, err1)
	}
}
