package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CalcRequestHandler struct {
	repo *repository.Repository
}

func NewCalcRequestHandler(repository *repository.Repository) *CalcRequestHandler {
	return &CalcRequestHandler{
		repo: repository,
	}
}

func (h *CalcRequestHandler) Register(router *gin.Engine) {
	router.GET("/calc-request/:id", h.GetCalcRequestByID)
}

func (h *CalcRequestHandler) GetCalcRequestByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	reqID, err := strconv.Atoi(lampIDStr)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	calcReqView, err := h.repo.CalcRequest.GetCalcRequestViewByID(reqID, h.repo.Lamp)
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
