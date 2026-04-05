package discord

import (
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

	cashSourceType := models.SourceTypeCash
	createInput := &models.CreateCashFlowInput{
		Date:        date,
		Type:        models.CashFlowType(input.Type),
		CategoryID:  categoryID,
		Amount:      input.Amount,
		Description: input.Description,
		SourceType:  &cashSourceType,
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
