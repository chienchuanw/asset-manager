package discord

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/google/uuid"
)

// CashFlowServiceAdapter bridges service.CashFlowService to CashFlowCreator.
type CashFlowServiceAdapter struct {
	svc service.CashFlowService
}

// NewCashFlowServiceAdapter wraps a CashFlowService for bot usage.
func NewCashFlowServiceAdapter(svc service.CashFlowService) *CashFlowServiceAdapter {
	return &CashFlowServiceAdapter{svc: svc}
}

func (a *CashFlowServiceAdapter) CreateCashFlowFromBot(input *BotCashFlowInput) (string, error) {
	date, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		date = time.Now()
	}

	categoryID, err := uuid.Parse(input.CategoryID)
	if err != nil {
		return "", err
	}

	sourceType := mapSourceType(input.SourceType)
	createInput := &models.CreateCashFlowInput{
		Date:        date,
		Type:        models.CashFlowType(input.Type),
		CategoryID:  categoryID,
		Amount:      input.Amount,
		Description: input.Description,
		SourceType:  &sourceType,
	}

	if input.SourceID != "" {
		sourceID, err := uuid.Parse(input.SourceID)
		if err == nil {
			createInput.SourceID = &sourceID
		}
	}

	record, err := a.svc.CreateCashFlow(createInput)
	if err != nil {
		return "", err
	}

	return record.ID.String(), nil
}

// CategoryRepoAdapter bridges repository.CategoryRepository to CategoryLoader.
type CategoryRepoAdapter struct {
	repo repository.CategoryRepository
}

// NewCategoryRepoAdapter wraps a CategoryRepository for bot usage.
func NewCategoryRepoAdapter(repo repository.CategoryRepository) *CategoryRepoAdapter {
	return &CategoryRepoAdapter{repo: repo}
}

func mapSourceType(sourceType string) models.SourceType {
	switch sourceType {
	case "bank_account":
		return models.SourceTypeBankAccount
	case "credit_card":
		return models.SourceTypeCreditCard
	case "cash":
		return models.SourceTypeCash
	default:
		return models.SourceTypeCash
	}
}

// AccountRepoAdapter bridges bank account and credit card repositories to AccountLoader.
type AccountRepoAdapter struct {
	bankRepo repository.BankAccountRepository
	cardRepo repository.CreditCardRepository
}

func NewAccountRepoAdapter(bankRepo repository.BankAccountRepository, cardRepo repository.CreditCardRepository) *AccountRepoAdapter {
	return &AccountRepoAdapter{bankRepo: bankRepo, cardRepo: cardRepo}
}

func (a *AccountRepoAdapter) LoadAccounts(sourceType string) ([]AccountInfo, error) {
	switch sourceType {
	case "bank_account":
		return a.loadBankAccounts()
	case "credit_card":
		return a.loadCreditCards()
	default:
		return nil, nil
	}
}

func (a *AccountRepoAdapter) loadBankAccounts() ([]AccountInfo, error) {
	accounts, err := a.bankRepo.GetAll(nil)
	if err != nil {
		return nil, err
	}
	result := make([]AccountInfo, 0, len(accounts))
	for _, acct := range accounts {
		label := fmt.Sprintf("%s *%s", acct.BankName, acct.AccountNumberLast4)
		result = append(result, AccountInfo{
			ID:   acct.ID.String(),
			Name: label,
			Type: "bank_account",
		})
	}
	return result, nil
}

func (a *AccountRepoAdapter) loadCreditCards() ([]AccountInfo, error) {
	cards, err := a.cardRepo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]AccountInfo, 0, len(cards))
	for _, card := range cards {
		label := fmt.Sprintf("%s %s *%s", card.IssuingBank, card.CardName, card.CardNumberLast4)
		result = append(result, AccountInfo{
			ID:   card.ID.String(),
			Name: label,
			Type: "credit_card",
		})
	}
	return result, nil
}

func (a *CategoryRepoAdapter) LoadCategories() ([]CategoryInfo, error) {
	categories, err := a.repo.GetAll(nil)
	if err != nil {
		return nil, err
	}

	result := make([]CategoryInfo, len(categories))
	for i, cat := range categories {
		result[i] = CategoryInfo{
			ID:   cat.ID.String(),
			Name: cat.Name,
			Type: string(cat.Type),
		}
	}
	return result, nil
}

// CashFlowQueryAdapter bridges service.CashFlowService to CashFlowQuerier.
type CashFlowQueryAdapter struct {
	svc service.CashFlowService
}

// NewCashFlowQueryAdapter wraps a CashFlowService for query usage.
func NewCashFlowQueryAdapter(svc service.CashFlowService) *CashFlowQueryAdapter {
	return &CashFlowQueryAdapter{svc: svc}
}

