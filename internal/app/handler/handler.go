package handler

import (
	"dia-backend/internal/app/repository"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(router *gin.Engine, repo *repository.Repository) {
	router.LoadHTMLGlob("./templates/**/*")
	router.Static("/static", "./resources")

	calcRequestHandler := NewCalcRequestHandler(repo)
	lampHandler := NewLampHandler(repo)
	lampsHandler := NewLampsHandler(repo)

	calcRequestHandler.Register(router)
	lampHandler.Register(router)
	lampsHandler.Register(router)
}
