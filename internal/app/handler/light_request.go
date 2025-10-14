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

type FinishRequestRequest struct {
	Action string `json:"action" binding:"required"`
}

// GetCartInfo godoc
// @Summary      Get cart information
// @Description  Get draft request ID and item count for current user
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  CartInfoResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /light-requests/cart [get]
func (h *RequestHandler) GetCartInfo(ctx *gin.Context) {
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

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

// FinishRequest godoc
// @Summary      Finish light request
// @Description  Resolve or reject a light request with calculations (moderator only)
// @Tags         light-requests
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Light Request ID"
// @Param        request body FinishRequestRequest true "Action data"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /light-requests/{id}/finish [put]
func (h *RequestHandler) FinishRequest(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}

	var req FinishRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if req.Action != "resolve" && req.Action != "reject" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Action must be 'resolve' or 'reject'"})
		return
	}

	var status uint8
	if req.Action == "resolve" {
		status = 4
	} else {
		status = 5
	}

	moderatorID := uint64(2)
	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if req.Action == "resolve" {
		request, err := h.repo.LightRequest.GetLightRequestByID(id, userID)
		if err == nil {
			for _, entry := range request.LightRequestToLamp {
				// N = (E * S) / Phi
				// E = 500 lux (standard office lighting)
				// S = area_m2
				// Phi = luminous_flux_lm
				requiredIlluminationLux := 500.0
				calculatedNumber := (requiredIlluminationLux * entry.AreaM2) / entry.Lamp.LuminousFluxLm
				number := uint64(calculatedNumber)

				if err := h.repo.LightRequest.UpdateRequestToLamp(id, entry.LampID, userID, nil, &number); err != nil {
					logrus.Errorf("Failed to update lamp number for lamp %d: %v", entry.LampID, err)
				}
			}
		}

		if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, status); err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Request resolved successfully"})
	} else {
		if err := h.repo.LightRequest.ResolveOrRejectRequest(id, moderatorID, status); err != nil {
			logrus.Error(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Request rejected successfully"})
	}
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

	userID, exists := GetUserIDFromContext(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.repo.LightRequest.DeleteRequest(id, userID); err != nil {
		logrus.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Request deleted successfully"})
}
