package discord

import (
	"errors"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type mockCashFlowQueryService struct {
	monthlySummary *models.MonthlyCashFlowSummary
	err            error
	calledYear     int
	calledMonth    int
}

func (m *mockCashFlowQueryService) CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	panic("unexpected call to CreateCashFlow")
}

func (m *mockCashFlowQueryService) GetCashFlow(id uuid.UUID) (*models.CashFlow, error) {
	panic("unexpected call to GetCashFlow")
}

func (m *mockCashFlowQueryService) ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	panic("unexpected call to ListCashFlows")
}

func (m *mockCashFlowQueryService) UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	panic("unexpected call to UpdateCashFlow")
}

func (m *mockCashFlowQueryService) DeleteCashFlow(id uuid.UUID) error {
	panic("unexpected call to DeleteCashFlow")
}

func (m *mockCashFlowQueryService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	panic("unexpected call to GetSummary")
}

func (m *mockCashFlowQueryService) GetMonthlySummaryWithComparison(year, month int) (*models.MonthlyCashFlowSummary, error) {
	m.calledYear = year
	m.calledMonth = month
	return m.monthlySummary, m.err
}

func (m *mockCashFlowQueryService) GetYearlySummaryWithComparison(year int) (*models.YearlyCashFlowSummary, error) {
	panic("unexpected call to GetYearlySummaryWithComparison")
}

type mockBankAccountQueryRepo struct {
	accounts []*models.BankAccount
	err      error
}

func (m *mockBankAccountQueryRepo) Create(input *models.CreateBankAccountInput) (*models.BankAccount, error) {
	panic("unexpected call to Create")
}

func (m *mockBankAccountQueryRepo) GetByID(id uuid.UUID) (*models.BankAccount, error) {
	panic("unexpected call to GetByID")
}

func (m *mockBankAccountQueryRepo) GetAll(currency *models.Currency) ([]*models.BankAccount, error) {
	return m.accounts, m.err
}

func (m *mockBankAccountQueryRepo) Update(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error) {
	panic("unexpected call to Update")
}

func (m *mockBankAccountQueryRepo) UpdateBalance(id uuid.UUID, amount float64) (*models.BankAccount, error) {
	panic("unexpected call to UpdateBalance")
}

func (m *mockBankAccountQueryRepo) Delete(id uuid.UUID) error {
	panic("unexpected call to Delete")
}

type mockCreditCardQueryRepo struct {
	cards []*models.CreditCard
	err   error
}

func (m *mockCreditCardQueryRepo) Create(input *models.CreateCreditCardInput) (*models.CreditCard, error) {
	panic("unexpected call to Create")
}

func (m *mockCreditCardQueryRepo) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	panic("unexpected call to GetByID")
}

func (m *mockCreditCardQueryRepo) GetAll() ([]*models.CreditCard, error) {
	return m.cards, m.err
}

func (m *mockCreditCardQueryRepo) GetByBillingDay(day int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetByBillingDay")
}

func (m *mockCreditCardQueryRepo) GetByPaymentDueDay(day int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetByPaymentDueDay")
}

func (m *mockCreditCardQueryRepo) GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetUpcomingBilling")
}

func (m *mockCreditCardQueryRepo) GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetUpcomingPayment")
}

func (m *mockCreditCardQueryRepo) Update(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error) {
	panic("unexpected call to Update")
}

func (m *mockCreditCardQueryRepo) UpdateUsedCredit(id uuid.UUID, amount float64) (*models.CreditCard, error) {
	panic("unexpected call to UpdateUsedCredit")
}

func (m *mockCreditCardQueryRepo) Delete(id uuid.UUID) error {
	panic("unexpected call to Delete")
}

