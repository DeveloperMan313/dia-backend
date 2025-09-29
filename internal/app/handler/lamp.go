package handler

import (
	"dia-backend/internal/app/ds"
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LampHandler struct {
	repo *repository.Repository
}

func NewLampHandler(repo *repository.Repository) *LampHandler {
	return &LampHandler{
		repo: repo,
	}
}

type CreateLampRequest struct {
	Title              string  `json:"title" binding:"required"`
	PowerW             float64 `json:"power_w" binding:"required"`
	LuminousFluxLm     float64 `json:"luminous_flux_lm" binding:"required"`
	ScatteringAngleDeg float64 `json:"scattering_angle_deg" binding:"required"`
}

type UpdateLampRequest struct {
	Title              *string  `json:"title"`
	PowerW             *float64 `json:"power_w"`
	LuminousFluxLm     *float64 `json:"luminous_flux_lm"`
	ScatteringAngleDeg *float64 `json:"scattering_angle_deg"`
}

func (h *LampHandler) GetLamps(ctx *gin.Context) {
	searchQuery := ctx.Query("title")

	lamps, err := h.repo.Lamp.GetLamps(searchQuery)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get lamps"})
		return
	}

	ctx.JSON(http.StatusOK, lamps)
}

func (h *LampHandler) GetLampByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	lamp, err := h.repo.Lamp.GetLampByID(id)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Lamp not found"})
		return
	}

	ctx.JSON(http.StatusOK, lamp)
}

func (h *LampHandler) CreateLamp(ctx *gin.Context) {
	var req CreateLampRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	lamp := &ds.Lamp{
		Title:              req.Title,
		PowerW:             req.PowerW,
		LuminousFluxLm:     req.LuminousFluxLm,
		ScatteringAngleDeg: req.ScatteringAngleDeg,
	}

	if err := h.repo.Lamp.CreateLamp(lamp); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lamp"})
		return
	}

	ctx.JSON(http.StatusCreated, lamp)
}

func (h *LampHandler) UpdateLamp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	var req UpdateLampRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	lampData := &ds.Lamp{}
	if req.Title != nil {
		lampData.Title = *req.Title
	}
	if req.PowerW != nil {
		lampData.PowerW = *req.PowerW
	}
	if req.LuminousFluxLm != nil {
		lampData.LuminousFluxLm = *req.LuminousFluxLm
	}
	if req.ScatteringAngleDeg != nil {
		lampData.ScatteringAngleDeg = *req.ScatteringAngleDeg
	}

	if err := h.repo.Lamp.UpdateLamp(id, lampData); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lamp"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lamp updated successfully"})
}

func (h *LampHandler) DeleteLamp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	if err := h.repo.Lamp.DeleteLamp(id); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lamp"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lamp deleted successfully"})
}

func (h *LampHandler) AddLampImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	if err := h.repo.Lamp.AddLampImage(id, file); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add image"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image added successfully"})
}

func (h *LampHandler) AddToDraftRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.Lamp.AddLampToDraftRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add lamp to draft request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lamp added to draft request successfully"})
}
