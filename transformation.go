package main

import (
	"io"

	"code.sajari.com/docconv/v2"
	"github.com/pkoukk/tiktoken-go"
	tiktoken_loader "github.com/pkoukk/tiktoken-go-loader"
)

func parseDocument(docPath string) (Text, error) {
	res, err := docconv.ConvertPath(docPath)
	if err != nil {
		return "", err
	}

	return Text(res.Body), nil
}

func parseFile(fileName string, file io.Reader) (Text, error) {
	res, err := docconv.Convert(file, docconv.MimeTypeByExtension(fileName), true)
	if err != nil {
		return "", err
	}

	return Text(res.Body), nil
}

func (text Text) tokenize() (tokens int) {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())

	tke, _ := tiktoken.EncodingForModel(GPTModel)

	token := tke.Encode(string(text), nil, nil)
	return len(token)
}
