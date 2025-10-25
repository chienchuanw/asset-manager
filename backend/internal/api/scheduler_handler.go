package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/scheduler"
	"github.com/gin-gonic/gin"
)

// SchedulerHandler 排程器 Handler
type SchedulerHandler struct {
	schedulerManager *scheduler.SchedulerManager
}

// NewSchedulerHandler 建立新的排程器 Handler
func NewSchedulerHandler(schedulerManager *scheduler.SchedulerManager) *SchedulerHandler {
	return &SchedulerHandler{
		schedulerManager: schedulerManager,
	}
}

// GetStatus 取得排程器狀態
// @Summary 取得排程器狀態
// @Description 取得排程器的當前狀態和下次執行時間
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=scheduler.SchedulerStatus}
// @Failure 500 {object} APIResponse
// @Router /api/scheduler/status [get]
func (h *SchedulerHandler) GetStatus(c *gin.Context) {
	status := h.schedulerManager.GetStatus()
	c.JSON(http.StatusOK, APIResponse{
		Data: status,
	})
}

// TriggerSnapshot 手動觸發快照任務
// @Summary 手動觸發快照任務
// @Description 立即執行每日快照任務
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=string}
// @Failure 500 {object} APIResponse
// @Router /api/scheduler/trigger/snapshot [post]
func (h *SchedulerHandler) TriggerSnapshot(c *gin.Context) {
	if err := h.schedulerManager.RunSnapshotNow(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "TRIGGER_SNAPSHOT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Snapshot task triggered successfully",
	})
}

// TriggerDiscordReport 手動觸發 Discord 報告任務
// @Summary 手動觸發 Discord 報告任務
// @Description 立即執行每日 Discord 報告任務
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=string}
// @Failure 500 {object} APIResponse
// @Router /api/scheduler/trigger/discord-report [post]
func (h *SchedulerHandler) TriggerDiscordReport(c *gin.Context) {
	if err := h.schedulerManager.RunDiscordReportNow(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "TRIGGER_DISCORD_REPORT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Discord report task triggered successfully",
	})
}

// ReloadDiscordSchedule 重新載入 Discord 排程
// @Summary 重新載入 Discord 排程
// @Description 當 Discord 設定變更時，重新載入排程
// @Tags scheduler
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=string}
// @Failure 500 {object} APIResponse
// @Router /api/scheduler/reload/discord [post]
func (h *SchedulerHandler) ReloadDiscordSchedule(c *gin.Context) {
	if err := h.schedulerManager.ReloadDiscordSchedule(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "RELOAD_DISCORD_SCHEDULE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Discord schedule reloaded successfully",
	})
}

