package handler

import (
	"dia-backend/internal/app/repository"
	"dia-backend/internal/app/role"
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

// GetCartInfo godoc
// @Summary      Get cart information
// @Description  Get draft request ID and item count for current user. For unauthorized/not found returns null values
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Success      200  {object}  CartInfoResponse  "Success. For unauthorized/empty cart: {\\"request_id\\": 0, \\"item_count\\": -1}"
// @Failure      500  {object}  map[string]string
// @Router       /light-requests/cart [get]
func (h *RequestHandler) GetCartInfo(ctx *gin.Context) {
	nullResponse := CartInfoResponse{
		RequestID: 0,
		ItemCount: -1,
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusOK, nullResponse)
		return
	}

	requestID, itemCount, err := h.repo.LightRequest.GetDraftRequestInfo(userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cart info"})
		return
	}

	if requestID == 0 {
		ctx.JSON(http.StatusOK, nullResponse)
		return
	}

	ctx.JSON(http.StatusOK, CartInfoResponse{
		RequestID: requestID,
		ItemCount: itemCount,
	})
}

// GetRequests godoc
// @Summary      List all light requests
// @Description  Get all light requests with optional filters
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        status     query  int     false  "Filter by status"
// @Param        date_from  query  string  false  "Filter from date (YYYY-MM-DD)"
// @Param        date_to    query  string  false  "Filter to date (YYYY-MM-DD)"
// @Success      200  {array}   ds.LightRequest
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-requests [get]
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

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	isMod := false
	userRole, exists := GetUserRoleFromContext(ctx)
	if exists && userRole == role.Moderator {
		isMod = true
	}

	requests, err := h.repo.LightRequest.GetLightRequests(statusFilter, dateFrom, dateTo, isMod, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get requests"})
		return
	}

	ctx.JSON(http.StatusOK, requests)
}

// GetRequestByID godoc
// @Summary      Get light request by ID
// @Description  Get detailed information about a specific light request
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Success      200  {object}  ds.LightRequest
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /light-requests/{id} [get]
func (h *RequestHandler) GetRequestByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	request, err := h.repo.LightRequest.GetLightRequestByID(id, userID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	ctx.JSON(http.StatusOK, request)
}

// UpdateRequest godoc
// @Summary      Update light request
// @Description  Update max total power for a draft light request
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Param        request body UpdateRequestRequest true "Update data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-requests/{id} [put]
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

// FormRequest godoc
// @Summary      Form light request
// @Description  Submit draft request for moderation
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /light-requests/{id}/form [put]
func (h *RequestHandler) FormRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.LightRequest.FormRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request formed successfully"})
}

// ResolveRequest godoc
// @Summary      Resolve light request
// @Description  Resolve light request with calculations (moderator only)
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /light-requests/{id}/resolve [put]
func (h *RequestHandler) ResolveRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	totalPower := h.repo.LightRequest.CalculateTotalPower(id)

	moderatorID, exists := GetUserIDFromContext(ctx)
	if !exists {
		panic("user not in context after passing auth check, something's wrong")
	}

	calculatedLamps := make(map[uint64]uint64)
	request, err := h.repo.LightRequest.GetLightRequestByID(id, moderatorID)

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

	if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, 4); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message": "Request resolved successfully",
		"calculated_data": gin.H{
			"total_power_w":    totalPower,
			"calculated_lamps": calculatedLamps,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

// RejectRequest godoc
// @Summary      Reject light request
// @Description  Reject light request
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /light-requests/{id}/reject [put]
func (h *RequestHandler) RejectRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	moderatorID, exists := GetUserIDFromContext(ctx)
	if !exists {
		panic("user not in context after passing auth check, something's wrong")
	}

	if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, 5); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
}

// DeleteRequest godoc
// @Summary      Delete light request
// @Description  Delete a draft light request
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-requests/{id} [delete]
func (h *RequestHandler) DeleteRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	moderatorID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.LightRequest.DeleteRequest(id, moderatorID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}
