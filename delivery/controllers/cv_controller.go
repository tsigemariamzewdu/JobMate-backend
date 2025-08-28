package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	usecase "github.com/tsigemariamzewdu/JobMate-backend/domain/interfaces/usecases"
)

type CVController struct {
	cvUsecase usecase.ICVUsecase
}

func NewCVController(u usecase.ICVUsecase) *CVController {
	return &CVController{cvUsecase: u}
}

type CVUploadRequest struct {
	UserID   string `json:"userId" binding:"required"`
	RawText  string `json:"rawText,omitempty"`
	FileName string `json:"fileName,omitempty"`
	File     *multipart.FileHeader
}

// POST /cv
func (c *CVController) UploadCV(ctx *gin.Context) {
	contentType := ctx.ContentType()
	var req CVUploadRequest

	if contentType == "application/json" {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}
	} else if contentType == "multipart/form-data" {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		req.File = file
		req.FileName = file.Filename
		req.UserID = ctx.PostForm("userId")
		if req.UserID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported content type"})
		return
	}

	createdCV, err := c.cvUsecase.Upload(ctx, req.UserID, req.RawText, req.File)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":        createdCV.ID,
		"userId":    createdCV.UserID,
		"fileName":  createdCV.FileName,
		"createdAt": createdCV.CreatedAt,
	})
}


// POST /cv/:id/analyze
func (c *CVController) AnalyzeCV(ctx *gin.Context) {
	cvID := ctx.Param("id")

	suggestions, err := c.cvUsecase.Analyze(ctx, cvID)
	if err != nil {
		if err == domain.ErrCVNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "CV not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cv_id":      cvID,
		"suggestions": suggestions,
	})
}