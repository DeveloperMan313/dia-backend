package handler

import (
	"dia-backend/internal/app/ds"
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

type LightRequestTemplateEntry struct {
	Lamp   ds.Lamp
	AreaM2 CompTextInput
	Number uint64
}

func NewLightRequestTemplateEntry(lightReqEntry *ds.LightRequestToLamp) *LightRequestTemplateEntry {
	return &LightRequestTemplateEntry{
		Lamp: lightReqEntry.Lamp,
		AreaM2: CompTextInput{
			ShowLabel:   true,
			Label:       "Площадь, м²",
			Type:        "text",
			Name:        "area-m2",
			Value:       strconv.FormatFloat(lightReqEntry.AreaM2, 'f', 2, 64),
			Placeholder: "Площадь",
		},
		Number: lightReqEntry.Number,
	}
}

type LightRequestTemplate struct {
	ID             uint64
	MaxTotalPowerW CompTextInput
	Entries        []LightRequestTemplateEntry
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

	lightReqTemplate := LightRequestTemplate{
		ID: lightRequest.ID,
		MaxTotalPowerW: CompTextInput{
			ShowLabel:   true,
			Label:       "Суммарная мощность, вт",
			Type:        "text",
			Name:        "max-total-power-w",
			Value:       strconv.FormatFloat(lightRequest.MaxTotalPowerW, 'f', 2, 64),
			Placeholder: "Мощность",
		},
	}

	for _, lightReqToLamp := range lightRequest.LightRequestToLamp {
		lightReqTemplate.Entries = append(lightReqTemplate.Entries, *NewLightRequestTemplateEntry(&lightReqToLamp))
	}

	ctx.HTML(http.StatusOK, "light_request.html", gin.H{
		"title":        "Составление заявки",
		"lightRequest": &lightReqTemplate,
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

	err = h.repo.LightRequest.DeleteRequest(requestID, 1)
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
