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

type mockCCPaymentCashFlowService struct {
	createResult *models.CashFlow
	err          error
	createInput  *models.CreateCashFlowInput
}

func (m *mockCCPaymentCashFlowService) CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	m.createInput = input
	return m.createResult, m.err
}

func (m *mockCCPaymentCashFlowService) GetCashFlow(id uuid.UUID) (*models.CashFlow, error) {
	panic("unexpected call to GetCashFlow")
}

func (m *mockCCPaymentCashFlowService) ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	panic("unexpected call to ListCashFlows")
}

func (m *mockCCPaymentCashFlowService) UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	panic("unexpected call to UpdateCashFlow")
}

func (m *mockCCPaymentCashFlowService) DeleteCashFlow(id uuid.UUID) error {
	panic("unexpected call to DeleteCashFlow")
}

func (m *mockCCPaymentCashFlowService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	panic("unexpected call to GetSummary")
}

func (m *mockCCPaymentCashFlowService) GetMonthlySummaryWithComparison(year, month int) (*models.MonthlyCashFlowSummary, error) {
	panic("unexpected call to GetMonthlySummaryWithComparison")
}

func (m *mockCCPaymentCashFlowService) GetYearlySummaryWithComparison(year int) (*models.YearlyCashFlowSummary, error) {
	panic("unexpected call to GetYearlySummaryWithComparison")
}

type mockCCPaymentCreditCardRepo struct {
	card     *models.CreditCard
	err      error
	calledID uuid.UUID
}

func (m *mockCCPaymentCreditCardRepo) Create(input *models.CreateCreditCardInput) (*models.CreditCard, error) {
	panic("unexpected call to Create")
}

func (m *mockCCPaymentCreditCardRepo) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	m.calledID = id
	return m.card, m.err
}

func (m *mockCCPaymentCreditCardRepo) GetAll() ([]*models.CreditCard, error) {
	panic("unexpected call to GetAll")
}

func (m *mockCCPaymentCreditCardRepo) GetByBillingDay(day int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetByBillingDay")
}

func (m *mockCCPaymentCreditCardRepo) GetByPaymentDueDay(day int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetByPaymentDueDay")
}

func (m *mockCCPaymentCreditCardRepo) GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetUpcomingBilling")
}

func (m *mockCCPaymentCreditCardRepo) GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error) {
	panic("unexpected call to GetUpcomingPayment")
}

func (m *mockCCPaymentCreditCardRepo) Update(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error) {
	panic("unexpected call to Update")
}

func (m *mockCCPaymentCreditCardRepo) UpdateUsedCredit(id uuid.UUID, amount float64) (*models.CreditCard, error) {
	panic("unexpected call to UpdateUsedCredit")
}

func (m *mockCCPaymentCreditCardRepo) Delete(id uuid.UUID) error {
	panic("unexpected call to Delete")
}

type mockCCPaymentCategoryRepo struct {
	categories []*models.CashFlowCategory
	err        error
}

