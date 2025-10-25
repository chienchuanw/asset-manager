package service

import (
	"fmt"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// SettingsService 設定服務介面
type SettingsService interface {
	// GetSettings 取得所有設定（群組格式）
	GetSettings() (*models.SettingsGroup, error)
	// UpdateSettings 更新設定（群組格式）
	UpdateSettings(input *models.UpdateSettingsGroupInput) (*models.SettingsGroup, error)
}

type settingsService struct {
	repo repository.SettingsRepository
}

// NewSettingsService 建立設定服務
func NewSettingsService(repo repository.SettingsRepository) SettingsService {
	return &settingsService{repo: repo}
}

// GetSettings 取得所有設定（群組格式）
func (s *settingsService) GetSettings() (*models.SettingsGroup, error) {
	// 取得所有設定
	settings, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	// 建立設定 map
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	// 轉換為群組格式
	group := &models.SettingsGroup{
		Discord: models.DiscordSettings{
			WebhookURL: settingsMap["discord_webhook_url"],
			Enabled:    settingsMap["discord_enabled"] == "true",
			ReportTime: settingsMap["discord_report_time"],
		},
		Allocation: models.AllocationSettings{
			TWStock:            parseFloat(settingsMap["target_allocation_tw_stock"]),
			USStock:            parseFloat(settingsMap["target_allocation_us_stock"]),
			Crypto:             parseFloat(settingsMap["target_allocation_crypto"]),
			RebalanceThreshold: parseFloat(settingsMap["rebalance_threshold"]),
		},
	}

	return group, nil
}

// UpdateSettings 更新設定（群組格式）
func (s *settingsService) UpdateSettings(input *models.UpdateSettingsGroupInput) (*models.SettingsGroup, error) {
	// 更新 Discord 設定
	if input.Discord != nil {
		if err := s.updateDiscordSettings(input.Discord); err != nil {
			return nil, err
		}
	}

	// 更新資產配置設定
	if input.Allocation != nil {
		if err := s.updateAllocationSettings(input.Allocation); err != nil {
			return nil, err
		}
	}

	// 回傳更新後的設定
	return s.GetSettings()
}

// updateDiscordSettings 更新 Discord 設定
func (s *settingsService) updateDiscordSettings(discord *models.DiscordSettings) error {
	// 更新 webhook URL
	if _, err := s.repo.Update("discord_webhook_url", &models.UpdateSettingInput{
		Value: discord.WebhookURL,
	}); err != nil {
		return fmt.Errorf("failed to update discord_webhook_url: %w", err)
	}

	// 更新 enabled
	enabledValue := "false"
	if discord.Enabled {
		enabledValue = "true"
	}
	if _, err := s.repo.Update("discord_enabled", &models.UpdateSettingInput{
		Value: enabledValue,
	}); err != nil {
		return fmt.Errorf("failed to update discord_enabled: %w", err)
	}

	// 更新 report time
	if _, err := s.repo.Update("discord_report_time", &models.UpdateSettingInput{
		Value: discord.ReportTime,
	}); err != nil {
		return fmt.Errorf("failed to update discord_report_time: %w", err)
	}

	return nil
}

// updateAllocationSettings 更新資產配置設定
func (s *settingsService) updateAllocationSettings(allocation *models.AllocationSettings) error {
	// 驗證總和是否為 100%
	total := allocation.TWStock + allocation.USStock + allocation.Crypto
	if total != 100 {
		return fmt.Errorf("allocation percentages must sum to 100, got %.2f", total)
	}

	// 更新台股配置
	if _, err := s.repo.Update("target_allocation_tw_stock", &models.UpdateSettingInput{
		Value: fmt.Sprintf("%.2f", allocation.TWStock),
	}); err != nil {
		return fmt.Errorf("failed to update target_allocation_tw_stock: %w", err)
	}

	// 更新美股配置
	if _, err := s.repo.Update("target_allocation_us_stock", &models.UpdateSettingInput{
		Value: fmt.Sprintf("%.2f", allocation.USStock),
	}); err != nil {
		return fmt.Errorf("failed to update target_allocation_us_stock: %w", err)
	}

	// 更新加密貨幣配置
	if _, err := s.repo.Update("target_allocation_crypto", &models.UpdateSettingInput{
		Value: fmt.Sprintf("%.2f", allocation.Crypto),
	}); err != nil {
		return fmt.Errorf("failed to update target_allocation_crypto: %w", err)
	}

	// 更新再平衡閾值
	if _, err := s.repo.Update("rebalance_threshold", &models.UpdateSettingInput{
		Value: fmt.Sprintf("%.2f", allocation.RebalanceThreshold),
	}); err != nil {
		return fmt.Errorf("failed to update rebalance_threshold: %w", err)
	}

	return nil
}

// parseFloat 解析浮點數字串
func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

