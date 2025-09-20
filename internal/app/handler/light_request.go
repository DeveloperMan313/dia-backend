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

type LightRequestHandler struct {
	repo *repository.Repository
}

func NewLightRequestHandler(repository *repository.Repository) *LightRequestHandler {
	return &LightRequestHandler{
		repo: repository,
	}
}

func (h *LightRequestHandler) Register(router *gin.Engine) {
	router.GET("/light-request/:id", h.GetLightRequestByID)
	router.POST("/light-request/:id", h.DeleteLightRequest)
}

func (h *LightRequestHandler) GetLightRequestByID(ctx *gin.Context) {
	lampIDStr := ctx.Param("id")
	reqID, err := strconv.ParseUint(lampIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	lightRequest, err := h.repo.LightRequest.GetLightRequestByID(reqID, 1)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "light_request.html", gin.H{
		"title":        "Просмотр заявки",
		"lightRequest": &lightRequest,
	})
}

func (h *LightRequestHandler) DeleteLightRequest(ctx *gin.Context) {
	requestIDStr := ctx.PostForm("request-id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusBadRequest)
		return
	}

	err = h.repo.LightRequest.DeleteLightRequest(requestID, 1)
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
