package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	service "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	pdfparser "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/pdf_parser"

	"github.com/unidoc/unioffice/document"
)

type FileTextExtractor struct{}

func NewFileTextExtractor() service.TextExtractor {
	return &FileTextExtractor{}
}

func (e *FileTextExtractor) Extract(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".pdf":
		return extractPDFText(file)
	case ".docx":
		return extractDocxText(file)
	case ".txt":
		return extractTxtText(file)
	default:
		return "", fmt.Errorf("unsupported file type: %s", ext)
	}
}

func extractPDFText(file multipart.File) (string, error) {
	tmpPath, err := saveTempFile(file, "upload-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	text, err := pdfparser.ParsePDF(tmpPath)
	if err != nil {
		return "", err
	}
	return text, nil
}

func extractDocxText(file multipart.File) (string, error) {
	tmp, err := saveTempFile(file, "upload-*.docx")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp)

	doc, err := document.Open(tmp)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	for _, para := range doc.Paragraphs() {
		for _, run := range para.Runs() {
			buf.WriteString(run.Text())
			buf.WriteString(" ")
		}
	}
	return buf.String(), nil
}

func extractTxtText(file multipart.File) (string, error) {
	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
