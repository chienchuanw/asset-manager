package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreditCardGroupHandler 信用卡群組 API handler
type CreditCardGroupHandler struct {
	service service.CreditCardGroupService
}

// NewCreditCardGroupHandler 建立新的信用卡群組 handler
func NewCreditCardGroupHandler(service service.CreditCardGroupService) *CreditCardGroupHandler {
	return &CreditCardGroupHandler{service: service}
}

// CreateCreditCardGroup 建立新的信用卡群組
// @Summary 建立信用卡群組
// @Description 建立新的信用卡群組,將多張信用卡組成共享額度群組
// @Tags credit-card-groups
// @Accept json
// @Produce json
// @Param group body models.CreateCreditCardGroupInput true "信用卡群組資料"
// @Success 201 {object} APIResponse{data=models.CreditCardGroupWithCards}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups [post]
func (h *CreditCardGroupHandler) CreateCreditCardGroup(c *gin.Context) {
	var input models.CreateCreditCardGroupInput

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

	// 呼叫 service 建立信用卡群組
	group, err := h.service.CreateCreditCardGroup(&input)
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
		Data: group,
	})
}

// GetCreditCardGroup 取得單筆信用卡群組
// @Summary 取得信用卡群組
// @Description 根據 ID 取得單筆信用卡群組及其包含的卡片
// @Tags credit-card-groups
// @Produce json
// @Param id path string true "信用卡群組 ID"
// @Success 200 {object} APIResponse{data=models.CreditCardGroupWithCards}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups/{id} [get]
func (h *CreditCardGroupHandler) GetCreditCardGroup(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid group ID format",
			},
		})
		return
	}

	// 呼叫 service 取得信用卡群組
	group, err := h.service.GetCreditCardGroup(id)
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
		Data: group,
	})
}

// ListCreditCardGroups 取得所有信用卡群組
// @Summary 取得所有信用卡群組
// @Description 取得所有信用卡群組列表
// @Tags credit-card-groups
// @Produce json
// @Success 200 {object} APIResponse{data=[]models.CreditCardGroupWithCards}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups [get]
func (h *CreditCardGroupHandler) ListCreditCardGroups(c *gin.Context) {
	// 呼叫 service 取得所有信用卡群組
	groups, err := h.service.ListCreditCardGroups()
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
		Data: groups,
	})
}

// UpdateCreditCardGroup 更新信用卡群組
// @Summary 更新信用卡群組
// @Description 更新信用卡群組資料
// @Tags credit-card-groups
// @Accept json
// @Produce json
// @Param id path string true "信用卡群組 ID"
// @Param group body models.UpdateCreditCardGroupInput true "更新的信用卡群組資料"
// @Success 200 {object} APIResponse{data=models.CreditCardGroup}
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups/{id} [put]
func (h *CreditCardGroupHandler) UpdateCreditCardGroup(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid group ID format",
			},
		})
		return
	}

	var input models.UpdateCreditCardGroupInput

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

	// 呼叫 service 更新信用卡群組
	group, err := h.service.UpdateCreditCardGroup(id, &input)
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
		Data: group,
	})
}

// DeleteCreditCardGroup 刪除信用卡群組
// @Summary 刪除信用卡群組
// @Description 刪除信用卡群組,群組內的卡片將恢復為獨立卡片
// @Tags credit-card-groups
// @Produce json
// @Param id path string true "信用卡群組 ID"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups/{id} [delete]
func (h *CreditCardGroupHandler) DeleteCreditCardGroup(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid group ID format",
			},
		})
		return
	}

	// 呼叫 service 刪除信用卡群組
	if err := h.service.DeleteCreditCardGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "DELETE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{"message": "Credit card group deleted successfully"},
	})
}

// AddCardsToGroup 新增卡片到群組
// @Summary 新增卡片到群組
// @Description 將一張或多張信用卡加入到現有群組
// @Tags credit-card-groups
// @Accept json
// @Produce json
// @Param id path string true "信用卡群組 ID"
// @Param cards body models.AddCardsToGroupInput true "要加入的卡片 ID 列表"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups/{id}/cards [post]
func (h *CreditCardGroupHandler) AddCardsToGroup(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid group ID format",
			},
		})
		return
	}

	var input models.AddCardsToGroupInput

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

	// 呼叫 service 新增卡片到群組
	if err := h.service.AddCardsToGroup(id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "ADD_CARDS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{"message": "Cards added to group successfully"},
	})
}

// RemoveCardsFromGroup 從群組移除卡片
// @Summary 從群組移除卡片
// @Description 將一張或多張信用卡從群組中移除,卡片將恢復為獨立卡片
// @Tags credit-card-groups
// @Accept json
// @Produce json
// @Param id path string true "信用卡群組 ID"
// @Param cards body models.RemoveCardsFromGroupInput true "要移除的卡片 ID 列表"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse{error=APIError}
// @Failure 404 {object} APIResponse{error=APIError}
// @Failure 500 {object} APIResponse{error=APIError}
// @Router /api/credit-card-groups/{id}/cards [delete]
func (h *CreditCardGroupHandler) RemoveCardsFromGroup(c *gin.Context) {
	// 解析 ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_ID",
				Message: "Invalid group ID format",
			},
		})
		return
	}

	var input models.RemoveCardsFromGroupInput

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

	// 呼叫 service 從群組移除卡片
	if err := h.service.RemoveCardsFromGroup(id, &input); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "REMOVE_CARDS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: gin.H{"message": "Cards removed from group successfully"},
	})
}

