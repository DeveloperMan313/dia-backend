package handler

import (
	"dia-backend/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RequestHandler struct {
	repo *repository.Repository
}

func NewRequestHandler(repo *repository.Repository) *RequestHandler {
	return &RequestHandler{
		repo: repo,
	}
}

type CartInfoResponse struct {
	RequestID uint64 `json:"request_id"`
	ItemCount int    `json:"item_count"`
}

type UpdateRequestRequest struct {
	MaxTotalPowerW *float64 `json:"max_total_power_w"`
}

func (h *RequestHandler) GetCartInfo(ctx *gin.Context) {
	userID := GetFixedUserID()
	requestID, itemCount, err := h.repo.LightRequest.GetDraftRequestInfo(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart info"})
		return
	}

	if requestID == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	ctx.JSON(http.StatusOK, CartInfoResponse{
		RequestID: requestID,
		ItemCount: itemCount,
	})
}

func (h *RequestHandler) GetRequests(ctx *gin.Context) {
	var statusFilter uint8
	if statusStr := ctx.Query("status"); statusStr != "" {
		if status, err := strconv.ParseUint(statusStr, 10, 8); err == nil {
			statusFilter = uint8(status)
		}
	}

	var dateFrom, dateTo *time.Time
	if dateFromStr := ctx.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &parsed
		}
	}
	if dateToStr := ctx.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &parsed
		}
	}

	requests, err := h.repo.LightRequest.GetLightRequests(statusFilter, dateFrom, dateTo)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) GetRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID := GetFixedUserID()
	request, err := h.repo.LightRequest.GetLightRequestByID(id, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	ctx.JSON(http.StatusOK, request)
}

func (h *RequestHandler) UpdateRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req UpdateRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := h.repo.LightRequest.UpdateLightRequest(id, req.MaxTotalPowerW); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request updated successfully"})
}

func (h *RequestHandler) FormRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.LightRequest.FormRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request formed successfully"})
}

func (h *RequestHandler) ResolveRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	totalPower := h.repo.LightRequest.CalculateTotalPower(id)

	deliveryDate := time.Now().AddDate(0, 1, 0)

	calculatedLamps := make(map[uint64]uint64)
	request, err := h.repo.LightRequest.GetLightRequestByID(id, GetFixedUserID())
	if err == nil {
		for _, entry := range request.LightRequestToLamp {
			// N = (E * S) / Phi
			// E = 500 lux (standard office lighting)
			// S = area_m2
			// Phi = luminous_flux_lm
			requiredIlluminationLux := 500.0
			calculatedNumber := (requiredIlluminationLux * entry.AreaM2) / entry.Lamp.LuminousFluxLm
			calculatedLamps[entry.LampID] = uint64(calculatedNumber)
		}
	}

	moderatorID := uint64(2)
	if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, 4); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message": "Request resolved successfully",
		"calculated_data": gin.H{
			"total_power_w":    totalPower,
			"delivery_date":    deliveryDate.Format("2006-01-02"),
			"calculated_lamps": calculatedLamps,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *RequestHandler) RejectRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	moderatorID := uint64(2)
	if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, 5); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
}

func (h *RequestHandler) DeleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID := GetFixedUserID()
	if err := h.repo.LightRequest.DeleteRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}
