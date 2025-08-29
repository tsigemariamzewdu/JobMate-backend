package pdfparser

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ParsePDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %v", err)
	}
	defer f.Close()

	var textBuilder strings.Builder
	totalPage := r.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, err := p.GetTextByRow()
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %v", pageIndex, err)
		}
		if len(rows) == 0 {
			continue
		}

		for _, row := range rows {
			var rowText strings.Builder
			prevX := -1.0
			for _, word := range row.Content {
				if prevX >= 0 && word.X-prevX > 1.5 {
					rowText.WriteString(" ")
				}
				rowText.WriteString(word.S)
				prevX = word.X + word.W
			}
			textBuilder.WriteString(strings.TrimSpace(rowText.String()) + "\n")
		}
		textBuilder.WriteString("\n")
	}

	return textBuilder.String(), nil
}
