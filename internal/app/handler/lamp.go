package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampHandler struct {
	LampRepository *repository.LampRepository
}

func NewLampHandler(lampRepo *repository.LampRepository) *LampHandler {
	return &LampHandler{
		LampRepository: lampRepo,
	}
}

func (h *LampHandler) GetLampByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	lampID, err := strconv.Atoi(lampIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	lamp, err := h.LampRepository.GetLampByID(lampID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "lamp.html", gin.H{
		"title": "О приборе",
		"lamp":  *lamp,
	})
}