func (a *CashFlowQueryAdapter) GetMonthlySummary(year, month int) (*MonthlySummaryResult, error) {
	summary, err := a.svc.GetMonthlySummaryWithComparison(year, month)
	if err != nil {
		return nil, err
	}
	if summary == nil {
		return &MonthlySummaryResult{TopCategories: []CategoryBreakdown{}}, nil
	}

	result := &MonthlySummaryResult{
		Year:          summary.Year,
		Month:         summary.Month,
		TotalIncome:   summary.TotalIncome,
		TotalExpense:  summary.TotalExpense,
		NetCashFlow:   summary.NetCashFlow,
		IncomeCount:   summary.IncomeCount,
		ExpenseCount:  summary.ExpenseCount,
		TopCategories: mapTopCategories(summary.ExpenseCategoryBreakdown),
	}
	if summary.ComparisonToPrev != nil {
		result.Comparison = &MonthComparisonResult{
			PreviousMonth:    summary.ComparisonToPrev.PreviousMonth,
			PreviousYear:     summary.ComparisonToPrev.PreviousYear,
			ExpenseChange:    summary.ComparisonToPrev.ExpenseChange,
			ExpenseChangePct: summary.ComparisonToPrev.ExpenseChangePct,
			IncomeChange:     summary.ComparisonToPrev.IncomeChange,
			IncomeChangePct:  summary.ComparisonToPrev.IncomeChangePct,
		}
	}
	return result, nil
}

// AccountBalanceQueryAdapter bridges repositories to AccountBalanceQuerier.
type AccountBalanceQueryAdapter struct {
	bankRepo repository.BankAccountRepository
	cardRepo repository.CreditCardRepository
}

// NewAccountBalanceQueryAdapter wraps repositories for balance queries.
func NewAccountBalanceQueryAdapter(bankRepo repository.BankAccountRepository, cardRepo repository.CreditCardRepository) *AccountBalanceQueryAdapter {
	return &AccountBalanceQueryAdapter{bankRepo: bankRepo, cardRepo: cardRepo}
}

func (a *AccountBalanceQueryAdapter) GetAllBalances() (*AccountBalancesResult, error) {
	result := &AccountBalancesResult{
		BankAccounts: []BankAccountBalance{},
		CreditCards:  []CreditCardBalance{},
	}

	bankAccounts, bankErr := a.bankRepo.GetAll(nil)
	if bankErr == nil {
		result.BankAccounts = mapBankAccountBalances(bankAccounts)
	} else {
		result.BankError = bankErr
	}

	creditCards, ccErr := a.cardRepo.GetAll()
	if ccErr == nil {
		result.CreditCards = mapCreditCardBalances(creditCards)
	} else {
		result.CCError = ccErr
	}

	if bankErr != nil && ccErr != nil {
		return result, errors.Join(bankErr, ccErr)
	}

	return result, nil
}

func mapTopCategories(categories []*models.CategorySummary) []CategoryBreakdown {
	if len(categories) == 0 {
		return []CategoryBreakdown{}
	}

	breakdowns := make([]CategoryBreakdown, 0, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		breakdowns = append(breakdowns, CategoryBreakdown{
			Name:   category.CategoryName,
			Amount: category.Amount,
			Count:  category.Count,
		})
	}

	sort.Slice(breakdowns, func(i, j int) bool {
		return breakdowns[i].Amount > breakdowns[j].Amount
	})
	if len(breakdowns) > 5 {
		breakdowns = breakdowns[:5]
	}
	return breakdowns
}

func mapBankAccountBalances(accounts []*models.BankAccount) []BankAccountBalance {
	if len(accounts) == 0 {
		return []BankAccountBalance{}
	}

	result := make([]BankAccountBalance, 0, len(accounts))
	for _, account := range accounts {
		if account == nil {
			continue
		}
		result = append(result, BankAccountBalance{
			Name:     account.BankName,
			Last4:    account.AccountNumberLast4,
			Currency: string(account.Currency),
			Balance:  account.Balance,
		})
	}
	return result
}

func mapCreditCardBalances(cards []*models.CreditCard) []CreditCardBalance {
	if len(cards) == 0 {
		return []CreditCardBalance{}
	}

	result := make([]CreditCardBalance, 0, len(cards))
	for _, card := range cards {
		if card == nil {
			continue
		}
		result = append(result, CreditCardBalance{
			Name:        fmt.Sprintf("%s %s", card.IssuingBank, card.CardName),
			Last4:       card.CardNumberLast4,
			CreditLimit: card.CreditLimit,
			UsedCredit:  card.UsedCredit,
			Remaining:   card.AvailableCredit(),
			UsagePct:    card.CreditUtilization(),
		})
	}
	return result
}
