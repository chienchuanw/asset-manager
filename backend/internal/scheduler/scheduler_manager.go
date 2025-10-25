package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/robfig/cron/v3"
)

// SchedulerManager 排程器管理器
// 統一管理所有排程任務（快照、Discord 報告、每日扣款等）
type SchedulerManager struct {
	cron                *cron.Cron
	snapshotService     service.AssetSnapshotService
	discordService      service.DiscordService
	settingsService     service.SettingsService
	holdingService      service.HoldingService
	rebalanceService    service.RebalanceService
	billingService      service.BillingService
	enabled             bool
	dailySnapshotTime   string // 格式: "HH:MM" (例如: "23:59")
	discordReportTime   string // 格式: "HH:MM" (例如: "09:00")
	dailyBillingTime    string // 格式: "HH:MM" (例如: "00:01")
	discordEnabled      bool
	mu                  sync.RWMutex
	snapshotJobID       cron.EntryID
	discordReportJobID  cron.EntryID
	dailyBillingJobID   cron.EntryID
}

// SchedulerManagerConfig 排程器管理器配置
type SchedulerManagerConfig struct {
	Enabled           bool   // 是否啟用排程
	DailySnapshotTime string // 每日快照時間 (格式: "HH:MM")
}

// NewSchedulerManager 建立新的排程器管理器
func NewSchedulerManager(
	snapshotService service.AssetSnapshotService,
	discordService service.DiscordService,
	settingsService service.SettingsService,
	holdingService service.HoldingService,
	rebalanceService service.RebalanceService,
	billingService service.BillingService,
	config SchedulerManagerConfig,
) *SchedulerManager {
	return &SchedulerManager{
		cron:              cron.New(),
		snapshotService:   snapshotService,
		discordService:    discordService,
		settingsService:   settingsService,
		holdingService:    holdingService,
		rebalanceService:  rebalanceService,
		billingService:    billingService,
		enabled:           config.Enabled,
		dailySnapshotTime: config.DailySnapshotTime,
		dailyBillingTime:  "00:01", // 預設在每天 00:01 執行扣款
	}
}

// Start 啟動排程器
func (m *SchedulerManager) Start() error {
	if !m.enabled {
		log.Println("Scheduler manager is disabled")
		return nil
	}

	// 啟動每日快照排程
	if err := m.startSnapshotSchedule(); err != nil {
		return fmt.Errorf("failed to start snapshot schedule: %w", err)
	}

	// 啟動 Discord 報告排程
	if err := m.startDiscordReportSchedule(); err != nil {
		log.Printf("Warning: Failed to start Discord report schedule: %v", err)
		// 不返回錯誤，因為 Discord 報告是可選的
	}

	// 啟動每日扣款排程
	if err := m.startDailyBillingSchedule(); err != nil {
		log.Printf("Warning: Failed to start daily billing schedule: %v", err)
		// 不返回錯誤，因為扣款排程是可選的
	}

	// 啟動 cron
	m.cron.Start()
	log.Println("Scheduler manager started successfully")

	return nil
}

// startSnapshotSchedule 啟動每日快照排程
func (m *SchedulerManager) startSnapshotSchedule() error {
	// 解析時間
	hour, minute, err := parseTime(m.dailySnapshotTime)
	if err != nil {
		return fmt.Errorf("invalid daily snapshot time: %w", err)
	}

	// 建立 cron 表達式 (每天指定時間執行)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 註冊每日快照任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		log.Println("Running daily snapshot task...")
		if err := m.snapshotService.CreateDailySnapshots(); err != nil {
			log.Printf("Error creating daily snapshots: %v", err)
		} else {
			log.Println("Daily snapshots created successfully")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add snapshot cron job: %w", err)
	}

	m.mu.Lock()
	m.snapshotJobID = jobID
	m.mu.Unlock()

	log.Printf("Daily snapshot schedule registered at %s", m.dailySnapshotTime)
	return nil
}

// startDiscordReportSchedule 啟動 Discord 報告排程
func (m *SchedulerManager) startDiscordReportSchedule() error {
	// 取得 Discord 設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		log.Println("Discord report is disabled in settings")
		return nil
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		log.Println("Discord webhook URL is not set")
		return nil
	}

	// 解析時間
	hour, minute, err := parseTime(settings.Discord.ReportTime)
	if err != nil {
		return fmt.Errorf("invalid Discord report time: %w", err)
	}

	// 建立 cron 表達式 (每天指定時間執行)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 註冊 Discord 報告任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		log.Println("Running daily Discord report task...")
		if err := m.sendDailyDiscordReport(); err != nil {
			log.Printf("Error sending Discord report: %v", err)
		} else {
			log.Println("Discord report sent successfully")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add Discord report cron job: %w", err)
	}

	// 更新狀態（使用鎖保護）
	m.mu.Lock()
	m.discordReportJobID = jobID
	m.discordEnabled = true
	m.discordReportTime = settings.Discord.ReportTime
	m.mu.Unlock()

	log.Printf("Discord report schedule registered at %s", settings.Discord.ReportTime)
	return nil
}

