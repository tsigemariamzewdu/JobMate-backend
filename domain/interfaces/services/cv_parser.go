package interfaces


import "mime/multipart"

type TextExtractor interface {
	Extract(fileHeader *multipart.FileHeader) (string, error)
}
