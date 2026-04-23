package handlers

import (
    "net/http"
    "time"

    "subscription-service/internal/models"
    "subscription-service/internal/repository"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

type SubscriptionHandler struct {
    repo *repository.SubscriptionRepository
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
    return &SubscriptionHandler{repo: repo}
}

// @Summary Create a new subscription
// @Description Create a new subscription record
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} models.SubscriptionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
    var req models.CreateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.WithError(err).Warn("Invalid request body for create subscription")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse dates
    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        log.WithError(err).Warn("Invalid start_date format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
        return
    }

    var endDate *time.Time
    if req.EndDate != nil && *req.EndDate != "" {
        parsed, err := time.Parse("01-2006", *req.EndDate)
        if err != nil {
            log.WithError(err).Warn("Invalid end_date format")
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
            return
        }
        endDate = &parsed
    }

    userID, err := uuid.Parse(req.UserID)
    if err != nil {
        log.WithError(err).Warn("Invalid user_id format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format, expected UUID"})
        return
    }

    subscription := &models.Subscription{
        ServiceName: req.ServiceName,
        Price:       req.Price,
        UserID:      userID,
        StartDate:   startDate,
        EndDate:     endDate,
    }

    if err := h.repo.Create(subscription); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
        return
    }

    response := mapToResponse(subscription)
    c.JSON(http.StatusCreated, response)
}

// @Summary Get subscription by ID
// @Description Get subscription details by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID (UUID)"
// @Success 200 {object} models.SubscriptionResponse
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        log.WithError(err).Warn("Invalid subscription ID format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
        return
    }

    subscription, err := h.repo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
        return
    }

    if subscription == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
        return
    }

    response := mapToResponse(subscription)
    c.JSON(http.StatusOK, response)
}

// @Summary Update subscription
// @Description Update subscription details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID (UUID)"
// @Param subscription body models.UpdateSubscriptionRequest true "Update data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        log.WithError(err).Warn("Invalid subscription ID format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
        return
    }

    var req models.UpdateSubscriptionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.WithError(err).Warn("Invalid request body for update subscription")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.repo.Update(id, &req); err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "subscription updated successfully"})
}

// @Summary Delete subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        log.WithError(err).Warn("Invalid subscription ID format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id format"})
        return
    }

    if err := h.repo.Delete(id); err != nil {
        if err.Error() == "subscription not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "subscription deleted successfully"})
}

// @Summary Get total cost of subscriptions
// @Description Calculate total cost of active subscriptions within date range with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.TotalCostRequest true "Filter criteria"
// @Success 200 {object} models.TotalCostResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subscriptions/total-cost [post]
func (h *SubscriptionHandler) GetTotalCost(c *gin.Context) {
    var req models.TotalCostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.WithError(err).Warn("Invalid request body for total cost calculation")
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    startDate, err := time.Parse("01-2006", req.StartDate)
    if err != nil {
        log.WithError(err).Warn("Invalid start_date format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, expected MM-YYYY"})
        return
    }

    endDate, err := time.Parse("01-2006", req.EndDate)
    if err != nil {
        log.WithError(err).Warn("Invalid end_date format")
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, expected MM-YYYY"})
        return
    }

    total, err := h.repo.GetTotalCost(startDate, endDate, req.UserID, req.ServiceName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total cost"})
        return
    }

    c.JSON(http.StatusOK, models.TotalCostResponse{TotalCost: total})
}

func mapToResponse(sub *models.Subscription) models.SubscriptionResponse {
    var endDateStr *string
    if sub.EndDate != nil {
        str := sub.EndDate.Format("01-2006")
        endDateStr = &str
    }

    return models.SubscriptionResponse{
        ID:          sub.ID.String(),
        ServiceName: sub.ServiceName,
        Price:       sub.Price,
        UserID:      sub.UserID.String(),
        StartDate:   sub.StartDate.Format("01-2006"),
        EndDate:     endDateStr,
        CreatedAt:   sub.CreatedAt.Format(time.RFC3339),
        UpdatedAt:   sub.UpdatedAt.Format(time.RFC3339),
    }
}