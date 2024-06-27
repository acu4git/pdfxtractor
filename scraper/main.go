package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
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
	inFile := "../extractor/target.pdf"
	f, err := os.Create(inFile)
	if err != nil {
		log.Fatal(err)
	}
	io.Copy(f, res.Body)
}
