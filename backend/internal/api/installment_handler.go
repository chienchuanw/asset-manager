package api

import (
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// InstallmentHandler 分期 API handler
type InstallmentHandler struct {
	service service.InstallmentService
}

// NewInstallmentHandler 建立新的分期 handler
func NewInstallmentHandler(service service.InstallmentService) *InstallmentHandler {
	return &InstallmentHandler{service: service}
}

// CreateInstallment 建立新的分期
// @Summary 建立分期
// @Description 建立新的分期
// @Tags installments
// @Accept json
// @Produce json
// @Param installment body models.CreateInstallmentInput true "分期資料"
// @Success 201 {object} APIResponse{data=models.Installment}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/installments [post]
func (h *InstallmentHandler) CreateInstallment(c *gin.Context) {
	var input models.CreateInstallmentInput

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

	// 呼叫 service 建立分期
	installment, err := h.service.CreateInstallment(&input)
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
		Data: installment,
	})
}

// GetInstallment 取得單筆分期
// @Summary 取得分期
// @Description 根據 ID 取得單筆分期
// @Tags installments
// @Produce json
// @Param id path string true "分期 ID"
// @Success 200 {object} APIResponse{data=models.Installment}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/installments/{id} [get]
func (h *InstallmentHandler) GetInstallment(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid installment ID format",
			},
		})
		return
	}

	// 呼叫 service 取得分期
	installment, err := h.service.GetInstallment(id)
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
		Data: installment,
	})
}

// ListInstallments 取得分期列表
// @Summary 取得分期列表
// @Description 取得分期列表，支援篩選和分頁
// @Tags installments
// @Produce json
// @Param status query string false "狀態篩選 (active, completed, cancelled)"
// @Param limit query int false "每頁筆數" default(100)
// @Param offset query int false "略過筆數" default(0)
// @Success 200 {object} APIResponse{data=[]models.Installment}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/installments [get]
func (h *InstallmentHandler) ListInstallments(c *gin.Context) {
	// 解析查詢參數
	filters := repository.InstallmentFilters{}

	// 狀態篩選
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.InstallmentStatus(statusStr)
		filters.Status = &status
	}

	// 分頁參數
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filters.Offset = offset
		}
	}

	// 呼叫 service 取得分期列表
	installments, err := h.service.ListInstallments(filters)
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
		Data: installments,
	})
}

// UpdateInstallment 更新分期
// @Summary 更新分期
// @Description 更新分期資料
// @Tags installments
// @Accept json
// @Produce json
// @Param id path string true "分期 ID"
// @Param installment body models.UpdateInstallmentInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.Installment}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/installments/{id} [put]
func (h *InstallmentHandler) UpdateInstallment(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid installment ID format",
			},
		})
		return
	}

	var input models.UpdateInstallmentInput

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

	// 呼叫 service 更新分期
	installment, err := h.service.UpdateInstallment(id, &input)
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
		Data: installment,
	})
}

// DeleteInstallment 刪除分期
// @Summary 刪除分期
// @Description 刪除分期
// @Tags installments
// @Produce json
// @Param id path string true "分期 ID"
// @Success 200 {object} APIResponse{data=string}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/installments/{id} [delete]
func (h *InstallmentHandler) DeleteInstallment(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid installment ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除分期
	if err := h.service.DeleteInstallment(id); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Installment deleted successfully",
	})
}

// GetCompletingSoon 取得即將完成的分期
// @Summary 取得即將完成的分期
// @Description 取得剩餘期數小於等於指定值的分期
// @Tags installments
// @Produce json
// @Param remaining_count query int false "剩餘期數" default(3)
// @Success 200 {object} APIResponse{data=[]models.Installment}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/installments/completing-soon [get]
func (h *InstallmentHandler) GetCompletingSoon(c *gin.Context) {
	// 解析剩餘期數參數
	remainingCount := 3 // 預設值
	if countStr := c.Query("remaining_count"); countStr != "" {
		if count, err := strconv.Atoi(countStr); err == nil {
			remainingCount = count
		}
	}

	// 呼叫 service 取得即將完成的分期
	installments, err := h.service.GetCompletingSoon(remainingCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: installments,
	})
}

