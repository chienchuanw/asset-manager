package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/robfig/cron/v3"
)

// SnapshotScheduler 資產快照排程器
type SnapshotScheduler struct {
	cron               *cron.Cron
	snapshotService    service.AssetSnapshotService
	enabled            bool
	dailySnapshotTime  string // 格式: "HH:MM" (例如: "23:59")
}

// SnapshotSchedulerConfig 排程器配置
type SnapshotSchedulerConfig struct {
	Enabled           bool   // 是否啟用排程
	DailySnapshotTime string // 每日快照時間 (格式: "HH:MM")
}

// NewSnapshotScheduler 建立新的快照排程器
func NewSnapshotScheduler(snapshotService service.AssetSnapshotService, config SnapshotSchedulerConfig) *SnapshotScheduler {
	return &SnapshotScheduler{
		cron:              cron.New(),
		snapshotService:   snapshotService,
		enabled:           config.Enabled,
		dailySnapshotTime: config.DailySnapshotTime,
	}
}

// Start 啟動排程器
func (s *SnapshotScheduler) Start() error {
	if !s.enabled {
		log.Println("Snapshot scheduler is disabled")
		return nil
	}

	// 解析時間
	hour, minute, err := parseTime(s.dailySnapshotTime)
	if err != nil {
		return fmt.Errorf("invalid daily snapshot time: %w", err)
	}

	// 建立 cron 表達式 (每天指定時間執行)
	cronExpr := fmt.Sprintf("%d %d * * *", minute, hour)

	// 註冊每日快照任務
	_, err = s.cron.AddFunc(cronExpr, func() {
		log.Println("Running daily snapshot task...")
		if err := s.snapshotService.CreateDailySnapshots(); err != nil {
			log.Printf("Error creating daily snapshots: %v", err)
		} else {
			log.Println("Daily snapshots created successfully")
		}
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	// 啟動 cron
	s.cron.Start()
	log.Printf("Snapshot scheduler started. Daily snapshots will be created at %s", s.dailySnapshotTime)

	return nil
}

// Stop 停止排程器
func (s *SnapshotScheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("Snapshot scheduler stopped")
	}
}

// RunNow 立即執行快照任務（用於測試或手動觸發）
func (s *SnapshotScheduler) RunNow() error {
	log.Println("Manually triggering snapshot task...")
	if err := s.snapshotService.CreateDailySnapshots(); err != nil {
		return fmt.Errorf("failed to create snapshots: %w", err)
	}
	log.Println("Snapshots created successfully")
	return nil
}

// parseTime 解析時間字串 (格式: "HH:MM")
func parseTime(timeStr string) (hour, minute int, err error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return 0, 0, err
	}
	return t.Hour(), t.Minute(), nil
}

