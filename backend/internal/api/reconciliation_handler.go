package api

import (
	"errors"
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// mapReconcileError 將 service 層 sentinel errors 映射成 HTTP status / API code。
func mapReconcileError(err error) (int, string) {
	switch {
	case errors.Is(err, service.ErrBankAccountNotFound),
		errors.Is(err, service.ErrCreditCardNotFound):
		return http.StatusNotFound, "NOT_FOUND"
	case errors.Is(err, service.ErrInvalidInput),
		errors.Is(err, service.ErrUsedExceedsLimit):
		return http.StatusBadRequest, "INVALID_INPUT"
	default:
		return http.StatusInternalServerError, "RECONCILE_FAILED"
	}
}

func writeReconcileError(c *gin.Context, err error) {
	status, code := mapReconcileError(err)
	c.JSON(status, APIResponse{Error: &APIError{Code: code, Message: err.Error()}})
}

// ReconciliationHandler 校準 API
type ReconciliationHandler struct {
	service service.ReconciliationService
}

// NewReconciliationHandler 建立新的 reconciliation handler
func NewReconciliationHandler(svc service.ReconciliationService) *ReconciliationHandler {
	return &ReconciliationHandler{service: svc}
}

// ReconcileBankAccount 校準銀行帳戶餘額
// @Summary 校準銀行帳戶餘額
// @Description 將銀行帳戶餘額校準為新值，並自動產生對帳 cash flow
// @Tags reconciliation
// @Accept json
// @Produce json
// @Param id path string true "銀行帳戶 ID"
// @Param input body models.ReconcileBankAccountInput true "校準資料"
// @Success 200 {object} APIResponse{data=models.ReconcileResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts/{id}/reconcile [post]
func (h *ReconciliationHandler) ReconcileBankAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{
			Code: "INVALID_ID", Message: "invalid bank account id",
		}})
		return
	}
	var input models.ReconcileBankAccountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{
			Code: "INVALID_INPUT", Message: err.Error(),
		}})
		return
	}
	res, err := h.service.ReconcileBankAccount(id, &input)
	if err != nil {
		writeReconcileError(c, err)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: res})
}

// ReconcileCreditCard 校準信用卡狀態
// @Summary 校準信用卡狀態
// @Description 將信用卡 used_credit 與/或 credit_limit 校準為新值
// @Tags reconciliation
// @Accept json
// @Produce json
// @Param id path string true "信用卡 ID"
// @Param input body models.ReconcileCreditCardInput true "校準資料"
// @Success 200 {object} APIResponse{data=models.ReconcileResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/{id}/reconcile [post]
func (h *ReconciliationHandler) ReconcileCreditCard(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{
			Code: "INVALID_ID", Message: "invalid credit card id",
		}})
		return
	}
	var input models.ReconcileCreditCardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{
			Code: "INVALID_INPUT", Message: err.Error(),
		}})
		return
	}
	res, err := h.service.ReconcileCreditCard(id, &input)
	if err != nil {
		writeReconcileError(c, err)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: res})
}

// ReconcileBatch 批次校準銀行帳戶與信用卡
// @Summary 批次校準
// @Description 一次校準多個銀行帳戶 / 信用卡，於單一 transaction 內執行
// @Tags reconciliation
// @Accept json
// @Produce json
// @Param input body models.ReconcileBatchInput true "批次校準資料"
// @Success 200 {object} APIResponse{data=[]models.ReconcileResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/user-management/reconcile-batch [post]
func (h *ReconciliationHandler) ReconcileBatch(c *gin.Context) {
	var input models.ReconcileBatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{
			Code: "INVALID_INPUT", Message: err.Error(),
		}})
		return
	}
	res, err := h.service.ReconcileBatch(&input)
	if err != nil {
		writeReconcileError(c, err)
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: res})
}
