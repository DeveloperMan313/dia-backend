package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
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

func (h *Handler) GetLamps(ctx *gin.Context) {
	var lamps []repository.Lamp
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		lamps, err = h.Repository.GetLamps()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		lamps, err = h.Repository.GetLampsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "lamps.html", gin.H{
		"lamps": lamps,
		"search": CompTextInput{
			ShowLabel:   false,
			Type:        "text",
			Name:        "query",
			Placeholder: "Поиск приборов",
			Value:       searchQuery,
		},
	})
}
