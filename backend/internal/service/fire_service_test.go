package service

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFireNetWorthSource mocks the latest-snapshot dependency
type MockFireNetWorthSource struct {
	mock.Mock
}

func (m *MockFireNetWorthSource) GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	args := m.Called(assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

// MockFireCashFlowSource mocks the cash-flow summary dependency
type MockFireCashFlowSource struct {
	mock.Mock
}

func (m *MockFireCashFlowSource) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CashFlowSummary), args.Error(1)
}

func f(v float64) *float64 { return &v }

func TestFireService_GetProjection_NormalCase(t *testing.T) {
	nw := new(MockFireNetWorthSource)
	cf := new(MockFireCashFlowSource)
	svc := NewFireService(nw, cf)

	in := models.FireProjectionInput{
		NetWorth:       f(1_000_000),
		AnnualExpenses: f(600_000),
		AnnualSavings:  f(500_000),
		ExpectedReturn: f(0.05),
		WithdrawalRate: f(0.04),
	}

	res, err := svc.GetProjection(in)

	assert.NoError(t, err)
	// fireNumber = 600000 / 0.04 = 15,000,000
	assert.InDelta(t, 15_000_000, res.FireNumber, 0.01)
	assert.InDelta(t, 1_000_000.0/15_000_000.0, res.ProgressPct, 0.0001)
	assert.NotNil(t, res.YearsToFI)
	assert.True(t, res.OnTrack)
	assert.NotEmpty(t, res.Projection)
	// Echoed inputs
	assert.Equal(t, 1_000_000.0, res.CurrentNetWorth)
	assert.Equal(t, 600_000.0, res.AnnualExpenses)
}

func TestFireService_GetProjection_NotOnTrackWhenSavingsNonPositive(t *testing.T) {
	svc := NewFireService(new(MockFireNetWorthSource), new(MockFireCashFlowSource))

	res, err := svc.GetProjection(models.FireProjectionInput{
		NetWorth:       f(100_000),
		AnnualExpenses: f(500_000),
		AnnualSavings:  f(-1000),
		WithdrawalRate: f(0.04),
	})

	assert.NoError(t, err)
	assert.Nil(t, res.YearsToFI)
	assert.False(t, res.OnTrack)
}

func TestFireService_GetProjection_AlreadyFIWhenExpensesZero(t *testing.T) {
	svc := NewFireService(new(MockFireNetWorthSource), new(MockFireCashFlowSource))

	res, err := svc.GetProjection(models.FireProjectionInput{
		NetWorth:       f(0),
		AnnualExpenses: f(0),
		AnnualSavings:  f(0),
		WithdrawalRate: f(0.04),
	})

	assert.NoError(t, err)
	assert.Equal(t, 0.0, res.FireNumber)
	assert.NotNil(t, res.YearsToFI)
	assert.Equal(t, 0, *res.YearsToFI)
	assert.True(t, res.OnTrack)
}

func TestFireService_GetProjection_ZeroReturnLinearGrowth(t *testing.T) {
	svc := NewFireService(new(MockFireNetWorthSource), new(MockFireCashFlowSource))

	// fireNumber = 100000/0.04 = 2,500,000; nw 500k; savings 500k/yr; r=0
	// years = ceil((2.5M-0.5M)/0.5M) = 4
	res, err := svc.GetProjection(models.FireProjectionInput{
		NetWorth:       f(500_000),
		AnnualExpenses: f(100_000),
		AnnualSavings:  f(500_000),
		ExpectedReturn: f(0),
		WithdrawalRate: f(0.04),
	})

	assert.NoError(t, err)
	assert.NotNil(t, res.YearsToFI)
	assert.Equal(t, 4, *res.YearsToFI)
}

