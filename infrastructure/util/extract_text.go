package utils

import (
	// "bytes"
	"fmt"
	"io"
	"log"

	// "log"
	"mime/multipart"
	"os"
	"path/filepath"

	// "github.com/lu4p/cat"
	"github.com/fumiama/go-docx"
	service "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/services"
	pdfparser "github.com/tsigemariamzewdu/JobMate-backend/infrastructure/pdf_parser"
	// "github.com/unidoc/unioffice/document"
)

type FileTextExtractor struct{}

func NewFileTextExtractor() service.TextExtractor {
	return &FileTextExtractor{}
}

func (e *FileTextExtractor) Extract(fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	switch ext {
	case ".pdf":
		file, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
		return extractPDFText(file)
	case ".docx":
		return extractDocxText(fileHeader)
	case ".txt":
		file, err := fileHeader.Open()
		if err != nil {
			return "", err
		}
		defer file.Close()
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

// func extractDocxText(file multipart.File) (string, error) {
// 	tmp, err := saveTempFile(file, "upload-*.docx")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer os.Remove(tmp)

// 	doc, err := document.Open(tmp)
// 	if err != nil {
// 		return "", err
// 	}

// 	var buf bytes.Buffer
// 	for _, para := range doc.Paragraphs() {
// 		for _, run := range para.Runs() {
// 			buf.WriteString(run.Text())
// 			buf.WriteString(" ")
// 		}
// 	}
// 	return buf.String(), nil
// }

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
