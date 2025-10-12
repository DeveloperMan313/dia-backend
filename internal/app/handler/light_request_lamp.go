package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestLampHandler struct {
	repo *repository.Repository
}

func NewRequestLampHandler(repo *repository.Repository) *RequestLampHandler {
	return &RequestLampHandler{
		repo: repo,
	}
}

type RemoveFromRequestRequest struct {
	RequestID uint64 `json:"request_id" binding:"required"`
	LampID    uint64 `json:"lamp_id" binding:"required"`
}

type UpdateRequestLampRequest struct {
	RequestID uint64   `json:"request_id" binding:"required"`
	LampID    uint64   `json:"lamp_id" binding:"required"`
	AreaM2    *float64 `json:"area_m2"`
	Number    *uint64  `json:"number"`
}

// @Summary      Remove lamp from request
// @Description  Remove a lamp from a draft light request
// @Tags         light-request-lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body RemoveFromRequestRequest true "Remove lamp data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-request-lamps [delete]
func (h *RequestLampHandler) RemoveFromRequest(ctx *gin.Context) {
	var req RemoveFromRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.LightRequest.RemoveLampFromRequest(req.RequestID, req.LampID, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove lamp from request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Lamp removed from request successfully"})
}

// @Summary      Update lamp in request
// @Description  Update area and/or number of lamps in a draft request
// @Tags         light-request-lamps
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body UpdateRequestLampRequest true "Update lamp data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-request-lamps [put]
func (h *RequestLampHandler) UpdateRequestLamp(ctx *gin.Context) {
	var req UpdateRequestLampRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.LightRequest.UpdateRequestToLamp(req.RequestID, req.LampID, userID, req.AreaM2, req.Number); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request-lamp"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request-lamp updated successfully"})
}
