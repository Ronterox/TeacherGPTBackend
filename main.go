package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
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

func handleFile(w http.ResponseWriter, r *http.Request) (Text, *multipart.FileHeader, error) {
    log.Println("Handling file...")
	file, handler, err := r.FormFile("file")
	if err != nil {
		sendError(err, http.StatusBadRequest, w)
		return "", handler, err
	}
	defer file.Close()

    log.Println("Parsing " + handler.Filename + "...")
	fileData, err := parseFile(handler.Filename, file)
	if err != nil {
		sendError(err, http.StatusInternalServerError, w)
		return "", handler, err
	}

	return Text(fileData), handler, nil
}

func handleMermaid(mermaid string) ([]byte, error) {
    mermaidUrl := generateMermaidInkUrl(mermaid)
    return getImageFromUrl(mermaidUrl)
}

func main() {
	http.HandleFunc("GET /api/template", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		temp, err := getJsonTemplate()
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}
		sendOk(w, []byte(temp))
	})

	http.HandleFunc("POST /api/parse", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Parsing file...")
		if fileData, _, err := handleFile(w, r); err == nil {
			sendOk(w, []byte(fileData))
		}
	})

	http.HandleFunc("POST /api/mermaid", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
        log.Println("Generating mermaid image...")

		body, err := io.ReadAll(r.Body)
        if err != nil {
            sendError(err, http.StatusBadRequest, w)
            return
        }

		img, err := handleMermaid(string(body))
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}

		sendOk(w, img)
	})

	http.HandleFunc("POST /api/summary", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
        log.Println("Generating summary image...")

        fileData, _, err := handleFile(w, r)
        if err != nil {
            return
        }

        log.Println("Generating mindmap...")
        mindMap, err := gptMindMap(fileData)
        if err != nil {
            sendError(err, http.StatusInternalServerError, w)
            return
        }

        bytesMap := bytes.Replace([]byte(mindMap), []byte("```mermaid"), []byte(""), 1)
        bytesMap = bytes.Replace(bytesMap, []byte("```"), []byte(""), 1)
        bytesMap = bytes.ReplaceAll(bytesMap, []byte("-"), []byte(" "))

        mindMap = strings.TrimSpace(string(bytesMap))
        log.Println("Mindmap:", mindMap)

        log.Println("Generating mermaid image...")
        img, err := handleMermaid(mindMap)
        if err != nil {
            sendError(err, http.StatusInternalServerError, w)
            return
        }

        sendOk(w, []byte(base64.StdEncoding.EncodeToString(img)))
	})

	http.HandleFunc("POST /api/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fileData, handler, err := handleFile(w, r)
		if err != nil {
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
