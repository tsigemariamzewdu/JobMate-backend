package routers

import (
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authController *controllers.AuthController)  *gin.Engine{
	router := gin.Default()
	authRoutes :=router.Group("/auth")
	{
		authRoutes.POST("/request-otp", authController.RequestOTP)
	}	
	return router
}