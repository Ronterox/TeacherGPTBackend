package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
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

// Handles the file sent in the request and returns the file data, the file header and an error if any
// If an error is returned, the file data will be the error code
func handleFile(r *http.Request) (Text, *multipart.FileHeader, error) {
	log.Println("Handling file...")
	file, handler, err := r.FormFile("file")
	if err != nil {
		return Text(strconv.Itoa(http.StatusBadRequest)), handler, err
	}
	defer file.Close()

	log.Println("Parsing " + handler.Filename + "...")
	fileData, err := parseFile(handler.Filename, file)
	if err != nil {
		return Text(strconv.Itoa(http.StatusInternalServerError)), handler, err
	}

	return Text(fileData), handler, nil
}

func generateMermaidImage(fileData Text) (string, error) {
	log.Println("Generating mindmap...")
	mindMap, err := gptMindMap(fileData)
	if err != nil {
		return "", err
	}

	bytesMap := bytes.Replace([]byte(mindMap), []byte("```mermaid"), []byte(""), 1)
	bytesMap = bytes.Replace(bytesMap, []byte("```"), []byte(""), 1)
	bytesMap = bytes.ReplaceAll(bytesMap, []byte("-"), []byte(" "))

	mindMap = strings.TrimSpace(string(bytesMap))
	log.Println("Mindmap:\n", mindMap)

	log.Println("Generating mermaid image...")
	return generateMermaidInkUrl(mindMap), nil
}

func imageToBase64(img []byte) string {
	return fmt.Sprintf("data:image/jpeg;base64,%v", base64.StdEncoding.EncodeToString(img))
}

func allowCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Headers, Origin,Accept, X-Requested-With, Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers")
}

func setJsonCORSHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	allowCORS(w)
}

func sendTemplate(w http.ResponseWriter, generateTemplate func() (string, error)) {
	setJsonCORSHeader(w)
	temp, err := generateTemplate()
	if err != nil {
		sendError(err, http.StatusInternalServerError, w)
		return
	}
	sendOk(w, []byte(temp))
}

func main() {
	http.HandleFunc("GET /api/template", func(w http.ResponseWriter, r *http.Request) { sendTemplate(w, getJsonTemplate) })

	http.HandleFunc("GET /api/template/open", func(w http.ResponseWriter, r *http.Request) { sendTemplate(w, getJsonTemplateOpen) })

	http.HandleFunc("POST /api/summary", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Generating summary image...")
		allowCORS(w)

		fileData, handler, err := handleFile(r)
		if code, _ := strconv.Atoi(string(fileData)); err != nil {
			sendError(err, code, w)
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

		img, err := generateMermaidImage(fileData)
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}

		sendOk(w, []byte(img))
	})

	http.HandleFunc("POST /api/generate/open", func(w http.ResponseWriter, r *http.Request) {
		setJsonCORSHeader(w)

		fileData, handler, err := handleFile(r)
		if code, _ := strconv.Atoi(string(fileData)); err != nil {
			sendError(err, code, w)
			return
		}

		exam, err := generateExam[QuestionOpen](fileData, handler.Filename)
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}

		sendOk(w, exam)
	})

	http.HandleFunc("POST /api/correct", func(w http.ResponseWriter, r *http.Request) {
		setJsonCORSHeader(w)
		var openQuestions []QuestionOpen
		if err := json.NewDecoder(r.Body).Decode(&openQuestions); err != nil {
			sendError(err, http.StatusBadRequest, w)
			return
		}

		stringJson, _ := json.MarshalIndent(openQuestions, "", "    ")
		res, err := gpt(string(stringJson), []gpt3.ChatCompletionRequestMessage{
			{Role: "system", Content: fmt.Sprintf("Return a valid json object with the correct answers for each question.\nThe scheme should follow the following:\n%v", stringJson)},
			{Role: "system", Content: "Make sure to write the comments in Spanish."}})

        if err != nil {
            sendError(err, http.StatusInternalServerError, w)
            return
        }

        sendOk(w, []byte(res))
	})

	http.HandleFunc("POST /api/generate", func(w http.ResponseWriter, r *http.Request) {
		setJsonCORSHeader(w)

		fileData, handler, err := handleFile(r)
		if code, _ := strconv.Atoi(string(fileData)); err != nil {
			sendError(err, code, w)
			return
		}

		exam, err := generateExam[QuestionSimple](Text(fileData), handler.Filename)
		if err != nil {
			sendError(err, http.StatusInternalServerError, w)
			return
		}

		sendOk(w, exam)
	})

	log.Printf("Server running on port %d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
