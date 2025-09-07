package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CalcRequestHandler struct {
	LampRepository        *repository.LampRepository
	CalcRequestRepository *repository.CalcRequestRepository
}

func NewCalcRequestHandler(lampRepo *repository.LampRepository, calcReqRepo *repository.CalcRequestRepository) *CalcRequestHandler {
	return &CalcRequestHandler{
		LampRepository:        lampRepo,
		CalcRequestRepository: calcReqRepo,
	}
}

func (h *CalcRequestHandler) GetCalcRequestByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	reqID, err := strconv.Atoi(lampIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	calcReqView, err := h.CalcRequestRepository.GetCalcRequestViewByID(reqID, h.LampRepository)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "calc_request.html", gin.H{
		"title": "Просмотр заявки",
		"view":  &calcReqView,
	})
}