func TestFireService_GetProjection_OverridesTakePrecedenceOverDerived(t *testing.T) {
	nw := new(MockFireNetWorthSource)
	cf := new(MockFireCashFlowSource)
	// Derived values that should be ignored because overrides are supplied
	nw.On("GetLatestSnapshot", models.SnapshotAssetTypeTotal).
		Return(&models.AssetSnapshot{ValueTWD: 999}, nil)
	cf.On("GetSummary", mock.Anything, mock.Anything).
		Return(&repository.CashFlowSummary{TotalExpense: 999, NetCashFlow: 999}, nil)

	svc := NewFireService(nw, cf)
	res, err := svc.GetProjection(models.FireProjectionInput{
		NetWorth:       f(1_000_000),
		AnnualExpenses: f(400_000),
		AnnualSavings:  f(300_000),
		WithdrawalRate: f(0.04),
	})

	assert.NoError(t, err)
	assert.Equal(t, 1_000_000.0, res.CurrentNetWorth)
	assert.Equal(t, 400_000.0, res.AnnualExpenses)
	assert.Equal(t, 300_000.0, res.AnnualSavings)
	nw.AssertNotCalled(t, "GetLatestSnapshot", mock.Anything)
	cf.AssertNotCalled(t, "GetSummary", mock.Anything, mock.Anything)
}

func TestFireService_GetProjection_DerivesFromDependenciesWhenNoOverrides(t *testing.T) {
	nw := new(MockFireNetWorthSource)
	cf := new(MockFireCashFlowSource)
	nw.On("GetLatestSnapshot", models.SnapshotAssetTypeTotal).
		Return(&models.AssetSnapshot{ValueTWD: 2_000_000}, nil)
	cf.On("GetSummary", mock.Anything, mock.Anything).
		Return(&repository.CashFlowSummary{TotalExpense: 480_000, NetCashFlow: 360_000}, nil)

	svc := NewFireService(nw, cf)
	res, err := svc.GetProjection(models.FireProjectionInput{})

	assert.NoError(t, err)
	assert.Equal(t, 2_000_000.0, res.CurrentNetWorth)
	assert.Equal(t, 480_000.0, res.AnnualExpenses)
	assert.Equal(t, 360_000.0, res.AnnualSavings)
	assert.Equal(t, models.DefaultExpectedReturn, res.ExpectedReturn)
	assert.Equal(t, models.DefaultWithdrawalRate, res.WithdrawalRate)
}

func TestFireService_GetProjection_NoSnapshotDefaultsToZeroNetWorth(t *testing.T) {
	nw := new(MockFireNetWorthSource)
	cf := new(MockFireCashFlowSource)
	nw.On("GetLatestSnapshot", models.SnapshotAssetTypeTotal).
		Return(nil, errors.New("no snapshot"))
	cf.On("GetSummary", mock.Anything, mock.Anything).
		Return(&repository.CashFlowSummary{TotalExpense: 300_000, NetCashFlow: 200_000}, nil)

	svc := NewFireService(nw, cf)
	res, err := svc.GetProjection(models.FireProjectionInput{})

	assert.NoError(t, err)
	assert.Equal(t, 0.0, res.CurrentNetWorth)
}

func TestFireService_GetProjection_RejectsNonPositiveWithdrawalRate(t *testing.T) {
	svc := NewFireService(new(MockFireNetWorthSource), new(MockFireCashFlowSource))

	_, err := svc.GetProjection(models.FireProjectionInput{
		AnnualExpenses: f(100_000),
		WithdrawalRate: f(0),
	})

	assert.Error(t, err)
}

func TestFireService_GetProjection_CapsProjectionAt60Years(t *testing.T) {
	svc := NewFireService(new(MockFireNetWorthSource), new(MockFireCashFlowSource))

	// Tiny savings, huge target -> never reached within 60 years
	res, err := svc.GetProjection(models.FireProjectionInput{
		NetWorth:       f(0),
		AnnualExpenses: f(10_000_000),
		AnnualSavings:  f(1),
		ExpectedReturn: f(0),
		WithdrawalRate: f(0.04),
	})

	assert.NoError(t, err)
	assert.Nil(t, res.YearsToFI)
	assert.False(t, res.OnTrack)
	assert.LessOrEqual(t, len(res.Projection), models.MaxProjectionYears+1)
}
