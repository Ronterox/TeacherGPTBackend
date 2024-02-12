package main

import (
	"code.sajari.com/docconv/v2"
	// "github.com/lukasjarosch/go-docx"
)

func parseDocx2(docPath string) (string, error) {
    res, err := docconv.ConvertPath(docPath)
    if err != nil {
        return "", err
    }

    return res.Body, nil
}

// func parseDocx(docPath string) (string, error) {
//     doc, err := docx.Open(docPath)
//     if err != nil {
//         return "", err
//     }
//
//     return docx.NewRunParser
// }
