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

// @Summary      List lamps
// @Description  Get list of all lamps with optional title filter
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Param        title  query  string  false  "Filter by lamp title"
// @Success      200  {array}   ds.Lamp
// @Failure      500  {object}  map[string]string
// @Router       /lamps [get]
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

// @Summary      Get lamp by ID
// @Description  Get detailed information about a specific lamp
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Lamp ID"
// @Success      200  {object}  ds.Lamp
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /lamps/{id} [get]
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

// @Summary      Create a new lamp
// @Description  Create a new lamp (moderator only)
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateLampRequest true "Lamp data"
// @Success      201  {object}  ds.Lamp
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lamps [post]
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

// @Summary      Update lamp
// @Description  Update lamp information (moderator only)
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Lamp ID"
// @Param        request body UpdateLampRequest true "Lamp update data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lamps/{id} [put]
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

// @Summary      Delete lamp
// @Description  Soft delete a lamp (moderator only)
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Lamp ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lamps/{id} [delete]
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

// @Summary      Add lamp image
// @Description  Upload or replace lamp image (moderator only)
// @Tags         lamps
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Lamp ID"
// @Param        image formData file true "Lamp image file"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lamps/{id}/image [post]
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

// @Summary      Add lamp to draft request
// @Description  Add lamp to user's draft light request
// @Tags         lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Lamp ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /lamps/{id}/draft [post]
func (h *LampHandler) AddToDraftRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lamp ID"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.Lamp.AddLampToDraftRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add lamp to draft request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lamp added to draft request successfully"})
}
