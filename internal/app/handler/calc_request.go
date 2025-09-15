package handler

import (
	"dia-backend/internal/app/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	router.POST("/calc-request/:id", h.DeleteCalcRequest)
}

func (h *CalcRequestHandler) GetCalcRequestByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	reqID, err := strconv.ParseUint(lampIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	calcRequest, err := h.repo.CalcRequest.GetCalcRequestByID(reqID, 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "calc_request.html", gin.H{
		"title":       "Просмотр заявки",
		"calcRequest": &calcRequest,
	})
}

func (h *CalcRequestHandler) DeleteCalcRequest(ctx *gin.Context) {
	requestIDStr := ctx.PostForm("request-id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.repo.CalcRequest.DeleteCalcRequest(requestID, 1)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/lamps")
}
