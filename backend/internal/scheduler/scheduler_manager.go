package scheduler

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/robfig/cron/v3"
)

// SchedulerManager 排程器管理器
// 統一管理所有排程任務（快照、Discord 報告、每日扣款等）
type SchedulerManager struct {
	cron                     *cron.Cron
	snapshotService          service.AssetSnapshotService
	discordService           service.DiscordService
	settingsService          service.SettingsService
	holdingService           service.HoldingService
	rebalanceService         service.RebalanceService
	billingService           service.BillingService
	exchangeRateService      service.ExchangeRateService
	creditCardService        service.CreditCardService // 新增信用卡服務
	cashFlowService          service.CashFlowService   // 新增現金流服務
	cashFlowReportLogRepo    repository.CashFlowReportLogRepository
	schedulerLogRepo         repository.SchedulerLogRepository
	enabled                  bool
	dailySnapshotTime        string // 格式: "HH:MM" (例如: "23:59")
	discordReportTime        string // 格式: "HH:MM" (例如: "09:00")
	dailyBillingTime         string // 格式: "HH:MM" (例如: "00:01")
	creditCardReminderTime   string // 格式: "HH:MM" (例如: "09:00")
	discordEnabled           bool
	mu                       sync.RWMutex
	snapshotJobID            cron.EntryID
	discordReportJobID       cron.EntryID
	dailyBillingJobID        cron.EntryID
	creditCardReminderJobID  cron.EntryID // 新增信用卡提醒任務 ID
	monthlyReportJobID       cron.EntryID // 月度現金流報告任務 ID
	yearlyReportJobID        cron.EntryID // 年度現金流報告任務 ID
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
	exchangeRateService service.ExchangeRateService,
	creditCardService service.CreditCardService,
	cashFlowService service.CashFlowService,
	schedulerLogRepo repository.SchedulerLogRepository,
	cashFlowReportLogRepo repository.CashFlowReportLogRepository,
	config SchedulerManagerConfig,
) *SchedulerManager {
	return &SchedulerManager{
		cron:                   cron.New(),
		snapshotService:        snapshotService,
		discordService:         discordService,
		settingsService:        settingsService,
		holdingService:         holdingService,
		rebalanceService:       rebalanceService,
		billingService:         billingService,
		exchangeRateService:    exchangeRateService,
		creditCardService:      creditCardService,
		cashFlowService:        cashFlowService,
		schedulerLogRepo:       schedulerLogRepo,
		cashFlowReportLogRepo:  cashFlowReportLogRepo,
		enabled:                config.Enabled,
		dailySnapshotTime:      config.DailySnapshotTime,
		dailyBillingTime:       "00:01", // 預設在每天 00:01 執行扣款
		creditCardReminderTime: "09:00", // 預設在每天 09:00 執行信用卡提醒
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

	// 啟動信用卡繳款提醒排程
	if err := m.startCreditCardReminderSchedule(); err != nil {
		log.Printf("Warning: Failed to start credit card reminder schedule: %v", err)
		// 不返回錯誤，因為提醒排程是可選的
	}

	// 啟動月度現金流報告排程
	if err := m.startMonthlyCashFlowReportSchedule(); err != nil {
		log.Printf("Warning: Failed to start monthly cash flow report schedule: %v", err)
		// 不返回錯誤，因為報告排程是可選的
	}

	// 啟動年度現金流報告排程
	if err := m.startYearlyCashFlowReportSchedule(); err != nil {
		log.Printf("Warning: Failed to start yearly cash flow report schedule: %v", err)
		// 不返回錯誤，因為報告排程是可選的
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
		startTime := time.Now()
		log.Printf("[%s] Running daily snapshot task...", startTime.Format("2006-01-02 15:04:05"))

		var taskErr error

		// 先更新今日匯率
		if m.exchangeRateService != nil {
			if err := m.exchangeRateService.RefreshTodayRate(); err != nil {
				log.Printf("Warning: Failed to refresh exchange rate: %v", err)
				// 不中斷流程，繼續建立快照
			} else {
				log.Println("Exchange rate refreshed successfully")
			}
		}

		// 建立每日快照
		if err := m.snapshotService.CreateDailySnapshots(); err != nil {
			log.Printf("Error creating daily snapshots: %v", err)
			taskErr = err
			m.sendFailureNotification("每日資產快照", err)
		} else {
			log.Println("Daily snapshots created successfully")
		}

		// 記錄執行結果
		m.logTaskExecution("daily_snapshot", startTime, taskErr)
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
		startTime := time.Now()
		log.Printf("[%s] Running daily Discord report task...", startTime.Format("2006-01-02 15:04:05"))

		taskErr := m.sendDailyDiscordReport()
		if taskErr != nil {
			log.Printf("Error sending Discord report: %v", taskErr)
			m.sendFailureNotification("Discord 每日報告", taskErr)
		} else {
			log.Println("Discord report sent successfully")
		}

		// 記錄執行結果
		m.logTaskExecution("discord_report", startTime, taskErr)
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

// startCreditCardReminderSchedule 啟動信用卡繳款提醒排程
func (m *SchedulerManager) startCreditCardReminderSchedule() error {
	// 解析時間
	hour, minute, err := parseTime(m.creditCardReminderTime)
	if err != nil {
		return fmt.Errorf("invalid credit card reminder time: %w", err)
	}

	// 建立 cron 表達式 (每天指定時間執行)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 註冊信用卡提醒任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		log.Println("Starting credit card payment reminder task...")
		if err := m.sendCreditCardPaymentReminder(); err != nil {
			log.Printf("Error in credit card payment reminder: %v", err)
			// 記錄失敗日誌
			if m.schedulerLogRepo != nil {
				now := time.Now()
				errMsg := fmt.Sprintf("Failed to send credit card payment reminder: %v", err)
				logEntry := &models.SchedulerLog{
					TaskName:     "credit_card_reminder",
					Status:       "failed",
					ErrorMessage: &errMsg,
					StartedAt:    now,
					CompletedAt:  &now,
				}
				if createErr := m.schedulerLogRepo.Create(logEntry); createErr != nil {
					log.Printf("Failed to create scheduler log: %v", createErr)
				}
			}
		} else {
			log.Println("Credit card payment reminder task completed successfully")
			// 記錄成功日誌
			if m.schedulerLogRepo != nil {
				now := time.Now()
				logEntry := &models.SchedulerLog{
					TaskName:    "credit_card_reminder",
					Status:      "success",
					StartedAt:   now,
					CompletedAt: &now,
				}
				if createErr := m.schedulerLogRepo.Create(logEntry); createErr != nil {
					log.Printf("Failed to create scheduler log: %v", createErr)
				}
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to add credit card reminder job: %w", err)
	}

	m.mu.Lock()
	m.creditCardReminderJobID = jobID
	m.mu.Unlock()

	log.Printf("Credit card payment reminder scheduled at %s (cron: %s)", m.creditCardReminderTime, cronExpr)
	return nil
}

// sendCreditCardPaymentReminder 發送信用卡繳款提醒
func (m *SchedulerManager) sendCreditCardPaymentReminder() error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		log.Println("Discord is disabled, skipping credit card payment reminder")
		return nil
	}

	// 取得明天需要繳款的信用卡
	creditCards, err := m.creditCardService.GetTomorrowPaymentDue()
	if err != nil {
		return fmt.Errorf("failed to get tomorrow payment due credit cards: %w", err)
	}

	// 如果沒有需要提醒的信用卡，直接返回
	if len(creditCards) == 0 {
		log.Println("No credit cards need payment reminder tomorrow")
		return nil
	}

	// 發送 Discord 提醒
	if err := m.discordService.SendCreditCardPaymentReminder(settings.Discord.WebhookURL, creditCards); err != nil {
		return fmt.Errorf("failed to send credit card payment reminder: %w", err)
	}

	log.Printf("Credit card payment reminder sent for %d card(s)", len(creditCards))
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

	// 先更新今日匯率
	if m.exchangeRateService != nil {
		if err := m.exchangeRateService.RefreshTodayRate(); err != nil {
			log.Printf("Warning: Failed to refresh exchange rate: %v", err)
			// 不中斷流程，繼續建立快照
		} else {
			log.Println("Exchange rate refreshed successfully")
		}
	}

	// 建立每日快照
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
		startTime := time.Now()
		log.Printf("[%s] Starting daily billing process...", startTime.Format("2006-01-02 15:04:05"))

		var taskErr error

		// 執行每日扣款
		result, err := m.billingService.ProcessDailyBilling(time.Now())
		if err != nil {
			log.Printf("Error processing daily billing: %v", err)
			taskErr = err
			m.sendFailureNotification("每日扣款處理", err)
		} else {
			log.Printf("Daily billing completed: %d subscriptions, %d installments, total: %.2f TWD",
				result.SubscriptionCount,
				result.InstallmentCount,
				result.TotalAmount,
			)

			// 發送 Discord 通知（如果啟用）
			m.sendDailyBillingNotification(result)
		}

		// 記錄執行結果
		m.logTaskExecution("daily_billing", startTime, taskErr)
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

// logTaskExecution 記錄任務執行結果
func (m *SchedulerManager) logTaskExecution(taskName string, startTime time.Time, taskErr error) {
	if m.schedulerLogRepo == nil {
		return // 如果沒有 repository，跳過記錄
	}

	completedAt := time.Now()
	duration := completedAt.Sub(startTime).Seconds()

	status := "success"
	var errorMsg *string
	if taskErr != nil {
		status = "failed"
		errStr := taskErr.Error()
		errorMsg = &errStr
	}

	logEntry := &models.SchedulerLog{
		TaskName:        taskName,
		Status:          status,
		ErrorMessage:    errorMsg,
		StartedAt:       startTime,
		CompletedAt:     &completedAt,
		DurationSeconds: &duration,
	}

	if err := m.schedulerLogRepo.Create(logEntry); err != nil {
		// 記錄失敗不應該影響主流程
		log.Printf("Warning: Failed to log scheduler execution: %v\n", err)
	}
}

// sendFailureNotification 發送任務失敗通知到 Discord
func (m *SchedulerManager) sendFailureNotification(taskName string, taskErr error) {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		log.Printf("Error getting settings for failure notification: %v", err)
		return
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled || settings.Discord.WebhookURL == "" {
		return
	}

	// 建立失敗通知訊息
	message := &models.DiscordMessage{
		Content: fmt.Sprintf("⚠️ **排程任務執行失敗**\n\n"+
			"**任務名稱：** %s\n"+
			"**時間：** %s\n"+
			"**錯誤訊息：** %s",
			taskName,
			time.Now().Format("2006-01-02 15:04:05"),
			taskErr.Error(),
		),
	}

	// 發送通知
	if err := m.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		log.Printf("Error sending failure notification: %v", err)
	}
}

// GetTaskSummaries 取得所有排程任務的執行摘要
func (m *SchedulerManager) GetTaskSummaries() ([]models.SchedulerLogSummary, error) {
	if m.schedulerLogRepo == nil {
		return nil, fmt.Errorf("scheduler log repository not available")
	}

	taskNames := []string{"daily_snapshot", "discord_report", "daily_billing"}
	summaries := make([]models.SchedulerLogSummary, 0, len(taskNames))

	for _, taskName := range taskNames {
		summary, err := m.schedulerLogRepo.GetSummaryByTaskName(taskName)
		if err != nil {
			log.Printf("Warning: Failed to get summary for task %s: %v", taskName, err)
			continue
		}
		if summary != nil {
			summaries = append(summaries, *summary)
		}
	}

	return summaries, nil
}

// startMonthlyCashFlowReportSchedule 啟動月度現金流報告排程
func (m *SchedulerManager) startMonthlyCashFlowReportSchedule() error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查是否啟用月度報告
	if !settings.Discord.MonthlyReportEnabled {
		log.Println("Monthly cash flow report is disabled in settings")
		return nil
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		log.Println("Discord is disabled, skipping monthly cash flow report schedule")
		return nil
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		log.Println("Discord webhook URL is not set")
		return nil
	}

	// 驗證日期設定（只允許 1-10 號）
	if settings.Discord.MonthlyReportDay < 1 || settings.Discord.MonthlyReportDay > 10 {
		return fmt.Errorf("invalid monthly report day: %d (must be 1-10)", settings.Discord.MonthlyReportDay)
	}

	// 建立 cron 表達式（每月指定日期的 09:00 執行）
	// 格式: "分 時 日 月 週"
	cronExpr := fmt.Sprintf("0 9 %d * *", settings.Discord.MonthlyReportDay)

	// 註冊月度報告任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		startTime := time.Now()
		log.Printf("[%s] Running monthly cash flow report task...", startTime.Format("2006-01-02 15:04:05"))

		// 計算上個月的年份和月份
		now := time.Now()
		year := now.Year()
		month := int(now.Month()) - 1
		if month < 1 {
			month = 12
			year--
		}

		taskErr := m.sendMonthlyCashFlowReport(year, month)
		if taskErr != nil {
			log.Printf("Error sending monthly cash flow report: %v", taskErr)
			m.sendFailureNotification("月度現金流報告", taskErr)
		} else {
			log.Println("Monthly cash flow report sent successfully")
		}

		// 記錄執行結果
		m.logTaskExecution("monthly_cash_flow_report", startTime, taskErr)
	})
	if err != nil {
		return fmt.Errorf("failed to add monthly cash flow report cron job: %w", err)
	}

	// 更新狀態
	m.mu.Lock()
	m.monthlyReportJobID = jobID
	m.mu.Unlock()

	log.Printf("Monthly cash flow report schedule registered at day %d of each month", settings.Discord.MonthlyReportDay)
	return nil
}

