package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SubscriptionHandler 訂閱 API handler
type SubscriptionHandler struct {
	service service.SubscriptionService
}

// NewSubscriptionHandler 建立新的訂閱 handler
func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription 建立新的訂閱
// @Summary 建立訂閱
// @Description 建立新的訂閱
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.CreateSubscriptionInput true "訂閱資料"
// @Success 201 {object} APIResponse{data=models.Subscription}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var input models.CreateSubscriptionInput

	// 綁定並驗證請求資料
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	// 呼叫 service 建立訂閱
	subscription, err := h.service.CreateSubscription(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "CREATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Data: subscription,
	})
}

// GetSubscription 取得單筆訂閱
// @Summary 取得訂閱
// @Description 根據 ID 取得單筆訂閱
// @Tags subscriptions
// @Produce json
// @Param id path string true "訂閱 ID"
// @Success 200 {object} APIResponse{data=models.Subscription}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid subscription ID format",
			},
		})
		return
	}

	// 呼叫 service 取得訂閱
	subscription, err := h.service.GetSubscription(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Error: &APIError{
				Code:    "NOT_FOUND",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: subscription,
	})
}

// ListSubscriptions 取得訂閱列表
// @Summary 取得訂閱列表
// @Description 取得訂閱列表，支援篩選和分頁
// @Tags subscriptions
// @Produce json
// @Param status query string false "狀態篩選 (active, cancelled)"
// @Param billing_cycle query string false "計費週期篩選 (monthly, quarterly, yearly)"
// @Param limit query int false "每頁筆數" default(100)
// @Param offset query int false "略過筆數" default(0)
// @Success 200 {object} APIResponse{data=[]models.Subscription}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	// 解析查詢參數
	filters := repository.SubscriptionFilters{}

	// 狀態篩選
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.SubscriptionStatus(statusStr)
		filters.Status = &status
	}

	// 分頁參數
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// 呼叫 service 取得訂閱列表
	subscriptions, err := h.service.ListSubscriptions(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "LIST_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: subscriptions,
	})
}

// UpdateSubscription 更新訂閱
// @Summary 更新訂閱
// @Description 更新訂閱資料
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "訂閱 ID"
// @Param subscription body models.UpdateSubscriptionInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.Subscription}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid subscription ID format",
			},
		})
		return
	}

	var input models.UpdateSubscriptionInput

	// 綁定並驗證請求資料
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	// 呼叫 service 更新訂閱
	subscription, err := h.service.UpdateSubscription(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "UPDATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: subscription,
	})
}

// DeleteSubscription 刪除訂閱
// @Summary 刪除訂閱
// @Description 刪除訂閱
// @Tags subscriptions
// @Produce json
// @Param id path string true "訂閱 ID"
// @Success 200 {object} APIResponse{data=string}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid subscription ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除訂閱
	if err := h.service.DeleteSubscription(id); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Subscription deleted successfully",
	})
}

// CancelSubscription 取消訂閱
// @Summary 取消訂閱
// @Description 取消訂閱（設定結束日期）
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "訂閱 ID"
// @Param input body CancelSubscriptionInput true "取消訂閱資料"
// @Success 200 {object} APIResponse{data=models.Subscription}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/subscriptions/{id}/cancel [post]
func (h *SubscriptionHandler) CancelSubscription(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid subscription ID format",
			},
		})
		return
	}

	var input CancelSubscriptionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	// 呼叫 service 取消訂閱
	subscription, err := h.service.CancelSubscription(id, input.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "CANCEL_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: subscription,
	})
}

// CancelSubscriptionInput 取消訂閱的輸入
type CancelSubscriptionInput struct {
	EndDate time.Time `json:"end_date" binding:"required"`
}

