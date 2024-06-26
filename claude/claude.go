package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Scholarship struct {
	Name             string
	Organization     string
	Deadline         string
	Amount           string
	EligibleStudents string
}

func main() {
	// 1. PDFからテキスト抽出
	content, err := readPdf("R6hpsyougakukinitiran20240624_page_1.pdf")
	if err != nil {
		fmt.Println("Error reading PDF:", err)
		return
	}
	fmt.Println(content)

	// 2. テキストの解析と構造化
	scholarships := parseScholarships(content)

	// 3. データの検索・フィルタリング
	// 例: 締め切りが "令和6年7月31日" の奨学金を検索
	filtered := filterScholarships(scholarships, func(s Scholarship) bool {
		return s.Deadline == "令和6年7月31日(水)"
	})

	// 結果の表示
	for _, s := range filtered {
		fmt.Printf("Name: %s\nOrganization: %s\nDeadline: %s\nAmount: %s\nEligible Students: %s\n\n",
			s.Name, s.Organization, s.Deadline, s.Amount, s.EligibleStudents)
	}
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

func parseScholarships(content string) []Scholarship {
	lines := strings.Split(content, "\n")
	var scholarships []Scholarship
	var current Scholarship

	for _, line := range lines {
		if strings.Contains(line, "奨学会名等") {
			if current.Name != "" {
				scholarships = append(scholarships, current)
			}
			current = Scholarship{}
			current.Name = strings.TrimSpace(line)
		} else if strings.Contains(line, "申請期限等") {
			current.Deadline = strings.TrimSpace(line)
		} else if strings.Contains(line, "年額・月額") {
			current.Amount = strings.TrimSpace(line)
		} else if strings.Contains(line, "対象(詳細)") {
			current.EligibleStudents = strings.TrimSpace(line)
		} else if strings.Contains(line, "住所") {
			current.Organization = strings.TrimSpace(line)
		}
	}

	if current.Name != "" {
		scholarships = append(scholarships, current)
	}

	return scholarships
}

func filterScholarships(scholarships []Scholarship, predicate func(Scholarship) bool) []Scholarship {
	var filtered []Scholarship
	for _, s := range scholarships {
		if predicate(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
