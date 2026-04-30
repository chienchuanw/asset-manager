package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// HoldingReconciliationHandler 持倉對帳 API
type HoldingReconciliationHandler struct {
	service service.HoldingReconciliationService
}

// NewHoldingReconciliationHandler 建立新的 holding reconciliation handler
func NewHoldingReconciliationHandler(svc service.HoldingReconciliationService) *HoldingReconciliationHandler {
	return &HoldingReconciliationHandler{service: svc}
}

// Preview godoc
// @Summary  預覽持倉對帳變動
// @Description 對帳前 dry-run，回傳每項 (asset_type, symbol) 的 prev/target/delta/action 預覽
// @Tags     holdings
// @Accept   json
// @Produce  json
// @Param    body body models.HoldingReconcileBatchInput true "對帳項目"
// @Success  200 {object} APIResponse{data=models.HoldingReconcilePreview}
// @Failure  400 {object} APIResponse
// @Failure  500 {object} APIResponse
// @Router   /holdings/reconcile/preview [post]
func (h *HoldingReconciliationHandler) Preview(c *gin.Context) {
	h.run(c, true)
}

// Execute godoc
// @Summary  執行持倉對帳
// @Description 為每個非 noop 項目寫入一筆 adjustment 交易；不結算實現損益、不影響現金流
// @Tags     holdings
// @Accept   json
// @Produce  json
// @Param    body body models.HoldingReconcileBatchInput true "對帳項目"
// @Success  200 {object} APIResponse{data=models.HoldingReconcilePreview}
// @Failure  400 {object} APIResponse
// @Failure  500 {object} APIResponse
// @Router   /holdings/reconcile [post]
func (h *HoldingReconciliationHandler) Execute(c *gin.Context) {
	h.run(c, false)
}

// ExecuteSingle godoc
// @Summary  執行單一 symbol 的持倉對帳
// @Description path 上的 symbol 會覆寫 body 內的 symbol（若有不一致）
// @Tags     holdings
// @Accept   json
// @Produce  json
// @Param    symbol path string true "標的代碼"
// @Param    body body models.HoldingReconcileItem true "對帳項目"
// @Success  200 {object} APIResponse{data=models.HoldingReconcilePreview}
// @Failure  400 {object} APIResponse
// @Failure  500 {object} APIResponse
// @Router   /holdings/{symbol}/reconcile [post]
func (h *HoldingReconciliationHandler) ExecuteSingle(c *gin.Context) {
	var item models.HoldingReconcileItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{Code: "INVALID_INPUT", Message: err.Error()}})
		return
	}
	if symbol := c.Param("symbol"); symbol != "" {
		item.Symbol = symbol
	}
	preview, err := h.service.ReconcileHoldings(c.Request.Context(), []models.HoldingReconcileItem{item}, false)
	h.respond(c, preview, err)
}

func (h *HoldingReconciliationHandler) run(c *gin.Context, dryRun bool) {
	var body models.HoldingReconcileBatchInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{Code: "INVALID_INPUT", Message: err.Error()}})
		return
	}
	preview, err := h.service.ReconcileHoldings(c.Request.Context(), body.Items, dryRun)
	h.respond(c, preview, err)
}

func (h *HoldingReconciliationHandler) respond(c *gin.Context, preview *models.HoldingReconcilePreview, err error) {
	if err != nil {
		if service.IsReconcileValidationError(err) {
			c.JSON(http.StatusBadRequest, APIResponse{Error: &APIError{Code: "INVALID_INPUT", Message: err.Error()}})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{Error: &APIError{Code: "RECONCILE_FAILED", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, APIResponse{Data: preview})
}