func TestCashFlowQueryAdapter_GetMonthlySummary(t *testing.T) {
	svc := &mockCashFlowQueryService{
		monthlySummary: &models.MonthlyCashFlowSummary{
			Year:         2026,
			Month:        4,
			TotalIncome:  50000,
			TotalExpense: 12345,
			NetCashFlow:  37655,
			IncomeCount:  2,
			ExpenseCount: 7,
			ExpenseCategoryBreakdown: []*models.CategorySummary{
				{CategoryName: "交通", Amount: 300, Count: 2},
				{CategoryName: "飲食", Amount: 1500, Count: 4},
				{CategoryName: "房租", Amount: 9000, Count: 1},
				{CategoryName: "娛樂", Amount: 800, Count: 3},
				{CategoryName: "醫療", Amount: 400, Count: 1},
				{CategoryName: "雜項", Amount: 200, Count: 2},
			},
			ComparisonToPrev: &models.MonthComparison{
				PreviousMonth:    3,
				PreviousYear:     2026,
				ExpenseChange:    345,
				ExpenseChangePct: 2.88,
				IncomeChange:     1000,
				IncomeChangePct:  2.04,
			},
		},
	}

	adapter := NewCashFlowQueryAdapter(svc)

	result, err := adapter.GetMonthlySummary(2026, 4)

	require.NoError(t, err)
	require.Equal(t, 2026, svc.calledYear)
	require.Equal(t, 4, svc.calledMonth)
	require.Equal(t, 2026, result.Year)
	require.Equal(t, 4, result.Month)
	require.Equal(t, 50000.0, result.TotalIncome)
	require.Equal(t, 12345.0, result.TotalExpense)
	require.Equal(t, 37655.0, result.NetCashFlow)
	require.Equal(t, 2, result.IncomeCount)
	require.Equal(t, 7, result.ExpenseCount)
	require.Len(t, result.TopCategories, 5)
	require.Equal(t, []CategoryBreakdown{
		{Name: "房租", Amount: 9000, Count: 1},
		{Name: "飲食", Amount: 1500, Count: 4},
		{Name: "娛樂", Amount: 800, Count: 3},
		{Name: "醫療", Amount: 400, Count: 1},
		{Name: "交通", Amount: 300, Count: 2},
	}, result.TopCategories)
	require.NotNil(t, result.Comparison)
	require.Equal(t, &MonthComparisonResult{
		PreviousMonth:    3,
		PreviousYear:     2026,
		ExpenseChange:    345,
		ExpenseChangePct: 2.88,
		IncomeChange:     1000,
		IncomeChangePct:  2.04,
	}, result.Comparison)
}

func TestCashFlowQueryAdapter_EmptyResult(t *testing.T) {
	svc := &mockCashFlowQueryService{monthlySummary: &models.MonthlyCashFlowSummary{}}
	adapter := NewCashFlowQueryAdapter(svc)

	result, err := adapter.GetMonthlySummary(2026, 4)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 0, result.Year)
	require.Equal(t, 0, result.Month)
	require.Equal(t, 0.0, result.TotalIncome)
	require.Equal(t, 0.0, result.TotalExpense)
	require.Equal(t, 0.0, result.NetCashFlow)
	require.Equal(t, 0, result.IncomeCount)
	require.Equal(t, 0, result.ExpenseCount)
	require.Empty(t, result.TopCategories)
	require.Nil(t, result.Comparison)
}

func TestAccountBalanceQueryAdapter_GetAllBalances(t *testing.T) {
	bankRepo := &mockBankAccountQueryRepo{accounts: []*models.BankAccount{
		{BankName: "First Bank", AccountNumberLast4: "1234", Currency: models.CurrencyTWD, Balance: 10000},
		{BankName: "Mega Bank", AccountNumberLast4: "5678", Currency: models.CurrencyUSD, Balance: 2500.5},
	}}
	cardRepo := &mockCreditCardQueryRepo{cards: []*models.CreditCard{
		{IssuingBank: "CTBC", CardName: "Line Pay", CardNumberLast4: "9999", CreditLimit: 100000, UsedCredit: 25000},
	}}
	adapter := NewAccountBalanceQueryAdapter(bankRepo, cardRepo)

	result, err := adapter.GetAllBalances()

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, []BankAccountBalance{
		{Name: "First Bank", Last4: "1234", Currency: "TWD", Balance: 10000},
		{Name: "Mega Bank", Last4: "5678", Currency: "USD", Balance: 2500.5},
	}, result.BankAccounts)
	require.Equal(t, []CreditCardBalance{
		{Name: "CTBC Line Pay", Last4: "9999", CreditLimit: 100000, UsedCredit: 25000, Remaining: 75000, UsagePct: 25},
	}, result.CreditCards)
	require.NoError(t, result.BankError)
	require.NoError(t, result.CCError)
}

func TestAccountBalanceQueryAdapter_PartialFailure(t *testing.T) {
	bankRepo := &mockBankAccountQueryRepo{accounts: []*models.BankAccount{
		{BankName: "First Bank", AccountNumberLast4: "1234", Currency: models.CurrencyTWD, Balance: 10000},
	}}
	cardRepo := &mockCreditCardQueryRepo{err: errors.New("credit cards unavailable")}
	adapter := NewAccountBalanceQueryAdapter(bankRepo, cardRepo)

	result, err := adapter.GetAllBalances()

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, []BankAccountBalance{{Name: "First Bank", Last4: "1234", Currency: "TWD", Balance: 10000}}, result.BankAccounts)
	require.Empty(t, result.CreditCards)
	require.NoError(t, result.BankError)
	require.ErrorContains(t, result.CCError, "credit cards unavailable")
}

func TestAccountBalanceQueryAdapter_NoAccounts(t *testing.T) {
	adapter := NewAccountBalanceQueryAdapter(&mockBankAccountQueryRepo{}, &mockCreditCardQueryRepo{})

	result, err := adapter.GetAllBalances()

	require.NoError(t, err)
	require.NotNil(t, result)
	require.NotNil(t, result.BankAccounts)
	require.NotNil(t, result.CreditCards)
	require.Empty(t, result.BankAccounts)
	require.Empty(t, result.CreditCards)
}
