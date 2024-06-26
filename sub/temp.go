package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

func extractFilename(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

func main() {
	// PDFファイルをオープン
	file, err := os.Open("../test.pdf")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Extract text from PDF
	fs, _ := file.Stat()
	fmt.Println("File size:", fs.Size())
	reader, err := pdf.NewReader(file, fs.Size())
	if err != nil {
		log.Fatalf("failed to create PDF reader: %v", err)
	}
	str := ""
	var extractedText strings.Builder
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		page := reader.Page(pageIndex)
		rows, err := page.GetTextByRow()
		text := "dummy"
		if err != nil {
			log.Printf("failed to extract text from page %d: %v", pageIndex, err)
			continue
		}
		for _, row := range rows {
			fmt.Println()
			for _, word := range row.Content {
				// fmt.Print(word.S)
				str = str + word.S + ", "
			}
			str = str + "\n"
		}
		extractedText.WriteString(text)
	}

	df, _ := os.Create("dump.csv")
	defer df.Close()
	df.WriteString(str)

	// Parse the extracted text and convert it to CSV format
	text := extractedText.String()
	lines := strings.Split(text, "\n")

	// Define the CSV headers
	headers := []string{
		"掲示日", "奨学金名等", "住所", "対象(学部・院)", "対象(詳細)",
		"年額・月額", "貸与・給付", "募集人員", "申請期限等", "担当窓口", "備考",
	}

	// Open a new CSV file
	csvFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("failed to create CSV file: %v", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write the CSV headers
	csvWriter.Write(headers)

	// Process the extracted text to fit into the CSV format
	var row []string
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		row = append(row, line)
		if len(row) == len(headers) {
			csvWriter.Write(row)
			row = nil
		}
	}

	fmt.Println("Text extracted and saved to output.csv")
}
