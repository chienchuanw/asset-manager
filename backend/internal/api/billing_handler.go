package api

import (
	"net/http"
	"time"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// BillingHandler 扣款 API handler
type BillingHandler struct {
	billingService service.BillingService
}

// NewBillingHandler 建立新的扣款 handler
func NewBillingHandler(billingService service.BillingService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
	}
}

// ProcessDailyBilling 處理每日扣款（手動觸發）
// @Summary 處理每日扣款
// @Description 手動觸發每日扣款處理（訂閱 + 分期）
// @Tags billing
// @Accept json
// @Produce json
// @Param input body ProcessDailyBillingInput false "扣款日期（選填，預設為今天）"
// @Success 200 {object} APIResponse{data=service.DailyBillingResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/billing/process-daily [post]
func (h *BillingHandler) ProcessDailyBilling(c *gin.Context) {
	var input ProcessDailyBillingInput

	// 綁定請求資料（選填）
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果沒有提供日期，使用今天
		input.Date = time.Now()
	}

	// 如果日期為零值，使用今天
	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	// 呼叫 service 處理每日扣款
	result, err := h.billingService.ProcessDailyBilling(input.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "PROCESS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: result,
	})
}

// ProcessSubscriptionBilling 處理訂閱扣款（手動觸發）
// @Summary 處理訂閱扣款
// @Description 手動觸發訂閱扣款處理
// @Tags billing
// @Accept json
// @Produce json
// @Param input body ProcessBillingInput false "扣款日期（選填，預設為今天）"
// @Success 200 {object} APIResponse{data=service.BillingResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/billing/process-subscriptions [post]
func (h *BillingHandler) ProcessSubscriptionBilling(c *gin.Context) {
	var input ProcessBillingInput

	// 綁定請求資料（選填）
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果沒有提供日期，使用今天
		input.Date = time.Now()
	}

	// 如果日期為零值，使用今天
	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	// 呼叫 service 處理訂閱扣款
	result, err := h.billingService.ProcessSubscriptionBilling(input.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "PROCESS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: result,
	})
}

// ProcessInstallmentBilling 處理分期扣款（手動觸發）
// @Summary 處理分期扣款
// @Description 手動觸發分期扣款處理
// @Tags billing
// @Accept json
// @Produce json
// @Param input body ProcessBillingInput false "扣款日期（選填，預設為今天）"
// @Success 200 {object} APIResponse{data=service.BillingResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/billing/process-installments [post]
func (h *BillingHandler) ProcessInstallmentBilling(c *gin.Context) {
	var input ProcessBillingInput

	// 綁定請求資料（選填）
	if err := c.ShouldBindJSON(&input); err != nil {
		// 如果沒有提供日期，使用今天
		input.Date = time.Now()
	}

	// 如果日期為零值，使用今天
	if input.Date.IsZero() {
		input.Date = time.Now()
	}

	// 呼叫 service 處理分期扣款
	result, err := h.billingService.ProcessInstallmentBilling(input.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "PROCESS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: result,
	})
}

// ProcessBillingInput 處理扣款的輸入
type ProcessBillingInput struct {
	Date time.Time `json:"date"`
}

// ProcessDailyBillingInput 處理每日扣款的輸入
type ProcessDailyBillingInput struct {
	Date time.Time `json:"date"`
}

