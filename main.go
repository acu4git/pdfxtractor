package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/unidoc/unipdf/v3/common/license"
)

func init() {
	apiKey := os.Getenv("UNIPDF_API_KEY")
	err := license.SetMeteredKey(apiKey)
	if err != nil {
		panic(err)
	}
}

func main() {
	pdfURL := "https://www.kit.ac.jp/wp/wp-content/uploads/2024/04/R6hpsyougakukinitiran20240408.pdf"

	filename := getFilename(pdfURL)

	err := downloadFile(pdfURL, filename)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	fmt.Println("File downloaded successfully")

	time.Sleep(time.Second * 5)
	if err = deleteFile(filename); err != nil {
		fmt.Println("Error deleting file:", err)
	}
}

func getFilename(path string) string {
	pathElems := strings.Split(path, "/")
	return pathElems[len(pathElems)-1]
}

func downloadFile(url string, filename string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func deleteFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}

	return nil
}