// startDiscordReportScheduleUnsafe 啟動 Discord 報告排程（不使用鎖，供內部使用）
// 注意：調用此方法前必須已經獲取鎖
func (m *SchedulerManager) startDiscordReportScheduleUnsafe() error {
	// 取得 Discord 設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		log.Println("Discord report is disabled in settings")
		return nil
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		log.Println("Discord webhook URL is not set")
		return nil
	}

	// 解析時間
	hour, minute, err := parseTime(settings.Discord.ReportTime)
	if err != nil {
		return fmt.Errorf("invalid Discord report time: %w", err)
	}

	// 建立 cron 表達式 (每天指定時間執行)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 註冊 Discord 報告任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		log.Println("Running daily Discord report task...")
		if err := m.sendDailyDiscordReport(); err != nil {
			log.Printf("Error sending Discord report: %v", err)
		} else {
			log.Println("Discord report sent successfully")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add Discord report cron job: %w", err)
	}

	// 更新狀態（不使用鎖，因為調用者已經獲取鎖）
	m.discordReportJobID = jobID
	m.discordEnabled = true
	m.discordReportTime = settings.Discord.ReportTime

	log.Printf("Discord report schedule registered at %s", settings.Discord.ReportTime)
	return nil
}

// sendDailyDiscordReport 發送每日 Discord 報告
func (m *SchedulerManager) sendDailyDiscordReport() error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 再次檢查 Discord 是否啟用（設定可能已變更）
	if !settings.Discord.Enabled {
		log.Println("Discord is disabled, skipping report")
		return nil
	}

	// 取得所有持倉
	holdings, err := m.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return fmt.Errorf("failed to get holdings: %w", err)
	}

	// 計算報告資料（與 discord_handler.go 相同的邏輯）
	reportData := m.buildReportData(holdings)

	// 檢查是否需要再平衡
	rebalanceCheck, err := m.rebalanceService.CheckRebalance()
	if err != nil {
		log.Printf("Warning: Failed to check rebalance: %v", err)
		// 不返回錯誤，繼續發送報告
	} else {
		// 將再平衡檢查結果加入報告
		reportData.RebalanceCheck = rebalanceCheck
	}

	// 格式化報告
	message := m.discordService.FormatDailyReport(reportData)

	// 發送訊息
	if err := m.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		return fmt.Errorf("failed to send Discord message: %w", err)
	}

	return nil
}

// buildReportData 建立報告資料
func (m *SchedulerManager) buildReportData(holdings []*models.Holding) *models.DailyReportData {
	var totalMarketValue, totalCost, totalUnrealizedPL float64
	byAssetType := make(map[string]*models.AssetTypePerformance)

	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
		totalCost += holding.TotalCost
		totalUnrealizedPL += holding.UnrealizedPL

		// 按資產類型分類
		assetTypeStr := string(holding.AssetType)
		if _, exists := byAssetType[assetTypeStr]; !exists {
			byAssetType[assetTypeStr] = &models.AssetTypePerformance{
				AssetType: assetTypeStr,
			}
		}
		perf := byAssetType[assetTypeStr]
		perf.MarketValue += holding.MarketValue
		perf.Cost += holding.TotalCost
		perf.UnrealizedPL += holding.UnrealizedPL
		perf.HoldingCount++
	}

	// 計算各資產類型的損益百分比
	for _, perf := range byAssetType {
		if perf.Cost > 0 {
			perf.UnrealizedPct = (perf.UnrealizedPL / perf.Cost) * 100
		}
	}

	// 計算總損益百分比
	totalUnrealizedPct := 0.0
	if totalCost > 0 {
		totalUnrealizedPct = (totalUnrealizedPL / totalCost) * 100
	}

	// 排序持倉（按市值降序）
	// 注意：這裡需要複製 holdings 以避免修改原始資料
	sortedHoldings := make([]*models.Holding, len(holdings))
	copy(sortedHoldings, holdings)

	// 簡單的冒泡排序（因為持倉數量通常不多）
	for i := 0; i < len(sortedHoldings); i++ {
		for j := i + 1; j < len(sortedHoldings); j++ {
			if sortedHoldings[i].MarketValue < sortedHoldings[j].MarketValue {
				sortedHoldings[i], sortedHoldings[j] = sortedHoldings[j], sortedHoldings[i]
			}
		}
	}

	// 取前 5 大持倉
	topHoldings := sortedHoldings
	if len(topHoldings) > 5 {
		topHoldings = topHoldings[:5]
	}

	return &models.DailyReportData{
		Date:               time.Now(),
		TotalMarketValue:   totalMarketValue,
		TotalCost:          totalCost,
		TotalUnrealizedPL:  totalUnrealizedPL,
		TotalUnrealizedPct: totalUnrealizedPct,
		HoldingCount:       len(holdings),
		TopHoldings:        topHoldings,
		ByAssetType:        byAssetType,
	}
}

