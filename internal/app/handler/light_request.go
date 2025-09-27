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

type LightRequestTemplateEntry struct {
	Lamp   repository.Lamp
	AreaM2 CompTextInput
	Number int
}

func NewLightRequestTemplateEntry(lightReqEntry *repository.LightRequestViewEntry) *LightRequestTemplateEntry {
	return &LightRequestTemplateEntry{
		Lamp: lightReqEntry.Lamp,
		AreaM2: CompTextInput{
			ShowLabel:   true,
			Label:       "Площадь, м²",
			Type:        "text",
			Name:        "area-m2",
			Value:       strconv.FormatFloat(float64(lightReqEntry.AreaM2), 'f', 2, 32),
			Placeholder: "Площадь",
		},
		Number: lightReqEntry.Number,
	}
}

type LightRequestTemplate struct {
	TotalPowerW CompTextInput
	Entries     []LightRequestTemplateEntry
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

	lightReqTemplate := LightRequestTemplate{
		TotalPowerW: CompTextInput{
			ShowLabel:   true,
			Label:       "Суммарная мощность, вт",
			Type:        "text",
			Name:        "max-total-power-w",
			Value:       strconv.FormatFloat(float64(lightReqView.LightRequest.MaxTotalPowerW), 'f', 2, 32),
			Placeholder: "Мощность",
		},
	}

	for _, item := range lightReqView.Entries {
		lightReqTemplate.Entries = append(lightReqTemplate.Entries, *NewLightRequestTemplateEntry(&item))
	}

	ctx.HTML(http.StatusOK, "light_request.html", gin.H{
		"title": "Составление заявки",
		"view":  &lightReqTemplate,
	})
}
