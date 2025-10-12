package api

import (
	"dia-backend/internal/app/handler"
	"dia-backend/internal/app/repository"
	"os"

	_ "dia-backend/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer() {
	repo, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}

	defer repository.CloseDBConn(repo)

	router := gin.Default()

	handler.RegisterHandlers(router, repo)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.Debug("Server started")

	router.Run(":" + os.Getenv("LISTEN_PORT"))

	logrus.Debug("Server stopped")
}
