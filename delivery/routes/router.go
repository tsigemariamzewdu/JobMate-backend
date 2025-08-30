package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tsigemariamzewdu/JobMate-backend/delivery/controllers"
	"github.com/tsigemariamzewdu/JobMate-backend/infrastructure/auth"
)

func SetupRouter(authMiddleware *auth.AuthMiddleware,
	uc *controllers.UserController,
	authController *controllers.AuthController,
	otpController *controllers.OtpController,
	oauthController *controllers.OAuth2Controller,
	cvController *controllers.CVController,
	chatController *controllers.ChatController,
	jobController *controllers.JobController,
) *gin.Engine {

	router := gin.Default()

	// register user + auth routes
	registerUserRoutes(router, authMiddleware, uc, authController)

	// add OTP route
	otpRoutes := router.Group("/auth")
	{
		otpRoutes.POST("/request-otp", otpController.RequestOTP)
	}

	// Auth routes
	authGroup := router.Group("/auth")
	NewAuthRouter(*authController, *authGroup)

	RegisterOAuthRoutes(router, oauthController)

	// Chat routes
	chatRoutes := router.Group("/chat")
	{
		chatRoutes.POST("", chatController.SendMessage)
		chatRoutes.GET("/history", chatController.GetConversationHistory)
	}

	//cv routes
	cvGroup := router.Group("/cv")
	NewCVRouter(*cvController, *cvGroup)

	// Job suggestion route
	jobRoutes := router.Group("/jobs")
	{
		jobRoutes.POST("/suggest", jobController.SuggestJobs)
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

func NewAuthRouter(authController controllers.AuthController, group gin.RouterGroup) {

	group.POST("/register", authController.Register)
	group.POST("/login", authController.Login)
	group.POST("/logout", authController.Logout)
}

func NewCVRouter(cvController controllers.CVController, group gin.RouterGroup) {
	group.POST("/", cvController.UploadCV)
	group.POST("/:id/analye", cvController.AnalyzeCV)
}

func RegisterOAuthRoutes(
	router *gin.Engine,
	oauthController *controllers.OAuth2Controller,
) {
	oauth := router.Group("/oauth")
	{
		oauth.GET("/:provider/login", oauthController.RedirectToProvider)
		oauth.GET("/:provider/callback", oauthController.HandleCallback)
	}
}
