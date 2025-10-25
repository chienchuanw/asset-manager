package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// RebalanceHandler 再平衡 Handler
type RebalanceHandler struct {
	service service.RebalanceService
}

// NewRebalanceHandler 建立再平衡 Handler
func NewRebalanceHandler(service service.RebalanceService) *RebalanceHandler {
	return &RebalanceHandler{
		service: service,
	}
}

// CheckRebalance 檢查是否需要再平衡
// @Summary 檢查是否需要再平衡
// @Description 檢查當前資產配置是否偏離目標配置，並提供再平衡建議
// @Tags rebalance
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=models.RebalanceCheck}
// @Failure 500 {object} APIResponse
// @Router /api/rebalance/check [get]
func (h *RebalanceHandler) CheckRebalance(c *gin.Context) {
	result, err := h.service.CheckRebalance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "CHECK_REBALANCE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: result,
	})
}

