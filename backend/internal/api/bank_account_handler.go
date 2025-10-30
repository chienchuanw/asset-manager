package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// BankAccountHandler 銀行帳戶 API handler
type BankAccountHandler struct {
	service service.BankAccountService
}

// NewBankAccountHandler 建立新的銀行帳戶 handler
func NewBankAccountHandler(service service.BankAccountService) *BankAccountHandler {
	return &BankAccountHandler{service: service}
}

// CreateBankAccount 建立新的銀行帳戶
// @Summary 建立銀行帳戶
// @Description 建立新的銀行帳戶
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param account body models.CreateBankAccountInput true "銀行帳戶資料"
// @Success 201 {object} APIResponse{data=models.BankAccount}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts [post]
func (h *BankAccountHandler) CreateBankAccount(c *gin.Context) {
	var input models.CreateBankAccountInput

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

	// 呼叫 service 建立銀行帳戶
	account, err := h.service.CreateBankAccount(&input)
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
		Data: account,
	})
}

// GetBankAccount 取得單筆銀行帳戶
// @Summary 取得銀行帳戶
// @Description 根據 ID 取得單筆銀行帳戶
// @Tags bank-accounts
// @Produce json
// @Param id path string true "銀行帳戶 ID"
// @Success 200 {object} APIResponse{data=models.BankAccount}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts/{id} [get]
func (h *BankAccountHandler) GetBankAccount(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid bank account ID format",
			},
		})
		return
	}

	// 呼叫 service 取得銀行帳戶
	account, err := h.service.GetBankAccount(id)
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
		Data: account,
	})
}

// ListBankAccounts 列出所有銀行帳戶
// @Summary 列出銀行帳戶
// @Description 列出所有銀行帳戶，可選擇性依幣別篩選
// @Tags bank-accounts
// @Produce json
// @Param currency query string false "幣別篩選 (TWD, USD)"
// @Success 200 {object} APIResponse{data=[]models.BankAccount}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts [get]
func (h *BankAccountHandler) ListBankAccounts(c *gin.Context) {
	// 取得查詢參數
	currencyStr := c.Query("currency")
	var currency *models.Currency
	if currencyStr != "" {
		curr := models.Currency(currencyStr)
		currency = &curr
	}

	// 呼叫 service 列出銀行帳戶
	accounts, err := h.service.ListBankAccounts(currency)
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
		Data: accounts,
	})
}

// UpdateBankAccount 更新銀行帳戶
// @Summary 更新銀行帳戶
// @Description 更新銀行帳戶資料
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "銀行帳戶 ID"
// @Param account body models.UpdateBankAccountInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.BankAccount}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts/{id} [put]
func (h *BankAccountHandler) UpdateBankAccount(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid bank account ID format",
			},
		})
		return
	}

	var input models.UpdateBankAccountInput

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

	// 呼叫 service 更新銀行帳戶
	account, err := h.service.UpdateBankAccount(id, &input)
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
		Data: account,
	})
}

// DeleteBankAccount 刪除銀行帳戶
// @Summary 刪除銀行帳戶
// @Description 刪除銀行帳戶
// @Tags bank-accounts
// @Produce json
// @Param id path string true "銀行帳戶 ID"
// @Success 200 {object} APIResponse{data=map[string]string}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/bank-accounts/{id} [delete]
func (h *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid bank account ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除銀行帳戶
	err = h.service.DeleteBankAccount(id)
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
			"message": "Bank account deleted successfully",
		},
	})
}

