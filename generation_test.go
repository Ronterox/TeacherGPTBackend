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

    _, err = generateExam[Question](pdf, "pdf", false)
    _, err1 = generateExam[Question](docx, "docx", false)

	if err != nil || err1 != nil {
		log.Fatal("Error:", err, err1)
	}
}

func TestMermaid(t *testing.T) {
	mermaid := `flowchart TD
   1(All started must be added to the list.\nIf started it will be finished no matter what\nYES or YES)
   2(Creativity is the most powerful and most important\nEs facil salirse de la norma, solo ve lo que hace la gente)
   3(Todo no es mas que un draft para lo que se viene despues\nHay que amar las cosas, pero apegarse no es una ventaja)
   4(Is always about the next thing, what comes after being done\nyou don't think about how will you feel doing a thing\nyou think about how it feels after it, and what comes)
   5(Luck is my ultimate ability, also my favorite skill of mine\nTogether with many others, I can take all risks)

   Rules -- "1 First 1" --- 
   1 -- "2 Second 2" --> 
   2 -- "3 Third 3" --> 
   3 -- "4 Fourth 4" --> 
   4 -- "5 Fifth 5" --> 5`

    url := generateMermaidInkUrl(mermaid)
    log.Println(url)

    img, err := getImageFromUrl(url)
    if err != nil {
        log.Fatal("Error:", err)
    }

    log.Println("Body:", string(img))
}
