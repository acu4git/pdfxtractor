package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var wtr = bufio.NewWriter(os.Stdout)

func main() {
	// ディレクトリをオープンする
	csvPath := "../extractor/csv"
	dir, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("failed to open the directory: %v", err)
	}
	defer dir.Close()

	// ディレクトリ下の情報を取り出す
	fileInfos, _ := dir.ReadDir(0)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// ファイルごとに処理する
	scholarshipInfos := []scholarship{}
	for _, fileInfo := range fileInfos {
		csvFilename := filepath.Join(csvPath, fileInfo.Name())

		// ファイルを開く
		f, err := os.Open(csvFilename)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", csvFilename, err)
		}
		defer f.Close()

		// ここでファイルに対する操作を行う
		r := csv.NewReader(f)
		for {
			info := scholarship{}
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			for i, cell := range record {
				if shouldIgnore(cell) {
					// fmt.Println("---------------------Ignore!----------------------")
					continue
				}

				cell = strings.Replace(cell, "\n", "", -1)
				switch i {
				case 0:
					info.updatedAt = cell
				case 1:
					info.association = cell
				case 2:
					info.address = cell
				case 3:
					info.target.course = cell
				case 4:
					info.target.detail = cell
				case 5:
					info.paymentInfo.amountInfo = cell
				case 6:
					info.paymentInfo.scholarshipType = cell
				case 7:
					info.capacity = cell
				case 8:
					info.deadline = cell
				case 9:
					info.pic = cell
				case 10:
					info.remark = cell
				}
				// fmt.Printf("%v個目: %s\n", i, cell)
			}
			if info.updatedAt != "" {
				scholarshipInfos = append(scholarshipInfos, info)
			}
		}

	}
	for i, str := range scholarshipInfos {

		fmt.Fprintf(wtr, "---------------%d件目---------------\n", i)
		fmt.Fprintln(wtr, "掲示日", str.updatedAt)
		fmt.Fprintln(wtr, "奨学会名等", str.association)
		fmt.Fprintln(wtr, "住所", str.address)
		fmt.Fprintln(wtr, "対象(学部・院)", str.target.course)
		fmt.Fprintln(wtr, "対象(詳細)", str.target.detail)
		fmt.Fprintln(wtr, "年額・月額", str.paymentInfo.amountInfo)
		fmt.Fprintln(wtr, "貸与・給付", str.paymentInfo.scholarshipType)
		fmt.Fprintln(wtr, "募集人員", str.capacity)
		fmt.Fprintln(wtr, "申請期限等", str.deadline)
		fmt.Fprintln(wtr, "担当窓口", str.pic)
		fmt.Fprintln(wtr, "備考", str.remark)
		fmt.Fprintln(wtr)
		wtr.Flush()
	}

	os.RemoveAll(csvPath)
}

func shouldIgnore(cell string) bool {
	ignoreStrings := []string{
		"掲示日", "奨学会名等", "住所", "対象(学部・院)", "対象(詳細)", "年額・月額", "給与・貸与", "募集人員", "申請期限等", "担当窓口", "備考",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
	}
	for _, str := range ignoreStrings {
		if cell == str {
			return true
		}
	}
	return false
}

type target struct {
	course string
	detail string
}

type paymentInfo struct {
	amountInfo      string
	scholarshipType string
}

type scholarship struct {
	updatedAt   string // いずれtime.Timeにする予定
	association string
	address     string
	target      target
	paymentInfo paymentInfo
	capacity    string
	deadline    string
	pic         string
	remark      string
}
