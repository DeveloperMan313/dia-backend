package handler

import (
	"dia-backend/internal/app/repository"
	"dia-backend/internal/app/role"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api")

	baseHandler := NewBaseHandler(repo)

	userHandler := NewUserHandler(repo)
	userRouter := apiRouter.Group("/users")
	{
		userRouter.POST("/register", userHandler.Register)
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/logout", baseHandler.WithAuthCheck(role.User, role.Moderator), userHandler.Logout)
		userRouter.GET("/profile", baseHandler.WithAuthCheck(role.User, role.Moderator), userHandler.GetProfile)
		userRouter.PUT("/profile", baseHandler.WithAuthCheck(role.User, role.Moderator), userHandler.UpdateProfile)
	}

	lampHandler := NewLampHandler(repo)
	lampRouter := apiRouter.Group("/lamps")
	{
		lampRouter.GET("", lampHandler.GetLamps)
		lampRouter.GET("/:id", lampHandler.GetLampByID)

		lampRouter.POST("", baseHandler.WithAuthCheck(role.Moderator), lampHandler.CreateLamp)
		lampRouter.PUT("/:id", baseHandler.WithAuthCheck(role.Moderator), lampHandler.UpdateLamp)
		lampRouter.DELETE("/:id", baseHandler.WithAuthCheck(role.Moderator), lampHandler.DeleteLamp)
		lampRouter.POST("/:id/image", baseHandler.WithAuthCheck(role.Moderator), lampHandler.AddLampImage)

		lampRouter.POST("/:id/draft", baseHandler.WithAuthCheck(role.User, role.Moderator), lampHandler.AddToDraftRequest)
	}

	requestHandler := NewRequestHandler(repo)
	requestRouter := apiRouter.Group("/light-requests")
	{
		requestRouter.GET("/cart", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.GetCartInfo)
		requestRouter.GET("/:id", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.GetRequestByID)
		requestRouter.PUT("/:id", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.UpdateRequest)
		requestRouter.PUT("/:id/form", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.FormRequest)
		requestRouter.DELETE("/:id", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.DeleteRequest)
		requestRouter.GET("", baseHandler.WithAuthCheck(role.User, role.Moderator), requestHandler.GetRequests)

		requestRouter.PUT("/:id/finish", baseHandler.WithAuthCheck(role.Moderator), requestHandler.FinishRequest)
	}

	requestLampHandler := NewRequestLampHandler(repo)
	requestLampRouter := apiRouter.Group("/light-request-lamps")
	requestLampRouter.Use(baseHandler.WithAuthCheck(role.User, role.Moderator))
	{
		requestLampRouter.DELETE("", requestLampHandler.RemoveFromRequest)
		requestLampRouter.PUT("", requestLampHandler.UpdateRequestLamp)
	}
}
