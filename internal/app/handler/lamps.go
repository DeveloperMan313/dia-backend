package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampsHandler struct {
	LampRepository        *repository.LampRepository
	CalcRequestRepository *repository.CalcRequestRepository
}

func NewLampsHandler(lampRepo *repository.LampRepository, calcReqRepo *repository.CalcRequestRepository) *LampsHandler {
	return &LampsHandler{
		LampRepository:        lampRepo,
		CalcRequestRepository: calcReqRepo,
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

	searchQuery := ctx.Query("query")
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

	calcRequestID := 1
	calcRequestEntryCnt, err := h.CalcRequestRepository.GetCalcRequestEntryCntByID(calcRequestID)
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
			Name:        "query",
			Placeholder: "Поиск приборов",
			Value:       searchQuery,
		},
		"calcRequestID":       calcRequestID,
		"calcRequestEntryCnt": calcRequestEntryCnt,
	})
}
