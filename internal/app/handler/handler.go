package handler

import (
	"dia-backend/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	router.LoadHTMLGlob("./templates/**/*")
	router.Static("/static", "./resources")

	lightRequestHandler := NewLightRequestHandler(repo)
	lampHandler := NewLampHandler(repo)
	lampsHandler := NewLampsHandler(repo)

	lightRequestHandler.Register(router)
	lampHandler.Register(router)
	lampsHandler.Register(router)
}
