package main

import (
	"fmt"

	"code.sajari.com/docconv/v2"
	"github.com/ledongthuc/pdf"
)

func parseDocument(docPath string) (string, error) {
    res, err := docconv.ConvertPath(docPath)
    if err != nil {
        return "", err
    }

    return res.Body, nil
}

func readPdfByRows(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer f.Close()

	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	result := ""
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			row_text := ""
			for _, word := range row.Content {
				row_text += word.S
			}
			result += fmt.Sprintf("%s\n", row_text)
		}
	}
	return result, nil
}
