package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampHandler struct {
	repo *repository.Repository
}

func NewLampHandler(repository *repository.Repository) *LampHandler {
	return &LampHandler{
		repo: repository,
	}
}

func (h *LampHandler) Register(router *gin.Engine) {
	router.GET("/lamp/:id", h.GetLampByID)
}

func (h *LampHandler) GetLampByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	lampID, err := strconv.ParseUint(lampIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	lamp, err := h.repo.Lamp.GetLampByID(lampID)
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
