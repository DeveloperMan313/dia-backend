package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LightRequestHandler struct {
	LampRepository         *repository.LampRepository
	LightRequestRepository *repository.LightRequestRepository
}

func NewLightRequestHandler(lampRepo *repository.LampRepository, lightReqRepo *repository.LightRequestRepository) *LightRequestHandler {
	return &LightRequestHandler{
		LampRepository:         lampRepo,
		LightRequestRepository: lightReqRepo,
	}
}

func (h *LightRequestHandler) GetLightRequestByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	reqID, err := strconv.Atoi(lampIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	lightReqView, err := h.LightRequestRepository.GetLightRequestViewByID(reqID, h.LampRepository)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "light_request.html", gin.H{
		"title": "Просмотр заявки",
		"view":  &lightReqView,
	})
}
