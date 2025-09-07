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
	calcReqRepo, err := repository.NewCalcRequestRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория запросов")
	}

	calcRequestHandler := handler.NewCalcRequestHandler(lampRepo, calcReqRepo)
	lampHandler := handler.NewLampHandler(lampRepo)
	lampsHandler := handler.NewLampsHandler(lampRepo, calcReqRepo)

	r := gin.Default()

	r.LoadHTMLGlob("./templates/**/*")
	r.Static("/static", "./resources")

	r.GET("/calc-request/:id", calcRequestHandler.GetCalcRequestByID)
	r.GET("/lamp/:id", lampHandler.GetLampByID)
	r.GET("/lamps", lampsHandler.GetLamps)

	r.Run(":8000")

	logrus.Debug("Server stopped")
}