// startYearlyCashFlowReportSchedule 啟動年度現金流報告排程
func (m *SchedulerManager) startYearlyCashFlowReportSchedule() error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查是否啟用年度報告
	if !settings.Discord.YearlyReportEnabled {
		log.Println("Yearly cash flow report is disabled in settings")
		return nil
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		log.Println("Discord is disabled, skipping yearly cash flow report schedule")
		return nil
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		log.Println("Discord webhook URL is not set")
		return nil
	}

	// 驗證日期設定（只允許 1-10 號）
	if settings.Discord.YearlyReportDay < 1 || settings.Discord.YearlyReportDay > 10 {
		return fmt.Errorf("invalid yearly report day: %d (must be 1-10)", settings.Discord.YearlyReportDay)
	}

	// 驗證月份設定（1-12）
	if settings.Discord.YearlyReportMonth < 1 || settings.Discord.YearlyReportMonth > 12 {
		return fmt.Errorf("invalid yearly report month: %d (must be 1-12)", settings.Discord.YearlyReportMonth)
	}

	// 建立 cron 表達式（每年指定月份和日期的 09:00 執行）
	// 格式: "分 時 日 月 週"
	cronExpr := fmt.Sprintf("0 9 %d %d *", settings.Discord.YearlyReportDay, settings.Discord.YearlyReportMonth)

	// 註冊年度報告任務
	jobID, err := m.cron.AddFunc(cronExpr, func() {
		startTime := time.Now()
		log.Printf("[%s] Running yearly cash flow report task...", startTime.Format("2006-01-02 15:04:05"))

		// 計算去年的年份
		year := time.Now().Year() - 1

		taskErr := m.sendYearlyCashFlowReport(year)
		if taskErr != nil {
			log.Printf("Error sending yearly cash flow report: %v", taskErr)
			m.sendFailureNotification("年度現金流報告", taskErr)
		} else {
			log.Println("Yearly cash flow report sent successfully")
		}

		// 記錄執行結果
		m.logTaskExecution("yearly_cash_flow_report", startTime, taskErr)
	})
	if err != nil {
		return fmt.Errorf("failed to add yearly cash flow report cron job: %w", err)
	}

	// 更新狀態
	m.mu.Lock()
	m.yearlyReportJobID = jobID
	m.mu.Unlock()

	log.Printf("Yearly cash flow report schedule registered at %d/%d of each year",
		settings.Discord.YearlyReportMonth, settings.Discord.YearlyReportDay)
	return nil
}

