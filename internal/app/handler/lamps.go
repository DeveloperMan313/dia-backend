package handler

import (
	"dia-backend/internal/app/ds"
	"dia-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampsHandler struct {
	repo *repository.Repository
}

func NewLampsHandler(repository *repository.Repository) *LampsHandler {
	return &LampsHandler{
		repo: repository,
	}
}

func (h *LampsHandler) Register(router *gin.Engine) {
	router.GET("/lamps", h.GetLamps)
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
	var lamps []ds.Lamp
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		lamps, err = h.repo.Lamp.GetLamps()
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	} else {
		lamps, err = h.repo.Lamp.GetLampsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
			ctx.Status(http.StatusNotFound)
			return
		}
	}

	calcRequestID := 1
	calcRequestEntryCnt, err := h.repo.CalcRequest.GetCalcRequestEntryCntByID(calcRequestID)
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
