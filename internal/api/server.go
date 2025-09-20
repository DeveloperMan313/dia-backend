package api

import (
	"dia-backend/internal/app/handler"
	"dia-backend/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	logrus.Debug("Server started")

	lampRepo, err := repository.NewLampRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория приборов")
	}
	lightReqRepo, err := repository.NewLightRequestRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория запросов")
	}

	lightRequestHandler := handler.NewLightRequestHandler(lampRepo, lightReqRepo)
	lampHandler := handler.NewLampHandler(lampRepo)
	lampsHandler := handler.NewLampsHandler(lampRepo, lightReqRepo)

	r := gin.Default()

	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./resources")

	r.GET("/light-request/:id", lightRequestHandler.GetLightRequestByID)
	r.GET("/lamp/:id", lampHandler.GetLampByID)
	r.GET("/lamps", lampsHandler.GetLamps)

	r.Run(":8000")

	logrus.Debug("Server stopped")
}
