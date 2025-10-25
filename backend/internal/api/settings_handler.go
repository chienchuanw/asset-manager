package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// SettingsHandler 設定 Handler
type SettingsHandler struct {
	service service.SettingsService
}

// NewSettingsHandler 建立設定 Handler
func NewSettingsHandler(service service.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		service: service,
	}
}

// GetSettings 取得所有設定
// @Summary 取得所有設定
// @Description 取得所有設定（群組格式）
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=models.SettingsGroup}
// @Failure 500 {object} APIResponse
// @Router /api/settings [get]
func (h *SettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_SETTINGS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: settings,
	})
}

// UpdateSettings 更新設定
// @Summary 更新設定
// @Description 更新設定（群組格式）
// @Tags settings
// @Accept json
// @Produce json
// @Param input body models.UpdateSettingsGroupInput true "更新設定輸入"
// @Success 200 {object} APIResponse{data=models.SettingsGroup}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/settings [put]
func (h *SettingsHandler) UpdateSettings(c *gin.Context) {
	var input models.UpdateSettingsGroupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	settings, err := h.service.UpdateSettings(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "UPDATE_SETTINGS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: settings,
	})
}

