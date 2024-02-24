package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
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

func generateMermaidImage(fileData Text) ([]byte, error) {
	log.Println("Generating mindmap...")
	mindMap, err := gptMindMap(fileData)
	if err != nil {
		return nil, err
	}

	bytesMap := bytes.Replace([]byte(mindMap), []byte("```mermaid"), []byte(""), 1)
	bytesMap = bytes.Replace(bytesMap, []byte("```"), []byte(""), 1)
	bytesMap = bytes.ReplaceAll(bytesMap, []byte("-"), []byte(" "))

	mindMap = strings.TrimSpace(string(bytesMap))
	log.Println("Mindmap:\n", mindMap)

	log.Println("Generating mermaid image...")
	mermaidUrl := generateMermaidInkUrl(mindMap)
	return getImageFromUrl(mermaidUrl)
}

func imageToBase64(img []byte) string {
	return fmt.Sprintf("data:image/jpeg;base64,%v", base64.StdEncoding.EncodeToString(img))
}

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func main() {
	var currentlyGeneratingFiles = make(map[string]int)

	http.HandleFunc("GET /api/template", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		allowCORS(w)
		temp, err := getJsonTemplate()
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}
		sendOk(w, []byte(temp))
	})

	http.HandleFunc("POST /api/summary", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Generating summary image...")
		allowCORS(w)

		fileData, handler, err := handleFile(w, r)
		if err != nil {
			sendError(err, http.StatusBadRequest, w)
			return
		}

		filePath := fmt.Sprintf("outputs/%v.jpg", handler.Filename)

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			log.Println("Caching image file...")
			img, err := os.ReadFile(filePath)
			if err != nil {
				sendError(err, http.StatusInternalServerError, w)
				return
			}
			sendOk(w, []byte(imageToBase64(img)))
			return
		}

        img, err := generateMermaidImage(fileData); 
        for ;err != nil; img, err = generateMermaidImage(fileData) {
            currentlyGeneratingFiles[handler.Filename]++
            if currentlyGeneratingFiles[handler.Filename] > 3 {
                sendError(err, http.StatusInternalServerError, w)
                return
            }
        }

		if err := generateDir("outputs"); err == nil {
			Text(img).save(filePath)
		}

		sendOk(w, []byte(imageToBase64(img)))
	})

	http.HandleFunc("POST /api/generate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		allowCORS(w)

		fileData, handler, err := handleFile(w, r)
		if err != nil {
			sendError(err, http.StatusBadRequest, w)
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