func (m *mockCCPaymentCategoryRepo) Create(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	panic("unexpected")
}
func (m *mockCCPaymentCategoryRepo) GetByID(id uuid.UUID) (*models.CashFlowCategory, error) {
	panic("unexpected")
}
func (m *mockCCPaymentCategoryRepo) GetAll(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	return m.categories, m.err
}
func (m *mockCCPaymentCategoryRepo) Update(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	panic("unexpected")
}
func (m *mockCCPaymentCategoryRepo) Delete(id uuid.UUID) error { panic("unexpected") }
func (m *mockCCPaymentCategoryRepo) IsInUse(id uuid.UUID) (bool, error) {
	panic("unexpected")
}
func (m *mockCCPaymentCategoryRepo) Reorder(input *models.ReorderCategoryInput) error {
	panic("unexpected")
}
func (m *mockCCPaymentCategoryRepo) GetMaxSortOrder(flowType models.CashFlowType) (int, error) {
	panic("unexpected")
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

func TestCCPaymentAdapter_CustomAmount(t *testing.T) {
	transferCatID := uuid.New()
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	recordID := uuid.New()
	svc := &mockCCPaymentCashFlowService{createResult: &models.CashFlow{ID: recordID}}
	catRepo := &mockCCPaymentCategoryRepo{categories: []*models.CashFlowCategory{{ID: transferCatID, Name: "移轉", Type: models.CashFlowTypeTransferOut}}}
	adapter := NewCreditCardPaymentAdapter(svc, &mockCCPaymentCreditCardRepo{}, catRepo)

	id, actualAmount, err := adapter.CreatePaymentFromBot(&BotCCPaymentInput{
		CreditCardID:  creditCardID.String(),
		BankAccountID: bankAccountID.String(),
		Amount:        15000,
		Date:          "2026-04-08",
		PaymentType:   "custom",
	})

	require.NoError(t, err)
	require.Equal(t, recordID.String(), id)
	require.Equal(t, 15000.0, actualAmount)
	require.NotNil(t, svc.createInput)
	require.Equal(t, models.CashFlowTypeTransferOut, svc.createInput.Type)
	require.Equal(t, transferCatID, svc.createInput.CategoryID)
	require.Equal(t, 15000.0, svc.createInput.Amount)
	require.Equal(t, time.Date(2026, 4, 8, 0, 0, 0, 0, time.UTC), svc.createInput.Date)
	require.NotNil(t, svc.createInput.SourceType)
	require.Equal(t, models.SourceTypeBankAccount, *svc.createInput.SourceType)
	require.NotNil(t, svc.createInput.SourceID)
	require.Equal(t, bankAccountID, *svc.createInput.SourceID)
	require.NotNil(t, svc.createInput.TargetType)
	require.Equal(t, models.SourceTypeCreditCard, *svc.createInput.TargetType)
	require.NotNil(t, svc.createInput.TargetID)
	require.Equal(t, creditCardID, *svc.createInput.TargetID)
}

func TestCCPaymentAdapter_FullPayment(t *testing.T) {
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	recordID := uuid.New()
	transferCatID := uuid.New()
	svc := &mockCCPaymentCashFlowService{createResult: &models.CashFlow{ID: recordID}}
	cardRepo := &mockCCPaymentCreditCardRepo{card: &models.CreditCard{ID: creditCardID, UsedCredit: 23500}}
	catRepo := &mockCCPaymentCategoryRepo{categories: []*models.CashFlowCategory{{ID: transferCatID, Type: models.CashFlowTypeTransferOut}}}
	adapter := NewCreditCardPaymentAdapter(svc, cardRepo, catRepo)

	id, actualAmount, err := adapter.CreatePaymentFromBot(&BotCCPaymentInput{
		CreditCardID:  creditCardID.String(),
		BankAccountID: bankAccountID.String(),
		Amount:        0,
		Date:          "2026-04-08",
		PaymentType:   "full",
	})

	require.NoError(t, err)
	require.Equal(t, recordID.String(), id)
	require.Equal(t, 23500.0, actualAmount)
	require.Equal(t, creditCardID, cardRepo.calledID)
	require.NotNil(t, svc.createInput)
	require.Equal(t, 23500.0, svc.createInput.Amount)
}

func TestCCPaymentAdapter_MinimumPayment(t *testing.T) {
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	recordID := uuid.New()
	transferCatID := uuid.New()
	svc := &mockCCPaymentCashFlowService{createResult: &models.CashFlow{ID: recordID}}
	catRepo := &mockCCPaymentCategoryRepo{categories: []*models.CashFlowCategory{{ID: transferCatID, Type: models.CashFlowTypeTransferOut}}}
	adapter := NewCreditCardPaymentAdapter(svc, &mockCCPaymentCreditCardRepo{}, catRepo)

	id, actualAmount, err := adapter.CreatePaymentFromBot(&BotCCPaymentInput{
		CreditCardID:  creditCardID.String(),
		BankAccountID: bankAccountID.String(),
		Amount:        3000,
		Date:          "2026-04-08",
		PaymentType:   "minimum",
	})

	require.NoError(t, err)
	require.Equal(t, recordID.String(), id)
	require.Equal(t, 3000.0, actualAmount)
	require.NotNil(t, svc.createInput)
	require.Equal(t, 3000.0, svc.createInput.Amount)
}

func TestCCPaymentAdapter_Failure(t *testing.T) {
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	transferCatID := uuid.New()
	svc := &mockCCPaymentCashFlowService{err: errors.New("create failed")}
	catRepo := &mockCCPaymentCategoryRepo{categories: []*models.CashFlowCategory{{ID: transferCatID, Type: models.CashFlowTypeTransferOut}}}
	adapter := NewCreditCardPaymentAdapter(svc, &mockCCPaymentCreditCardRepo{}, catRepo)

	id, actualAmount, err := adapter.CreatePaymentFromBot(&BotCCPaymentInput{
		CreditCardID:  creditCardID.String(),
		BankAccountID: bankAccountID.String(),
		Amount:        15000,
		Date:          "2026-04-08",
		PaymentType:   "custom",
	})

	require.ErrorContains(t, err, "create failed")
	require.Empty(t, id)
	require.Equal(t, 15000.0, actualAmount)
}

func TestCCPaymentAdapter_FullPayment_ZeroUsedCredit(t *testing.T) {
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	recordID := uuid.New()
	transferCatID := uuid.New()
	svc := &mockCCPaymentCashFlowService{createResult: &models.CashFlow{ID: recordID}}
	cardRepo := &mockCCPaymentCreditCardRepo{card: &models.CreditCard{ID: creditCardID, UsedCredit: 0}}
	catRepo := &mockCCPaymentCategoryRepo{categories: []*models.CashFlowCategory{{ID: transferCatID, Type: models.CashFlowTypeTransferOut}}}
	adapter := NewCreditCardPaymentAdapter(svc, cardRepo, catRepo)

	id, actualAmount, err := adapter.CreatePaymentFromBot(&BotCCPaymentInput{
		CreditCardID:  creditCardID.String(),
		BankAccountID: bankAccountID.String(),
		Amount:        0,
		Date:          "2026-04-08",
		PaymentType:   "full",
	})

	require.NoError(t, err)
	require.Equal(t, recordID.String(), id)
	require.Equal(t, 0.0, actualAmount)
	require.NotNil(t, svc.createInput)
	require.Equal(t, 0.0, svc.createInput.Amount)
}
