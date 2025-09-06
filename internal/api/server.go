package api

import (
	"dia-backend/internal/app/handler"
	"dia-backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	logrus.Debug("Server started")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./resources")

	r.GET("/lamps", handler.GetLamps)

	r.Run("0.0.0.0:8000")

	logrus.Debug("Server stopped")
}
