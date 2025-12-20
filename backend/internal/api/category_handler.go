package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CategoryHandler 現金流分類 API handler
type CategoryHandler struct {
	service service.CategoryService
}

// NewCategoryHandler 建立新的現金流分類 handler
func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

// CreateCategory 建立新的分類
// @Summary 建立分類
// @Description 建立新的現金流分類
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryInput true "分類資料"
// @Success 201 {object} APIResponse{data=models.CashFlowCategory}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input models.CreateCategoryInput

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

	// 呼叫 service 建立分類
	category, err := h.service.CreateCategory(&input)
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
		Data: category,
	})
}

// GetCategory 取得單筆分類
// @Summary 取得分類
// @Description 根據 ID 取得單筆分類
// @Tags categories
// @Produce json
// @Param id path string true "分類 ID"
// @Success 200 {object} APIResponse{data=models.CashFlowCategory}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid category ID format",
			},
		})
		return
	}

	// 呼叫 service 取得分類
	category, err := h.service.GetCategory(id)
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
		Data: category,
	})
}

// ListCategories 取得分類列表
// @Summary 取得分類列表
// @Description 取得所有分類，支援類型篩選
// @Tags categories
// @Produce json
// @Param type query string false "現金流類型 (income/expense)"
// @Success 200 {object} APIResponse{data=[]models.CashFlowCategory}
// @Failure 400 {object} APIResponse{error=APIError}
// @Router /api/categories [get]
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var flowType *models.CashFlowType

	// 類型篩選
	if typeStr := c.Query("type"); typeStr != "" {
		ft := models.CashFlowType(typeStr)
		flowType = &ft
	}

	// 呼叫 service 取得分類列表
	categories, err := h.service.ListCategories(flowType)
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
		Data: categories,
	})
}

// UpdateCategory 更新分類
// @Summary 更新分類
// @Description 更新分類（僅限自訂分類）
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "分類 ID"
// @Param category body models.UpdateCategoryInput true "更新資料"
// @Success 200 {object} APIResponse{data=models.CashFlowCategory}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid category ID format",
			},
		})
		return
	}

	var input models.UpdateCategoryInput

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

	// 呼叫 service 更新分類
	category, err := h.service.UpdateCategory(id, &input)
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
		Data: category,
	})
}

// DeleteCategory 刪除分類
// @Summary 刪除分類
// @Description 刪除分類（僅限自訂分類）
// @Tags categories
// @Produce json
// @Param id path string true "分類 ID"
// @Success 204 "No Content"
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid category ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除分類
	if err := h.service.DeleteCategory(id); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ReorderCategories 批次更新分類排序
// @Summary 重新排序分類
// @Description 批次更新分類的排序順序
// @Tags categories
// @Accept json
// @Produce json
// @Param orders body models.ReorderCategoryInput true "排序資料"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/categories/reorder [put]
func (h *CategoryHandler) ReorderCategories(c *gin.Context) {
	var input models.ReorderCategoryInput

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

	// 呼叫 service 重新排序
	if err := h.service.ReorderCategories(&input); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "REORDER_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: nil,
	})
}

