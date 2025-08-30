package fileparser

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/fumiama/go-docx"
	"github.com/ledongthuc/pdf"
	service "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
)

type FileTextExtractor struct{}

func NewFileTextExtractor() service.TextExtractor {
	return &FileTextExtractor{}
}

func (e *FileTextExtractor) Extract(fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".pdf":
		return extractPDFText(fileHeader)
	case ".docx":
		return extractDocxText(fileHeader)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func extractPDFText(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	tmpPath, err := saveTempFile(file, "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	f, r, err := pdf.Open(tmpPath)
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

func extractDocxText(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	tmpPath, err := saveTempFile(file, "upload-*.docx")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	f, err := os.Open(tmpPath)
	if err != nil {
		return "", fmt.Errorf("failed to open temp docx: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat docx: %w", err)
	}

	doc, err := docx.Parse(f, info.Size())
	if err != nil {
		return "", fmt.Errorf("failed to parse docx: %w", err)
	}

	var text string
	for _, it := range doc.Document.Body.Items {
		switch v := it.(type) {
		case *docx.Paragraph:
			text += v.String() + "\n"
		case *docx.Table:

			text += v.String() + "\n"
		}
	}

	log.Print(text)
	return text, nil
}

func saveTempFile(file multipart.File, pattern string) (string, error) {
	tmp, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	_, err = io.Copy(tmp, file)
	if err != nil {
		return "", err
	}
	return tmp.Name(), nil
}
