package handler

import (
	"bytes"
	"dia-backend/internal/app/repository"
	"dia-backend/internal/app/role"
	"encoding/json"
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

type AsyncUpdateRequestLampsRequest struct {
	Key                string `json:"key" binding:"required"`
	LightRequestToLamp []struct {
		LampID uint64 `json:"lamp_id" binding:"required"`
		Number uint64 `json:"number" binding:"required"`
	} `json:"light_request_to_lamp" binding:"required"`
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

	moderatorID, exists := GetUserIDFromContext(ctx)
	if !exists {
		panic("user not in context after passing auth check, something's wrong")
	}

	request, err := h.repo.LightRequest.GetLightRequestByID(id, moderatorID)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get request data"})
		return
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize request"})
		return
	}

	calcServiceURL := h.repo.Config.CalcService.URL
	if calcServiceURL == "" {
		logrus.Error("CalcService URL is not configured")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Calculation service not configured"})
		return
	}

	if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, 4); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := http.Post(calcServiceURL, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call calculation service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Calculation service returned status: %d", resp.StatusCode)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Calculation service failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request resolved successfully"})
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

// AsyncUpdateRequestLamp godoc
// @Summary      Asynchronously update lamps in request
// @Description  Update multiple lamps in a draft light request asynchronously using secret key
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Light Request ID"
// @Param        request body AsyncUpdateRequestLampsRequest true "Update lamps data"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-requests/{id}/async-update [put]
func (h *RequestHandler) AsyncUpdateRequestLamp(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req AsyncUpdateRequestLampsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if req.Key != h.repo.Config.CalcService.Key {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	lamps := make([]repository.AsyncUpdateLampRequest, len(req.LightRequestToLamp))

	for i, lamp := range req.LightRequestToLamp {
		lamps[i] = repository.AsyncUpdateLampRequest{
			LampID: lamp.LampID,
			Number: lamp.Number,
		}
	}

	if err := h.repo.LightRequest.AsyncUpdateRequestLampNumbers(id, lamps); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request lamps"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request lamps updated successfully"})
}
