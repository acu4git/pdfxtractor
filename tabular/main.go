package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	csvPath := "../extractor/csv"
	dir, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("failed to open the directory: %v", err)
	}
	defer dir.Close()

	fileInfos, _ := dir.ReadDir(0)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// ファイルごとに処理する
	for _, fileInfo := range fileInfos {
		csvfile := filepath.Join(csvPath, fileInfo.Name())

		// ファイルを開く
		go func() {
			file, err := os.Open(csvfile)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", csvfile, err)
			}
			defer file.Close()
		}()

		// ここでファイルに対する操作を行う
		// 例: ファイルの内容を読み込む、書き込む、処理するなど
		// この例ではファイルの名前を出力する
		fmt.Println("Opened file:", fileInfo.Name())

	}
}

type target struct {
	course string
	detail string
}

type scholarship struct {
	updatedAt   string // time.Timeにする予定
	association string
	address     string
	target      target
	paymentInfo string
	capacity    string
	deadline    string
	pic         string
	remark      string
}