// Stop 停止排程器
func (m *SchedulerManager) Stop() {
	if m.cron != nil {
		m.cron.Stop()
		log.Println("Scheduler manager stopped")
	}
}

// RunSnapshotNow 立即執行快照任務（用於測試或手動觸發）
func (m *SchedulerManager) RunSnapshotNow() error {
	log.Println("Manually triggering snapshot task...")
	if err := m.snapshotService.CreateDailySnapshots(); err != nil {
		return fmt.Errorf("failed to create snapshots: %w", err)
	}
	log.Println("Snapshots created successfully")
	return nil
}

// RunDiscordReportNow 立即執行 Discord 報告任務（用於測試或手動觸發）
func (m *SchedulerManager) RunDiscordReportNow() error {
	log.Println("Manually triggering Discord report task...")
	if err := m.sendDailyDiscordReport(); err != nil {
		return fmt.Errorf("failed to send Discord report: %w", err)
	}
	log.Println("Discord report sent successfully")
	return nil
}

// ReloadDiscordSchedule 重新載入 Discord 排程（當設定變更時使用）
func (m *SchedulerManager) ReloadDiscordSchedule() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 移除舊的 Discord 報告任務
	if m.discordReportJobID != 0 {
		m.cron.Remove(m.discordReportJobID)
		m.discordReportJobID = 0
		m.discordEnabled = false
		log.Println("Removed old Discord report schedule")
	}

	// 重新啟動 Discord 報告排程（使用不帶鎖的版本）
	if err := m.startDiscordReportScheduleUnsafe(); err != nil {
		return fmt.Errorf("failed to reload Discord schedule: %w", err)
	}

	return nil
}

// GetStatus 取得排程器狀態
func (m *SchedulerManager) GetStatus() SchedulerStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return SchedulerStatus{
		Enabled:            m.enabled,
		SnapshotTime:       m.dailySnapshotTime,
		DiscordEnabled:     m.discordEnabled,
		DiscordReportTime:  m.discordReportTime,
		NextSnapshotRun:    m.getNextRunTime(m.snapshotJobID),
		NextDiscordRun:     m.getNextRunTime(m.discordReportJobID),
	}
}

// getNextRunTime 取得下次執行時間
func (m *SchedulerManager) getNextRunTime(jobID cron.EntryID) time.Time {
	if jobID == 0 {
		return time.Time{}
	}
	entry := m.cron.Entry(jobID)
	return entry.Next
}

// startDailyBillingSchedule 啟動每日扣款排程
func (m *SchedulerManager) startDailyBillingSchedule() error {
	// 解析時間
	hour, minute, err := parseTime(m.dailyBillingTime)
	if err != nil {
		return fmt.Errorf("invalid daily billing time: %w", err)
	}

	// 建立 cron 表達式（每天在指定時間執行）
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 新增排程任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		log.Println("Starting daily billing process...")

		// 執行每日扣款
		result, err := m.billingService.ProcessDailyBilling(time.Now())
		if err != nil {
			log.Printf("Error processing daily billing: %v", err)
			return
		}

		log.Printf("Daily billing completed: %d subscriptions, %d installments, total: %.2f TWD",
			result.SubscriptionCount,
			result.InstallmentCount,
			result.TotalAmount,
		)

		// 發送 Discord 通知（如果啟用）
		m.sendDailyBillingNotification(result)
	})

	if err != nil {
		return fmt.Errorf("failed to add daily billing job: %w", err)
	}

	m.mu.Lock()
	m.dailyBillingJobID = jobID
	m.mu.Unlock()

	log.Printf("Daily billing schedule started (will run at %s)", m.dailyBillingTime)
	return nil
}

// sendDailyBillingNotification 發送每日扣款通知到 Discord
func (m *SchedulerManager) sendDailyBillingNotification(result *service.DailyBillingResult) {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		log.Printf("Error getting settings for billing notification: %v", err)
		return
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		return
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		return
	}

	// 檢查是否啟用每日扣款通知
	// 這裡假設 settings 中有 notification_daily_billing 設定
	// 如果沒有，可以直接發送通知

	// 發送通知
	if err := m.discordService.SendDailyBillingNotification(settings.Discord.WebhookURL, result); err != nil {
		log.Printf("Error sending daily billing notification: %v", err)
	} else {
		log.Println("Daily billing notification sent to Discord")
	}
}

// SchedulerStatus 排程器狀態
type SchedulerStatus struct {
	Enabled           bool      `json:"enabled"`
	SnapshotTime      string    `json:"snapshot_time"`
	DiscordEnabled    bool      `json:"discord_enabled"`
	DiscordReportTime string    `json:"discord_report_time"`
	NextSnapshotRun   time.Time `json:"next_snapshot_run"`
	NextDiscordRun    time.Time `json:"next_discord_run"`
}

