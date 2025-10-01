package api

import (
	"dia-backend/internal/app/handler"
	"dia-backend/internal/app/repository"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	repo, err := repository.NewRepository()
	if err != nil {
		panic(err)
	}

	defer repository.CloseDBConn(repo)

	router := gin.Default()

	handler.RegisterHandlers(router, repo)

	logrus.Debug("Server started")

	router.Run(":" + os.Getenv("LISTEN_PORT"))

	logrus.Debug("Server stopped")
}
