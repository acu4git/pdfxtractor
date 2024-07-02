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

type userTarget int

const (
	bachelor userTarget = iota + 1
	master
	other
	all
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
	postDate    string
	association string
	address     string
	target      target
	paymentInfo paymentInfo
	capacity    string
	deadline    string
	pic         string
	remark      string
}

func getUserInput() userTarget {
	fmt.Println()
	fmt.Println("表示させたい奨学金一覧の対象について数字(1~4)で入力してください．")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("1. 学部生")
	fmt.Println("2. 大学院生")
	fmt.Println("3. その他")
	fmt.Println("4. 全部")
	fmt.Println("-------------------------------------------------------------------")

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

	return userTarget(num)
}

// 無視する項目の処理
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

// 元号表記(令和yy年mm月dd日)の日付をyyyy-mm-dd形式に変換
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
	// ユーザー入力部
	tagetInfo := getUserInput()

	// 現在時刻の取得
	curr := time.Now()

	// それぞれのCSVファイルについて処理を行う
	scholarshipInfos := []scholarship{}
	csvPath := "../extractor/csv/"
	err := filepath.WalkDir(csvPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合はスキップ
		if d.IsDir() {
			return nil
		}

		// // デバッグ出力
		// fmt.Printf("file 1: %s\n", d.Name())

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// レコードの処理
		r := csv.NewReader(f)
		for {
			info := scholarship{}
			// 1つのレコードを取得
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

				// 下処理
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
					info.postDate = field
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

			// 掲示日が古すぎないかどうか
			if info.postDate == "" {
				continue
			}
			date, err := parseDate(info.postDate)
			if err != nil {
				fmt.Println(err)
				continue
			}
			old, _ := time.Parse("2006-01-02", date)
			diff := curr.Sub(old)
			// 掲示日から1年以上経っている情報は捨てる
			if int(diff.Hours()/24) > 365 {
				continue
			}

			// 指定した対象を含んだデータのみ抽出
			switch tagetInfo {
			case bachelor:
				if !strings.Contains(info.target.course, "学部") {
					continue
				}
			case master:
				if !strings.Contains(info.target.course, "大学院") {
					continue
				}
			case other:
				if !strings.Contains(info.target.course, "その他") {
					continue
				}
			case all: //何もしない
			default:
				fmt.Println("このメッセージは表示されないはずだぴょん:", "targetInfo")
			}

			// 申請期限が切れているかどうか
			dateInfo, _, found := strings.Cut(info.deadline, "(")
			if found {
				date, err := parseDate(dateInfo)
				if err != nil {
					fmt.Println(err)
					continue
				}
				if date < curr.Format("2006-01-02") {
					continue
				}
			}

			// 条件が合致していれば追加
			scholarshipInfos = append(scholarshipInfos, info)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("../result.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(outFile, "更新日: %d年%d月%d日\n\n", curr.Year(), curr.Month(), curr.Day())
	for i, str := range scholarshipInfos {
		fmt.Fprintf(outFile, "---------------%d件目---------------\n", i+1)
		fmt.Fprintln(outFile, "***掲示日***\n", str.postDate)
		fmt.Fprintln(outFile, "\n***奨学会名等***\n", str.association)
		fmt.Fprintln(outFile, "\n***住所***\n", str.address)
		fmt.Fprintln(outFile, "\n***対象(学部・院)***\n", str.target.course)
		fmt.Fprintln(outFile, "\n***対象(詳細)***\n", str.target.detail)
		fmt.Fprintln(outFile, "\n***年額・月額***\n", str.paymentInfo.amountInfo)
		fmt.Fprintln(outFile, "\n***貸与・給付***\n", str.paymentInfo.scholarshipType)
		fmt.Fprintln(outFile, "\n***募集人員***\n", str.capacity)
		fmt.Fprintln(outFile, "\n***申請期限等***\n", str.deadline)
		fmt.Fprintln(outFile, "\n***担当窓口***\n", str.pic)
		fmt.Fprintln(outFile, "\n***備考***\n", str.remark)
		fmt.Fprintln(outFile)
	}

	fmt.Println("Saved to result.txt!")
	os.RemoveAll(csvPath)
}
