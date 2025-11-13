package api

import (
	"log"
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// HoldingHandler 持倉 API Handler
type HoldingHandler struct {
	holdingService service.HoldingService
}

// NewHoldingHandler 建立新的 Holding Handler
func NewHoldingHandler(holdingService service.HoldingService) *HoldingHandler {
	return &HoldingHandler{
		holdingService: holdingService,
	}
}

// GetAllHoldings 取得所有持倉
// @Summary 取得所有持倉
// @Description 取得所有持倉列表，支援按資產類型和標的代碼篩選
// @Tags holdings
// @Accept json
// @Produce json
// @Param asset_type query string false "資產類型 (cash, tw-stock, us-stock, crypto)"
// @Param symbol query string false "標的代碼"
// @Success 200 {object} map[string]interface{} "成功返回持倉列表"
// @Failure 500 {object} map[string]interface{} "伺服器錯誤"
// @Router /api/holdings [get]
func (h *HoldingHandler) GetAllHoldings(c *gin.Context) {
	log.Println("=== [DEBUG] GetAllHoldings API called ===")

	// 解析查詢參數
	var filters models.HoldingFilters

	// 資產類型篩選
	if assetTypeStr := c.Query("asset_type"); assetTypeStr != "" {
		assetType := models.AssetType(assetTypeStr)
		filters.AssetType = &assetType
		log.Printf("[DEBUG] Filter by asset_type: %s", assetTypeStr)
	}

	// 標的代碼篩選
	if symbol := c.Query("symbol"); symbol != "" {
		filters.Symbol = &symbol
		log.Printf("[DEBUG] Filter by symbol: %s", symbol)
	}

	log.Println("[DEBUG] Calling holdingService.GetAllHoldings...")

	// 呼叫 Service 層
	holdings, err := h.holdingService.GetAllHoldings(filters)
	if err != nil {
		log.Printf("[ERROR] GetAllHoldings failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	log.Printf("[DEBUG] GetAllHoldings success, returned %d holdings", len(holdings))

	// 返回成功結果
	c.JSON(http.StatusOK, gin.H{
		"data":  holdings,
		"error": nil,
	})
}

// GetHoldingBySymbol 取得單一標的持倉
// @Summary 取得單一標的持倉
// @Description 根據標的代碼取得持倉詳情
// @Tags holdings
// @Accept json
// @Produce json
// @Param symbol path string true "標的代碼"
// @Success 200 {object} map[string]interface{} "成功返回持倉詳情"
// @Failure 400 {object} map[string]interface{} "請求參數錯誤"
// @Failure 500 {object} map[string]interface{} "伺服器錯誤"
// @Router /api/holdings/{symbol} [get]
func (h *HoldingHandler) GetHoldingBySymbol(c *gin.Context) {
	// 取得路徑參數
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"data":  nil,
			"error": gin.H{
				"code":    "INVALID_PARAMETER",
				"message": "symbol is required",
			},
		})
		return
	}

	// 呼叫 Service 層
	holding, err := h.holdingService.GetHoldingBySymbol(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"data":  nil,
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// 返回成功結果
	c.JSON(http.StatusOK, gin.H{
		"data":  holding,
		"error": nil,
	})
}

