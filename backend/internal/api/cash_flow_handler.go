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

// CashFlowHandler 現金流記錄 API handler
type CashFlowHandler struct {
	service service.CashFlowService
}

// NewCashFlowHandler 建立新的現金流記錄 handler
func NewCashFlowHandler(service service.CashFlowService) *CashFlowHandler {
	return &CashFlowHandler{service: service}
}

// CreateCashFlow 建立新的現金流記錄
// @Summary 建立現金流記錄
// @Description 建立新的現金流記錄
// @Tags cash-flows
// @Accept json
// @Produce json
// @Param cash_flow body models.CreateCashFlowInput true "現金流記錄資料"
// @Success 201 {object} APIResponse{data=models.CashFlow}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/cash-flows [post]
func (h *CashFlowHandler) CreateCashFlow(c *gin.Context) {
	var input models.CreateCashFlowInput

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

	// 呼叫 service 建立現金流記錄
	cashFlow, err := h.service.CreateCashFlow(&input)
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
		Data: cashFlow,
	})
}

// GetCashFlow 取得單筆現金流記錄
// @Summary 取得現金流記錄
// @Description 根據 ID 取得單筆現金流記錄
// @Tags cash-flows
// @Produce json
// @Param id path string true "現金流記錄 ID"
// @Success 200 {object} APIResponse{data=models.CashFlow}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/cash-flows/{id} [get]
func (h *CashFlowHandler) GetCashFlow(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid cash flow ID format",
			},
		})
		return
	}

	// 呼叫 service 取得現金流記錄
	cashFlow, err := h.service.GetCashFlow(id)
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
		Data: cashFlow,
	})
}

// ListCashFlows 取得現金流記錄列表
// @Summary 取得現金流記錄列表
// @Description 取得所有現金流記錄，支援篩選
// @Tags cash-flows
// @Produce json
// @Param type query string false "現金流類型 (income/expense)"
// @Param category_id query string false "分類 ID"
// @Param start_date query string false "開始日期 (YYYY-MM-DD)"
// @Param end_date query string false "結束日期 (YYYY-MM-DD)"
// @Param limit query int false "每頁筆數"
// @Param offset query int false "偏移量"
// @Success 200 {object} APIResponse{data=[]models.CashFlow}
// @Failure 400 {object} APIResponse{error=APIError}
// @Router /api/cash-flows [get]
func (h *CashFlowHandler) ListCashFlows(c *gin.Context) {
	// 解析查詢參數
	filters := repository.CashFlowFilters{}

	// 類型篩選
	if typeStr := c.Query("type"); typeStr != "" {
		flowType := models.CashFlowType(typeStr)
		filters.Type = &flowType
	}

	// 分類篩選
	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_CATEGORY_ID",
					Message: "Invalid category ID format",
				},
			})
			return
		}
		filters.CategoryID = &categoryID
	}

	// 日期範圍篩選
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_START_DATE",
					Message: "Invalid start date format, use YYYY-MM-DD",
				},
			})
			return
		}
		filters.StartDate = &startDate
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_END_DATE",
					Message: "Invalid end date format, use YYYY-MM-DD",
				},
			})
			return
		}
		filters.EndDate = &endDate
	}

	// 分頁參數
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_LIMIT",
					Message: "Invalid limit parameter",
				},
			})
			return
		}
		filters.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_OFFSET",
					Message: "Invalid offset parameter",
				},
			})
			return
		}
		filters.Offset = offset
	}

	// 呼叫 service 取得現金流記錄列表
	cashFlows, err := h.service.ListCashFlows(filters)
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
		Data: cashFlows,
	})
}

// UpdateCashFlow 更新現金流記錄
// @Summary 更新現金流記錄
// @Description 更新現金流記錄
// @Tags cash-flows
// @Accept json
// @Produce json
// @Param id path string true "現金流記錄 ID"
// @Param cash_flow body models.UpdateCashFlowInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.CashFlow}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/cash-flows/{id} [put]
func (h *CashFlowHandler) UpdateCashFlow(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid cash flow ID format",
			},
		})
		return
	}

	var input models.UpdateCashFlowInput

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

	// 呼叫 service 更新現金流記錄
	cashFlow, err := h.service.UpdateCashFlow(id, &input)
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
		Data: cashFlow,
	})
}

// DeleteCashFlow 刪除現金流記錄
// @Summary 刪除現金流記錄
// @Description 刪除現金流記錄
// @Tags cash-flows
// @Produce json
// @Param id path string true "現金流記錄 ID"
// @Success 204 "No Content"
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/cash-flows/{id} [delete]
func (h *CashFlowHandler) DeleteCashFlow(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid cash flow ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除現金流記錄
	if err := h.service.DeleteCashFlow(id); err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetSummary 取得現金流摘要
// @Summary 取得現金流摘要
// @Description 取得指定日期範圍的現金流摘要統計
// @Tags cash-flows
// @Produce json
// @Param start_date query string true "開始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} APIResponse{data=repository.CashFlowSummary}
// @Failure 400 {object} APIResponse{error=APIError}
// @Router /api/cash-flows/summary [get]
func (h *CashFlowHandler) GetSummary(c *gin.Context) {
	// 解析日期參數
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "MISSING_PARAMETERS",
				Message: "start_date and end_date are required",
			},
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_START_DATE",
				Message: "Invalid start date format, use YYYY-MM-DD",
			},
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_END_DATE",
				Message: "Invalid end date format, use YYYY-MM-DD",
			},
		})
		return
	}

	// 呼叫 service 取得摘要
	summary, err := h.service.GetSummary(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "SUMMARY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: summary,
	})
}

