package discord

import (
	"fmt"
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
