package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = 8080

func sendError(err error, statusCode int, w http.ResponseWriter) {
    w.WriteHeader(statusCode)
    w.Write([]byte(err.Error()))
    log.Println(err)
}

func sendOk(w http.ResponseWriter, data []byte) {
    w.WriteHeader(http.StatusOK)
    w.Write(data)
}

func main() {
	http.HandleFunc("GET /template", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		temp, err := getJsonTemplate()
		if err != nil {
            sendError(err, http.StatusInternalServerError, w)
			return
		}
        sendOk(w, []byte(temp))
	})

	// Receive file and generate exam
	http.HandleFunc("POST /generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		file, handler, err := r.FormFile("file")
		if err != nil {
            sendError(err, http.StatusBadRequest, w)
			return
		}
		defer file.Close()

		fileData, err := parseFile(handler.Filename, file)
		if err != nil {
            sendError(err, http.StatusInternalServerError, w)
			return
		}

		exam, err := generateExam(Text(fileData), handler.Filename)
		if err != nil {
            sendError(err, http.StatusInternalServerError, w)
			return
		}

        sendOk(w, exam)
	})

	log.Printf("Server running on port %d", PORT)
    http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
