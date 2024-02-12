package main

import (
	"log"
)

func main() {
    filePath := "samples/tema1.1_introducción_computadoras.docx"
    text, err := parseDocx2(filePath)

    if err != nil {
        log.Fatal(err)
    }

    log.Println(text)

	// fileName := "tema_3.1._gestión_de_la_memoria_paginación_y_segmentación.pdf"
	// dirname := "samples/"
	//
	// text, err := readPdfByRows(dirname + fileName)
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// log.Println(text)

	// completion, err := gpt("What is the meaning of life?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(completion)
}
