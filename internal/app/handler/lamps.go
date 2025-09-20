package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampsHandler struct {
	LampRepository         *repository.LampRepository
	LightRequestRepository *repository.LightRequestRepository
}

func NewLampsHandler(lampRepo *repository.LampRepository, lightReqRepo *repository.LightRequestRepository) *LampsHandler {
	return &LampsHandler{
		LampRepository:         lampRepo,
		LightRequestRepository: lightReqRepo,
	}
}

type CompTextInput struct {
	ShowLabel   bool
	Label       string
	Type        string
	Name        string
	Value       string
	Placeholder string
}

func (h *LampsHandler) GetLamps(ctx *gin.Context) {
	var lamps []repository.Lamp
	var err error

	searchQuery := ctx.Query("title")
	if searchQuery == "" {
		lamps, err = h.LampRepository.GetLamps()
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	} else {
		lamps, err = h.LampRepository.GetLampsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	lightRequestID := 1
	lightRequestEntryCnt, err := h.LightRequestRepository.GetLightRequestEntryCntByID(lightRequestID)
	if err != nil {
		logrus.Error(err)
		ctx.Status(http.StatusNotFound)
		return
	}

	ctx.HTML(http.StatusOK, "lamps.html", gin.H{
		"title": "Приборы",
		"lamps": lamps,
		"search": CompTextInput{
			ShowLabel:   false,
			Type:        "text",
			Name:        "title",
			Placeholder: "Поиск приборов",
			Value:       searchQuery,
		},
		"lightRequestID":       lightRequestID,
		"lightRequestEntryCnt": lightRequestEntryCnt,
	})
}
