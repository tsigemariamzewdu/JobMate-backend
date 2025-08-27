package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth"
)

func SetupRouter(authMiddleware *auth.AuthMiddleware, uc *controllers.UserController, authUC domain.IAuthUsecase) *gin.Engine {
	router := gin.Default()

	// auth controller from domain usecase
	authController := controllers.NewAuthController(authUC)

	// register user + auth routes
	registerUserRoutes(router, authMiddleware, uc, authController)

	// add OTP route 
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/request-otp", authController.RequestOTP)
	}

	return router
}

func registerUserRoutes(router *gin.Engine, authMiddleware *auth.AuthMiddleware, uc *controllers.UserController, authController *controllers.AuthController) {
	userRoutes := router.Group("/users", authMiddleware.Middleware())
	{
		userRoutes.GET("/me", uc.GetProfile)
		userRoutes.POST("/me", uc.UpdateProfile)
	}
	
	// refresh token
	router.POST("/auth/refresh", authController.RefreshToken)
}
