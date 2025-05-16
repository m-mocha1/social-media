package main

import (
	"encoding/base64"
	"io"
	"mime/multipart"
	"os"
)

func readImg(path string, file multipart.File) ([]byte, error) {
	if path != "" {
		img, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer img.Close()
		imgbyte, err := io.ReadAll(img)
		if err != nil {
			return nil, err
		}
		return imgbyte, nil
	} else {
		imgbyte, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return imgbyte, nil
	}
}
func encodeToBase64(imgData []byte) string {
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(imgData)
}
