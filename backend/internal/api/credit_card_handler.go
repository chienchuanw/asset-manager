package api

import (
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreditCardHandler 信用卡 API handler
type CreditCardHandler struct {
	service service.CreditCardService
}

// NewCreditCardHandler 建立新的信用卡 handler
func NewCreditCardHandler(service service.CreditCardService) *CreditCardHandler {
	return &CreditCardHandler{service: service}
}

// CreateCreditCard 建立新的信用卡
// @Summary 建立信用卡
// @Description 建立新的信用卡
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param card body models.CreateCreditCardInput true "信用卡資料"
// @Success 201 {object} APIResponse{data=models.CreditCard}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards [post]
func (h *CreditCardHandler) CreateCreditCard(c *gin.Context) {
	var input models.CreateCreditCardInput

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

	// 呼叫 service 建立信用卡
	card, err := h.service.CreateCreditCard(&input)
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
		Data: card,
	})
}

// GetCreditCard 取得單筆信用卡
// @Summary 取得信用卡
// @Description 根據 ID 取得單筆信用卡
// @Tags credit-cards
// @Produce json
// @Param id path string true "信用卡 ID"
// @Success 200 {object} APIResponse{data=models.CreditCard}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/{id} [get]
func (h *CreditCardHandler) GetCreditCard(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid credit card ID format",
			},
		})
		return
	}

	// 呼叫 service 取得信用卡
	card, err := h.service.GetCreditCard(id)
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
		Data: card,
	})
}

// ListCreditCards 列出所有信用卡
// @Summary 列出信用卡
// @Description 列出所有信用卡
// @Tags credit-cards
// @Produce json
// @Success 200 {object} APIResponse{data=[]models.CreditCard}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards [get]
func (h *CreditCardHandler) ListCreditCards(c *gin.Context) {
	// 呼叫 service 列出信用卡
	cards, err := h.service.ListCreditCards()
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
		Data: cards,
	})
}

// GetUpcomingBilling 取得即將到來的帳單日信用卡
// @Summary 取得即將到來的帳單日信用卡
// @Description 取得未來 N 天內的帳單日信用卡
// @Tags credit-cards
// @Produce json
// @Param days_ahead query int false "未來天數 (預設: 7)"
// @Success 200 {object} APIResponse{data=[]models.CreditCard}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/upcoming-billing [get]
func (h *CreditCardHandler) GetUpcomingBilling(c *gin.Context) {
	// 取得查詢參數
	daysAheadStr := c.DefaultQuery("days_ahead", "7")
	daysAhead, err := strconv.Atoi(daysAheadStr)
	if err != nil || daysAhead < 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_PARAMETER",
				Message: "days_ahead must be a non-negative integer",
			},
		})
		return
	}

	// 呼叫 service 取得即將到來的帳單日信用卡
	cards, err := h.service.GetUpcomingBilling(daysAhead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "QUERY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: cards,
	})
}

// GetUpcomingPayment 取得即將到來的繳款截止日信用卡
// @Summary 取得即將到來的繳款截止日信用卡
// @Description 取得未來 N 天內的繳款截止日信用卡
// @Tags credit-cards
// @Produce json
// @Param days_ahead query int false "未來天數 (預設: 7)"
// @Success 200 {object} APIResponse{data=[]models.CreditCard}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/upcoming-payment [get]
func (h *CreditCardHandler) GetUpcomingPayment(c *gin.Context) {
	// 取得查詢參數
	daysAheadStr := c.DefaultQuery("days_ahead", "7")
	daysAhead, err := strconv.Atoi(daysAheadStr)
	if err != nil || daysAhead < 0 {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_PARAMETER",
				Message: "days_ahead must be a non-negative integer",
			},
		})
		return
	}

	// 呼叫 service 取得即將到來的繳款截止日信用卡
	cards, err := h.service.GetUpcomingPayment(daysAhead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "QUERY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: cards,
	})
}

// UpdateCreditCard 更新信用卡
// @Summary 更新信用卡
// @Description 更新信用卡資料
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param id path string true "信用卡 ID"
// @Param card body models.UpdateCreditCardInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.CreditCard}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/{id} [put]
func (h *CreditCardHandler) UpdateCreditCard(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid credit card ID format",
			},
		})
		return
	}

	var input models.UpdateCreditCardInput

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

	// 呼叫 service 更新信用卡
	card, err := h.service.UpdateCreditCard(id, &input)
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
		Data: card,
	})
}

// DeleteCreditCard 刪除信用卡
// @Summary 刪除信用卡
// @Description 刪除信用卡
// @Tags credit-cards
// @Produce json
// @Param id path string true "信用卡 ID"
// @Success 200 {object} APIResponse{data=map[string]string}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-cards/{id} [delete]
func (h *CreditCardHandler) DeleteCreditCard(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid credit card ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除信用卡
	err = h.service.DeleteCreditCard(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: map[string]string{
			"message": "Credit card deleted successfully",
		},
	})
}

