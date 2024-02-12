package main

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

func readPdfPlain(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
    defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
    b, err := r.GetPlainText()
    if err != nil {
        return "", err
    }
    buf.ReadFrom(b)
	return buf.String(), nil
}

func readPdfStyle(path string) (string, error) {
	f, r, err := pdf.Open(path)
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		var lastTextStyle pdf.Text
		texts := p.Content().Text
		for _, text := range texts {
			if isSameSentence(text, lastTextStyle) {
				lastTextStyle.S = lastTextStyle.S + text.S
			} else {
				fmt.Printf("Font: %s, Font-size: %f, x: %f, y: %f, content: %s \n", lastTextStyle.Font, lastTextStyle.FontSize, lastTextStyle.X, lastTextStyle.Y, lastTextStyle.S)
				lastTextStyle = text
			}
		}
	}
	return "", nil
}

func isSameSentence(t1, t2 pdf.Text) bool {
    return t1.Y == t2.Y
}

func readPdfByRows(path string) (string, error) {
	f, r, err := pdf.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return "", err
	}
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
		    println(">>>> row: ", row.Position)
		    for _, word := range row.Content {
		        fmt.Println(word.S)
		    }
		}
	}
	return "", nil
}