// sendMonthlyCashFlowReport 發送月度現金流報告
func (m *SchedulerManager) sendMonthlyCashFlowReport(year, month int) error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查是否已經成功發送過
	if m.cashFlowReportLogRepo != nil {
		latestLog, err := m.cashFlowReportLogRepo.GetLatestByType(models.CashFlowReportTypeMonthly, year, month)
		if err != nil {
			log.Printf("Warning: Failed to check latest report log: %v", err)
		} else if latestLog != nil && latestLog.Success {
			log.Printf("Monthly report for %d/%d already sent successfully, skipping", year, month)
			return nil
		}
	}

	// 取得月度摘要（包含比較資料）
	summary, err := m.cashFlowService.GetMonthlySummaryWithComparison(year, month)
	if err != nil {
		return m.recordReportFailure(models.CashFlowReportTypeMonthly, year, month, err)
	}

	// 格式化 Discord 訊息
	message := m.discordService.FormatMonthlyCashFlowReport(summary)

	// 發送訊息
	if err := m.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		return m.recordReportFailure(models.CashFlowReportTypeMonthly, year, month, err)
	}

	// 記錄成功發送
	return m.recordReportSuccess(models.CashFlowReportTypeMonthly, year, month)
}

// sendYearlyCashFlowReport 發送年度現金流報告
func (m *SchedulerManager) sendYearlyCashFlowReport(year int) error {
	// 取得設定
	settings, err := m.settingsService.GetSettings()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// 檢查是否已經成功發送過（年度報告的 month 設為 0）
	if m.cashFlowReportLogRepo != nil {
		latestLog, err := m.cashFlowReportLogRepo.GetLatestByType(models.CashFlowReportTypeYearly, year, 0)
		if err != nil {
			log.Printf("Warning: Failed to check latest report log: %v", err)
		} else if latestLog != nil && latestLog.Success {
			log.Printf("Yearly report for %d already sent successfully, skipping", year)
			return nil
		}
	}

	// 取得年度摘要（包含比較資料）
	summary, err := m.cashFlowService.GetYearlySummaryWithComparison(year)
	if err != nil {
		return m.recordReportFailure(models.CashFlowReportTypeYearly, year, 0, err)
	}

	// 格式化 Discord 訊息
	message := m.discordService.FormatYearlyCashFlowReport(summary)

	// 發送訊息
	if err := m.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		return m.recordReportFailure(models.CashFlowReportTypeYearly, year, 0, err)
	}

	// 記錄成功發送
	return m.recordReportSuccess(models.CashFlowReportTypeYearly, year, 0)
}

// recordReportSuccess 記錄報告發送成功
func (m *SchedulerManager) recordReportSuccess(reportType models.CashFlowReportType, year, month int) error {
	if m.cashFlowReportLogRepo == nil {
		return nil
	}

	input := &models.CreateCashFlowReportLogInput{
		ReportType: reportType,
		Year:       year,
		Month:      &month,
		Success:    true,
		RetryCount: 0,
	}

	_, err := m.cashFlowReportLogRepo.Create(input)
	if err != nil {
		log.Printf("Warning: Failed to create report log: %v", err)
	}

	return nil
}

// recordReportFailure 記錄報告發送失敗
func (m *SchedulerManager) recordReportFailure(reportType models.CashFlowReportType, year, month int, sendErr error) error {
	if m.cashFlowReportLogRepo == nil {
		return sendErr
	}

	errMsg := sendErr.Error()
	input := &models.CreateCashFlowReportLogInput{
		ReportType: reportType,
		Year:       year,
		Month:      &month,
		Success:    false,
		ErrorMsg:   &errMsg,
		RetryCount: 0,
	}

	_, err := m.cashFlowReportLogRepo.Create(input)
	if err != nil {
		log.Printf("Warning: Failed to create report log: %v", err)
	}

	return sendErr
}
