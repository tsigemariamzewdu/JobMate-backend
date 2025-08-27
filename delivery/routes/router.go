package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"
	"github.com/tsigemariamzewdu/JobMate-backend/domain"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth"
)

func SetupRouter(authMiddleware *auth.AuthMiddleware, uc *controllers.UserController, authUC domain.IAuthUsecase) *gin.Engine {
	router := gin.Default()

	authController := controllers.NewAuthController(authUC)
	registerUserRoutes(router, authMiddleware, uc, authController)
	return router
}

func registerUserRoutes(router *gin.Engine, authMiddleware *auth.AuthMiddleware, uc *controllers.UserController, authController *controllers.AuthController) {
	userRoutes := router.Group("/users", authMiddleware.Middleware())
	{
		userRoutes.GET("/me", uc.GetProfile)
		userRoutes.POST("/me", uc.UpdateProfile)
	}
	
	//auth routes
	router.POST("/auth/refresh", authController.RefreshToken)
}
