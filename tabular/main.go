package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	// ディレクトリをオープンする
	csvPath := "../extractor/csv"
	dir, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("failed to open the directory: %v", err)
	}
	defer dir.Close()

	// ディレクトリ下の情報を取り出す(順番は保証されていない)
	fileInfos, _ := dir.ReadDir(0)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// ソートしてファイル番号を昇順にする
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].Name() < fileInfos[j].Name()
	})

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
			for i, field := range record {
				if shouldIgnore(field) {
					continue
				}

				field = strings.Replace(field, " ", "", -1)
				field = strings.Replace(field, "①", "(1)", -1)
				field = strings.Replace(field, "②", "(2)", -1)
				field = strings.Replace(field, "③", "(3)", -1)
				field = strings.Replace(field, "④", "(4)", -1)
				field = strings.Replace(field, "⑤", "(5)", -1)
				field = strings.Replace(field, "⑥", "(6)", -1)
				field = strings.Replace(field, "⑦", "(7)", -1)
				field = strings.Replace(field, "⑧", "(8)", -1)
				field = strings.Replace(field, "⑨", "(9)", -1)

				switch i {
				case 0:
					info.updatedAt = field
				case 1:
					info.association = field
				case 2:
					info.address = field
				case 3:
					info.target.course = field
				case 4:
					info.target.detail = field
				case 5:
					info.paymentInfo.amountInfo = field
				case 6:
					info.paymentInfo.scholarshipType = field
				case 7:
					info.capacity = field
				case 8:
					info.deadline = field
				case 9:
					info.pic = field
				case 10:
					info.remark = field
				}
				// fmt.Printf("%v個目: %s\n", i, field)
			}
			if info.updatedAt != "" {
				scholarshipInfos = append(scholarshipInfos, info)
			}
		}

	}
	for i, str := range scholarshipInfos {
		fmt.Printf("---------------%d件目---------------\n", i+1)
		fmt.Println("***掲示日***\n", str.updatedAt)
		fmt.Println("\n***奨学会名等***\n", str.association)
		fmt.Println("\n***住所***\n", str.address)
		fmt.Println("\n***対象(学部・院)***\n", str.target.course)
		fmt.Println("\n***対象(詳細)***\n", str.target.detail)
		fmt.Println("\n***年額・月額***\n", str.paymentInfo.amountInfo)
		fmt.Println("\n***貸与・給付***\n", str.paymentInfo.scholarshipType)
		fmt.Println("\n***募集人員***\n", str.capacity)
		fmt.Println("\n***申請期限等***\n", str.deadline)
		fmt.Println("\n***担当窓口***\n", str.pic)
		fmt.Println("\n***備考***\n", str.remark)
		fmt.Println()
	}

	os.RemoveAll(csvPath)
}

func shouldIgnore(field string) bool {
	ignoreStrings := []string{
		"掲示日", "奨学会名等", "住所", "対象(学部・院)", "対象(詳細)", "年額・月額", "給与・貸与", "募集人員", "申請期限等", "担当窓口", "備考",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
	}
	for _, str := range ignoreStrings {
		if field == str {
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
	updatedAt   string
	association string
	address     string
	target      target
	paymentInfo paymentInfo
	capacity    string
	deadline    string
	pic         string
	remark      string
}
