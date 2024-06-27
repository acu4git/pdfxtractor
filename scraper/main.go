package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ledongthuc/pdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// 指定したURL(url)のレスポンスから，検索パターン(pattern)に合致するURLを1つ返す
func scrape(url, pattern string) (string, error) {
	// スクレイピング対象のページにリクエストを送る
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch the URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch the URL: status code %d", res.StatusCode)
	}

	// goqueryドキュメントを作成
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse the HTML: %v", err)
	}

	// 検索パターン(pattern)をGoの正規表現としてコンパイル
	re := regexp.MustCompile(pattern)

	// リンクを探す
	found := false
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, exists := item.Attr("href")
		if exists && re.MatchString(href) {
			fmt.Printf("found link: %s\n", href)
			url = href
			found = true
		}
	})

	if !found {
		return "", fmt.Errorf("not found link")
	}

	return url, nil
}

func extractFilename(url string) string {
	tokens := strings.Split(url, "/")
	return tokens[len(tokens)-1]
}

func main() {
	// スクレイピング
	url := "https://www.kit.ac.jp/campus_index/life_fee/scholarship/minkanscholarship/" // 対象PDFファイルが掲載されているHP
	linkPattern := `https://www\.kit\.ac\.jp/wp/wp-content/uploads/\d{4}/\d{2}/.*hpsyougakukinitiran.*\.pdf`
	url, err := scrape(url, linkPattern)
	if err != nil {
		log.Fatal(err)
	}

	// PDFレスポンスを受け取る
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("failed to fetch the URL: status code %d", res.StatusCode)
	}

	// 一旦PDFファイルをローカルに保存
	inFile := extractFilename(url)
	f, err := os.Create(inFile)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(f, res.Body)

	os.Mkdir("dump", 0777)
	// defer os.RemoveAll("dump")
	path := "./dump"

	// PDFを分割（いらないかも）
	conf := model.NewDefaultConfiguration()
	selectedPages := []string{"1-"} // Extract text from all pages
	err = api.ExtractPagesFile(inFile, path, selectedPages, conf)
	if err != nil {
		log.Fatal(err)
	}

	// PDFファイルをオープン
	file, err := os.Open("test.pdf")
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

	var extractedText strings.Builder
	for pageIndex := 1; pageIndex <= reader.NumPage(); pageIndex++ {
		page := reader.Page(pageIndex)
		rows, err := page.GetTextByRow()
		if err != nil {
			log.Printf("failed to extract text from page %d: %v", pageIndex, err)
			continue
		}
		str := ""
		for _, row := range rows {
			for _, word := range row.Content {
				if word.S == " " {
					continue
				}
				str += word.S
			}
			str += "\n"
		}
		extractedText.WriteString(str)
	}

	// Parse the extracted text and convert it to CSV format
	text := extractedText.String()
	fmt.Println(text)
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
