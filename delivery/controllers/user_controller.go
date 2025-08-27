package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/dto"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
)

type UserController struct {
	userUsecase domain.IUserUsecase
}

func NewUserController(userUsecase domain.IUserUsecase) *UserController {
	return &UserController{userUsecase: userUsecase}
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateUserProfile
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request payload",
		})
		return
	}

	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	user := &domain.User{
		UserID:            userID,
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		PreferredLanguage: req.PreferredLanguage,
		EducationLevel:    req.EducationLevel,
		FieldOfStudy:      req.FieldOfStudy,
		YearsExperience:   req.YearsExperience,
		CareerInterests:   req.CareerInterests,
		CareerGoals:       req.CareerGoals,
		ProfilePicture:    req.ProfilePicture,
	}

	updatedUser, err := c.userUsecase.UpdateProfile(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update profile",
		})
		return
	}

	resp := dto.UserProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		User:    dto.ToUserProfileView(updatedUser),
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Unauthorized",
		})
		return
	}

	user, err := c.userUsecase.GetProfile(ctx, userID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch user profile",
		})
		return
	}

	resp := dto.UserProfileResponse{
		Success: true,
		Message: "User profile retrieved successfully",
		User:    dto.ToUserProfileView(user),
	}

	ctx.JSON(http.StatusOK, resp)
}
