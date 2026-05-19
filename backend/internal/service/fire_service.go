package service

import (
	"errors"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
)

// ErrInvalidWithdrawalRate 提領率必須為正數
var ErrInvalidWithdrawalRate = errors.New("withdrawal rate must be greater than zero")

// fireNetWorthSource 提供最新淨值快照（介面隔離，僅取所需方法）
type fireNetWorthSource interface {
	GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error)
}

// fireCashFlowSource 提供現金流摘要（介面隔離，僅取所需方法）
type fireCashFlowSource interface {
	GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error)
}

// FireService FIRE 投影計算服務
type FireService interface {
	GetProjection(input models.FireProjectionInput) (*models.FireProjectionResult, error)
}

type fireService struct {
	netWorthSource fireNetWorthSource
	cashFlowSource fireCashFlowSource
}

// NewFireService 建立 FireService
func NewFireService(netWorthSource fireNetWorthSource, cashFlowSource fireCashFlowSource) FireService {
	return &fireService{
		netWorthSource: netWorthSource,
		cashFlowSource: cashFlowSource,
	}
}

func (s *fireService) GetProjection(input models.FireProjectionInput) (*models.FireProjectionResult, error) {
	expectedReturn := models.DefaultExpectedReturn
	if input.ExpectedReturn != nil {
		expectedReturn = *input.ExpectedReturn
	}

	withdrawalRate := models.DefaultWithdrawalRate
	if input.WithdrawalRate != nil {
		withdrawalRate = *input.WithdrawalRate
	}
	if withdrawalRate <= 0 {
		return nil, ErrInvalidWithdrawalRate
	}

	// 推導現金流相關預設值僅在需要時查詢一次
	var summary *repository.CashFlowSummary
	needSummary := input.AnnualExpenses == nil || input.AnnualSavings == nil
	if needSummary {
		end := time.Now()
		start := end.AddDate(-1, 0, 0)
		// 摘要查詢失敗時退回 0，不阻斷投影
		summary, _ = s.cashFlowSource.GetSummary(start, end)
	}

	annualExpenses := 0.0
	if input.AnnualExpenses != nil {
		annualExpenses = *input.AnnualExpenses
	} else if summary != nil {
		annualExpenses = summary.TotalExpense
	}

	annualSavings := 0.0
	if input.AnnualSavings != nil {
		annualSavings = *input.AnnualSavings
	} else if summary != nil {
		annualSavings = summary.NetCashFlow
	}

	netWorth := 0.0
	if input.NetWorth != nil {
		netWorth = *input.NetWorth
	} else if snapshot, err := s.netWorthSource.GetLatestSnapshot(models.SnapshotAssetTypeTotal); err == nil && snapshot != nil {
		netWorth = snapshot.ValueTWD
	}

	fireNumber := annualExpenses / withdrawalRate

	progressPct := 1.0
	if fireNumber > 0 {
		progressPct = netWorth / fireNumber
	}

	yearsToFI, projection := project(netWorth, fireNumber, annualSavings, expectedReturn)

	return &models.FireProjectionResult{
		CurrentNetWorth: netWorth,
		AnnualExpenses:  annualExpenses,
		AnnualSavings:   annualSavings,
		ExpectedReturn:  expectedReturn,
		WithdrawalRate:  withdrawalRate,
		FireNumber:      fireNumber,
		ProgressPct:     progressPct,
		YearsToFI:       yearsToFI,
		OnTrack:         yearsToFI != nil,
		Projection:      projection,
	}, nil
}

// project 逐年推算淨值，回傳達到 FI 的年數（未達標為 nil）與投影資料點。
func project(netWorth, fireNumber, annualSavings, expectedReturn float64) (*int, []models.FireProjectionPoint) {
	projection := []models.FireProjectionPoint{
		{Year: 0, ProjectedNetWorth: netWorth, FireTarget: fireNumber},
	}

	if netWorth >= fireNumber {
		zero := 0
		return &zero, projection
	}

	// 無正儲蓄則視為無法達標（依設計規格）
	if annualSavings <= 0 {
		return nil, projection
	}

	nw := netWorth
	for year := 1; year <= models.MaxProjectionYears; year++ {
		nw = nw*(1+expectedReturn) + annualSavings
		projection = append(projection, models.FireProjectionPoint{
			Year:              year,
			ProjectedNetWorth: nw,
			FireTarget:        fireNumber,
		})
		if nw >= fireNumber {
			y := year
			return &y, projection
		}
	}

	return nil, projection
}
