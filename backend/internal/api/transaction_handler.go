package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TransactionHandler 交易記錄 API handler
type TransactionHandler struct {
	service          service.TransactionService
	csvImportService service.CSVImportService
}

// NewTransactionHandler 建立新的交易記錄 handler
func NewTransactionHandler(service service.TransactionService, csvImportService service.CSVImportService) *TransactionHandler {
	return &TransactionHandler{
		service:          service,
		csvImportService: csvImportService,
	}
}

// APIResponse 統一的 API 回應格式
type APIResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIError API 錯誤格式
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// CreateTransaction 建立新的交易記錄
// @Summary 建立交易記錄
// @Description 建立新的交易記錄
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.CreateTransactionInput true "交易記錄資料"
// @Success 201 {object} APIResponse{data=models.Transaction}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var input models.CreateTransactionInput

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

	// 呼叫 service 建立交易記錄
	transaction, err := h.service.CreateTransaction(&input)
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
		Data: transaction,
	})
}

// CreateTransactionsBatch 批次建立交易記錄
// @Summary 批次建立交易記錄
// @Description 批次建立多筆交易記錄（全有或全無）
// @Tags transactions
// @Accept json
// @Produce json
// @Param batch body models.BatchCreateTransactionsInput true "批次交易記錄資料"
// @Success 201 {object} APIResponse{data=[]models.Transaction}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/transactions/batch [post]
func (h *TransactionHandler) CreateTransactionsBatch(c *gin.Context) {
	var input models.BatchCreateTransactionsInput

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

	// 呼叫 service 批次建立交易記錄
	transactions, err := h.service.CreateTransactionsBatch(input.Transactions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "BATCH_CREATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Data: transactions,
	})
}

// GetTransaction 取得單筆交易記錄
// @Summary 取得交易記錄
// @Description 根據 ID 取得單筆交易記錄
// @Tags transactions
// @Produce json
// @Param id path string true "交易記錄 ID"
// @Success 200 {object} APIResponse{data=models.Transaction}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid transaction ID format",
			},
		})
		return
	}

	// 呼叫 service 取得交易記錄
	transaction, err := h.service.GetTransaction(id)
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
		Data: transaction,
	})
}

// ListTransactions 取得交易記錄列表
// @Summary 取得交易記錄列表
// @Description 取得所有交易記錄，支援篩選
// @Tags transactions
// @Produce json
// @Param asset_type query string false "資產類型"
// @Param transaction_type query string false "交易類型"
// @Param symbol query string false "代碼"
// @Param start_date query string false "開始日期 (YYYY-MM-DD)"
// @Param end_date query string false "結束日期 (YYYY-MM-DD)"
// @Param limit query int false "每頁筆數"
// @Param offset query int false "偏移量"
// @Success 200 {object} APIResponse{data=[]models.Transaction}
// @Failure 400 {object} APIResponse{error=APIError}
// @Router /api/transactions [get]
func (h *TransactionHandler) ListTransactions(c *gin.Context) {
	// 解析查詢參數
	filters := repository.TransactionFilters{}

	// 資產類型篩選
	if assetTypeStr := c.Query("asset_type"); assetTypeStr != "" {
		assetType := models.AssetType(assetTypeStr)
		filters.AssetType = &assetType
	}

	// 交易類型篩選
	if transactionTypeStr := c.Query("transaction_type"); transactionTypeStr != "" {
		transactionType := models.TransactionType(transactionTypeStr)
		filters.TransactionType = &transactionType
	}

	// 代碼篩選
	if symbol := c.Query("symbol"); symbol != "" {
		filters.Symbol = &symbol
	}

	// 日期範圍篩選
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, APIResponse{
				Error: &APIError{
					Code:    "INVALID_DATE",
					Message: "Invalid start_date format, expected YYYY-MM-DD",
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
					Code:    "INVALID_DATE",
					Message: "Invalid end_date format, expected YYYY-MM-DD",
				},
			})
			return
		}
		filters.EndDate = &endDate
	}

	// 分頁參數
	if limit := c.Query("limit"); limit != "" {
		var limitInt int
		if _, err := fmt.Sscanf(limit, "%d", &limitInt); err == nil {
			filters.Limit = limitInt
		}
	}

	if offset := c.Query("offset"); offset != "" {
		var offsetInt int
		if _, err := fmt.Sscanf(offset, "%d", &offsetInt); err == nil {
			filters.Offset = offsetInt
		}
	}

	// 呼叫 service 取得交易記錄列表
	transactions, err := h.service.ListTransactions(filters)
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
		Data: transactions,
	})
}

// UpdateTransaction 更新交易記錄
// @Summary 更新交易記錄
// @Description 更新指定的交易記錄
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "交易記錄 ID"
// @Param transaction body models.UpdateTransactionInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.Transaction}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid transaction ID format",
			},
		})
		return
	}

	// 綁定並驗證請求資料
	var input models.UpdateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	// 呼叫 service 更新交易記錄
	transaction, err := h.service.UpdateTransaction(id, &input)
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
		Data: transaction,
	})
}

// DeleteTransaction 刪除交易記錄
// @Summary 刪除交易記錄
// @Description 刪除指定的交易記錄
// @Tags transactions
// @Produce json
// @Param id path string true "交易記錄 ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid transaction ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除交易記錄
	err = h.service.DeleteTransaction(id)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{"message": "Transaction deleted successfully"},
	})
}

// DownloadCSVTemplate 下載 CSV 樣板
// @Summary 下載 CSV 樣板
// @Description 下載交易記錄 CSV 匯入樣板檔案
// @Tags transactions
// @Produce text/csv
// @Success 200 {file} string "CSV 樣板檔案"
// @Router /api/transactions/template [get]
func (h *TransactionHandler) DownloadCSVTemplate(c *gin.Context) {
	// 生成 CSV 樣板
	csvContent := h.csvImportService.GenerateTemplate()

	// 設定回應 header
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=transaction_template.csv")

	// 返回 CSV 內容
	c.String(http.StatusOK, csvContent)
}

// ParseCSV 解析 CSV 檔案
// @Summary 解析 CSV 檔案
// @Description 解析上傳的 CSV 檔案並返回解析後的交易資料或錯誤訊息
// @Tags transactions
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "CSV 檔案"
// @Success 200 {object} APIResponse{data=models.CSVImportResult}
// @Failure 400 {object} APIResponse{error=APIError}
// @Router /api/transactions/parse-csv [post]
func (h *TransactionHandler) ParseCSV(c *gin.Context) {
	// 取得上傳的檔案
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_FILE",
				Message: "無法讀取上傳的檔案",
			},
		})
		return
	}

	// 開啟檔案
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "FILE_OPEN_ERROR",
				Message: "無法開啟檔案",
			},
		})
		return
	}
	defer f.Close()

	// 解析 CSV
	result := h.csvImportService.ParseCSV(f)

	// 返回結果
	c.JSON(http.StatusOK, APIResponse{
		Data: result,
	})
}
