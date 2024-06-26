package main

import (
	"fmt"
	"log"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetLanguage("eng", "jpn")

	var fullText string
	imgPath := "../ocr/test_00.jpg"
	client.SetImage(imgPath)
	text, err := client.Text()
	if err != nil {
		log.Fatalf("failed to perform OCR on image %s: %v", imgPath, err)
	}
	fullText += text
	fmt.Println(fullText)
}
