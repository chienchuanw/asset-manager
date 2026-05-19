package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// FireHandler FIRE 投影 API Handler
type FireHandler struct {
	fireService service.FireService
}

// NewFireHandler 建立新的 FireHandler
func NewFireHandler(fireService service.FireService) *FireHandler {
	return &FireHandler{fireService: fireService}
}

// optionalFloatQuery 解析可選的浮點數查詢參數；缺少或無法解析時回傳 nil。
func optionalFloatQuery(c *gin.Context, key string) *float64 {
	raw := c.Query(key)
	if raw == "" {
		return nil
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return nil
	}
	return &v
}

// GetProjection 取得 FIRE 投影
// @Summary 取得 FIRE 投影
// @Description 依據自動推導的淨值與現金流（可由查詢參數覆寫）計算 FIRE 數字、進度與達標年數
// @Tags fire
// @Accept json
// @Produce json
// @Param netWorth query number false "目前淨值（TWD），覆寫最新資產快照"
// @Param annualExpenses query number false "年支出（TWD），覆寫近 12 個月推導值"
// @Param annualSavings query number false "年儲蓄（TWD），覆寫近 12 個月推導值"
// @Param expectedReturn query number false "預期年化報酬率（小數），預設 0.05"
// @Param withdrawalRate query number false "安全提領率（小數），預設 0.04"
// @Success 200 {object} models.FireProjectionResult
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/fire/projection [get]
func (h *FireHandler) GetProjection(c *gin.Context) {
	input := models.FireProjectionInput{
		NetWorth:       optionalFloatQuery(c, "netWorth"),
		AnnualExpenses: optionalFloatQuery(c, "annualExpenses"),
		AnnualSavings:  optionalFloatQuery(c, "annualSavings"),
		ExpectedReturn: optionalFloatQuery(c, "expectedReturn"),
		WithdrawalRate: optionalFloatQuery(c, "withdrawalRate"),
	}

	result, err := h.fireService.GetProjection(input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidWithdrawalRate) {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_WITHDRAWAL_RATE",
					Message: err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "FIRE_PROJECTION_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{Data: result})
}
