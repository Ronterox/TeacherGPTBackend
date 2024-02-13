package main

import (
	"fmt"

	"code.sajari.com/docconv/v2"
	"github.com/pkoukk/tiktoken-go"
	"github.com/pkoukk/tiktoken-go-loader"
)

func parseDocument(docPath string) (Text, error) {
	res, err := docconv.ConvertPath(docPath)
	if err != nil {
		return "", err
	}

	return Text(res.Body), nil
}

func parseAndWrite(filePath string, parser func(string) (Text, error), outPath string) (Text, error) {
	fmt.Println("Parsing file:", filePath)
	text, err := parser(filePath)
	if err != nil {
		return "", err
	}
	return text, text.save(outPath)
}

func (text Text) tokenize() (tokens int, err error) {
	tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())

	tke, err := tiktoken.EncodingForModel(GPTModel)
	if err != nil {
		return 0, fmt.Errorf("getEncoding: %v", err)
	}

	token := tke.Encode(string(text), nil, nil)
	return len(token), nil
}
