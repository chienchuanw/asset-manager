// Package mcpserver exposes read-only portfolio data to AI agents over MCP.
// Adapters are thin wrappers over the existing service layer (same pattern as
// internal/discord) — they select the read methods the tools need and return
// the domain models, which already carry JSON tags with raw numeric values.
package mcpserver

import (
	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/google/uuid"
)

// --- holdings ---

type holdingReader interface {
	GetAllHoldings(models.HoldingFilters) (*service.HoldingServiceResult, error)
	GetHoldingBySymbol(string) (*models.Holding, error)
}

type HoldingsAdapter struct{ svc holdingReader }

func NewHoldingsAdapter(svc holdingReader) *HoldingsAdapter { return &HoldingsAdapter{svc: svc} }

func (a *HoldingsAdapter) List() ([]*models.Holding, error) {
	res, err := a.svc.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		return nil, err
	}
	return res.Holdings, nil
}

func (a *HoldingsAdapter) Get(symbol string) (*models.Holding, error) {
	return a.svc.GetHoldingBySymbol(symbol)
}

// --- transactions ---

type transactionReader interface {
	ListTransactions(repository.TransactionFilters) ([]*models.Transaction, error)
	GetTransaction(uuid.UUID) (*models.Transaction, error)
}

type TransactionsAdapter struct{ svc transactionReader }

func NewTransactionsAdapter(svc transactionReader) *TransactionsAdapter {
	return &TransactionsAdapter{svc: svc}
}

func (a *TransactionsAdapter) List(f repository.TransactionFilters) ([]*models.Transaction, error) {
	return a.svc.ListTransactions(f)
}

func (a *TransactionsAdapter) Get(id string) (*models.Transaction, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return a.svc.GetTransaction(uid)
}

// --- analytics ---

type analyticsReader interface {
	GetSummary(models.TimeRange) (*models.AnalyticsSummary, error)
	GetTopAssets(models.TimeRange, int) ([]*models.TopAsset, error)
}

type AnalyticsAdapter struct{ svc analyticsReader }

func NewAnalyticsAdapter(svc analyticsReader) *AnalyticsAdapter { return &AnalyticsAdapter{svc: svc} }

type AnalyticsResult struct {
	Summary   *models.AnalyticsSummary `json:"summary"`
	TopAssets []*models.TopAsset       `json:"top_assets"`
}

func (a *AnalyticsAdapter) Get(tr models.TimeRange) (*AnalyticsResult, error) {
	sum, err := a.svc.GetSummary(tr)
	if err != nil {
		return nil, err
	}
	top, err := a.svc.GetTopAssets(tr, 5)
	if err != nil {
		return nil, err
	}
	return &AnalyticsResult{Summary: sum, TopAssets: top}, nil
}

// --- allocation ---

type allocationReader interface {
	GetAllocationByType() ([]models.AllocationByType, error)
}

type AllocationAdapter struct{ svc allocationReader }

func NewAllocationAdapter(svc allocationReader) *AllocationAdapter {
	return &AllocationAdapter{svc: svc}
}

func (a *AllocationAdapter) ByType() ([]models.AllocationByType, error) {
	return a.svc.GetAllocationByType()
}

// --- performance trend ---

type trendReader interface {
	GetLatestTrend(int) ([]models.PerformanceTrendPoint, error)
}

type TrendAdapter struct{ svc trendReader }

func NewTrendAdapter(svc trendReader) *TrendAdapter { return &TrendAdapter{svc: svc} }

func (a *TrendAdapter) Latest(days int) ([]models.PerformanceTrendPoint, error) {
	return a.svc.GetLatestTrend(days)
}
