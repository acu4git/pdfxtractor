package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type userTarget uint16

const (
	bachelor = iota
	master
)

type target struct {
	course string
	detail string
}

type paymentInfo struct {
	amountInfo      string
	scholarshipType string
}

type scholarship struct {
	postDay     string
	association string
	address     string
	target      target
	paymentInfo paymentInfo
	capacity    string
	deadline    string
	pic         string
	remark      string
}

func getUserInput() int {
	fmt.Println("表示させたい奨学金一覧の対象について数字(1~4)で入力してください．")
	fmt.Println("1. 学部生")
	fmt.Println("2. 大学院生")
	fmt.Println("3. その他")
	fmt.Println("4. 全部")
	fmt.Println()

	var input string
	var num int
	var err error

	for {
		fmt.Print(" > ")
		reader := bufio.NewReader(os.Stdin)
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			continue
		}
		num, err = strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			fmt.Println("入力に不適切な文字が含まれています．")
			continue
		}
		if num < 1 || 4 < num {
			fmt.Println("1, 2, 3, 4のいずれかを入力してください．")
			continue
		}

		break
	}

	return num
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

// 元号表記の日付をyyyy-mm-dd形式に変換
func parseDate(dateStr string) (string, error) {
	// 元年からのオフセット(適時追加)
	const reiwaStartYear = 2019

	// 令和の日付を分割して解析
	parts := strings.Split(dateStr, "年")
	if len(parts) != 2 {
		return "", fmt.Errorf("無効なフォーマットです: %s", dateStr)
	}

	// 年の部分を取得し、整数に変換
	reiwaYear, err := strconv.Atoi(strings.TrimPrefix(parts[0], "令和"))
	if err != nil {
		return "", fmt.Errorf("無効な年: %s", parts[0])
	}

	// 日付の部分を分割して月日を取得
	dateParts := strings.Split(parts[1], "月")
	if len(dateParts) != 2 {
		return "", fmt.Errorf("無効なフォーマットです: %s", parts[1])
	}
	month, err := strconv.Atoi(dateParts[0])
	if err != nil {
		return "", fmt.Errorf("無効な月: %s", dateParts[0])
	}
	day, err := strconv.Atoi(strings.TrimSuffix(dateParts[1], "日"))
	if err != nil {
		return "", fmt.Errorf("無効な日: %s", dateParts[1])
	}

	// 西暦年を計算
	gregorianYear := reiwaStartYear + reiwaYear - 1

	// timeパッケージを使用してフォーマット
	date := time.Date(gregorianYear, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return date.Format("2006-01-02"), nil
}

func main() {
	fmt.Println("")

	scholarshipInfos := []scholarship{}

	// それぞれのCSVファイルについて処理を行う
	csvPath := "../extractor/csv/"
	err := filepath.WalkDir(csvPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// ディレクトリの場合はスキップ
			return nil
		}

		// // デバッグ出力
		// fmt.Printf("file 1: %s\n", d.Name())

		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}

		r := csv.NewReader(f)
		for {
			info := scholarship{}
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			for i, field := range record {
				if shouldIgnore(field) {
					continue
				}

				field = strings.TrimSpace(field)
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
					info.postDay = field
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

			dateInfo, _, found := strings.Cut(info.deadline, "(")
			// 申請期限がある場合について
			if found {
				date, err := parseDate(dateInfo)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if date < time.Now().Format("2006-01-02") {
					continue
				}
			}
			if info.postDay != "" {
				scholarshipInfos = append(scholarshipInfos, info)
			}

		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	for i, str := range scholarshipInfos {
		fmt.Printf("---------------%d件目---------------\n", i+1)
		fmt.Println("***掲示日***\n", str.postDay)
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

	// os.RemoveAll(csvPath)
}
