package handler

import (
	"dia-backend/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	apiRouter := router.Group("/api")

	lampHandler := NewLampHandler(repo)
	lampRouter := apiRouter.Group("/lamps")
	{
		lampRouter.GET("", lampHandler.GetLamps)
		lampRouter.GET("/:id", lampHandler.GetLampByID)
		lampRouter.POST("", lampHandler.CreateLamp)
		lampRouter.PUT("/:id", lampHandler.UpdateLamp)
		lampRouter.DELETE("/:id", lampHandler.DeleteLamp)
		lampRouter.POST("/:id/image", lampHandler.AddLampImage)
		lampRouter.POST("/:id/draft", lampHandler.AddToDraftRequest)
	}

	requestHandler := NewRequestHandler(repo)
	requestRouter := apiRouter.Group("/light-requests")
	{
		requestRouter.GET("/cart", requestHandler.GetCartInfo)
		requestRouter.GET("", requestHandler.GetRequests)
		requestRouter.GET("/:id", requestHandler.GetRequestByID)
		requestRouter.PUT("/:id", requestHandler.UpdateRequest)
		requestRouter.PUT("/:id/form", requestHandler.FormRequest)
		requestRouter.PUT("/:id/resolve", requestHandler.ResolveRequest)
		requestRouter.PUT("/:id/reject", requestHandler.RejectRequest)
		requestRouter.DELETE("/:id", requestHandler.DeleteRequest)
	}

	requestLampHandler := NewRequestLampHandler(repo)
	requestLampRouter := apiRouter.Group("/light-request-lamps")
	{
		requestLampRouter.DELETE("", requestLampHandler.RemoveFromRequest)
		requestLampRouter.PUT("", requestLampHandler.UpdateRequestLamp)
	}

	userHandler := NewUserHandler(repo)
	userRouter := apiRouter.Group("/users")
	{
		userRouter.POST("/register", userHandler.Register)
		userRouter.GET("/profile", userHandler.GetProfile)
		userRouter.PUT("/profile", userHandler.UpdateProfile)
		userRouter.POST("/login", userHandler.Login)
		userRouter.POST("/logout", userHandler.Logout)
	}
}
