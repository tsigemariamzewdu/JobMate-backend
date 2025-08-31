package controllers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/utils"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	usecase "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
)

const (
	maxUploadBytes = 10 << 20
)

type CVController struct {
	cvUsecase usecase.ICVUsecase
}

func NewCVController(u usecase.ICVUsecase) *CVController {
	return &CVController{cvUsecase: u}
}

type CVUploadRequest struct {
	
	RawText string                `json:"rawText" form:"rawText"`
	File    *multipart.FileHeader `form:"file"`
}

// POST /cv
func (c *CVController) UploadCV(ctx *gin.Context) {

	userID:=ctx.GetString("user_id")
	// Limit request size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxUploadBytes)

	var req CVUploadRequest
	if err := ctx.ShouldBind(&req); err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			ctx.JSON(http.StatusRequestEntityTooLarge, utils.ErrorPayload("File exceeds max size of 10MB", nil))
			return
		}
		ctx.JSON(http.StatusBadRequest, utils.ErrorPayload("Invalid input", err.Error()))
		return

	}

	if req.RawText == "" && req.File == nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorPayload("Either rawText or file must be provided", nil))
		return
	}

	if req.RawText != "" && req.File != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorPayload("Cannot provide both rawText and file at the same time", nil))
		return
	}

	// validate file
	if req.File != nil {
		if req.File.Size <= 0 || req.File.Size > maxUploadBytes {
			ctx.JSON(http.StatusRequestEntityTooLarge, utils.ErrorPayload("File exceeds max size of 10MB", nil))
			return
		}

		// Sniff MIME type
		file, err := req.File.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorPayload("Could not open file", err.Error()))
			return
		}
		defer file.Close()

		head := make([]byte, 512)
		n, _ := io.ReadFull(file, head)
		mime := http.DetectContentType(head[:n])
		ext := strings.ToLower(path.Ext(req.File.Filename))
		allowed := map[string]bool{
			"application/pdf": true,
		}
		if !allowed[mime] && ext != ".docx" {
			ctx.JSON(http.StatusUnsupportedMediaType, utils.ErrorPayload(
				"Only PDF or DOCX files are allowed",
				map[string]any{"detected": mime},
			))
			return
		}

		// Sanitize filename
		req.File.Filename = path.Base(req.File.Filename)
	}

	createdCV, err := c.cvUsecase.Upload(ctx,userID, req.RawText, req.File)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCVNotFound):
			ctx.JSON(http.StatusNotFound, utils.ErrorPayload("CV not found", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorPayload("Failed to upload CV", err.Error()))
		}
		return
	}

	ctx.Header("Location", "/cv/"+createdCV.ID)

	ctx.JSON(http.StatusCreated, utils.SuccessPayload("CV uploaded successfully", gin.H{
		"cvId":      createdCV.ID,
		"userId":    createdCV.UserID,
		"fileName":  createdCV.FileName,
		"createdAt": createdCV.CreatedAt,
	}))
}

// POST /cv/:id/analyze
func (c *CVController) AnalyzeCV(ctx *gin.Context) {
	cvID := ctx.Param("id")

	suggestions, err := c.cvUsecase.Analyze(ctx, cvID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCVID):
			ctx.JSON(http.StatusBadRequest, utils.ErrorPayload("Invalid CV ID", nil))
		case errors.Is(err, domain.ErrCVNotFound):
			ctx.JSON(http.StatusNotFound, utils.ErrorPayload("CV not found", nil))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.ErrorPayload("Failed to analyze CV", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessPayload("CV analyzed successfully", gin.H{
		"cvId":        cvID,
		"suggestions": suggestions,
	}))
}
