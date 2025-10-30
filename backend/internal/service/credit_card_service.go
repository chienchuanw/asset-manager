package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// CreditCardService 信用卡業務邏輯介面
type CreditCardService interface {
	CreateCreditCard(input *models.CreateCreditCardInput) (*models.CreditCard, error)
	GetCreditCard(id uuid.UUID) (*models.CreditCard, error)
	ListCreditCards() ([]*models.CreditCard, error)
	GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error)
	GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error)
	UpdateCreditCard(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error)
	DeleteCreditCard(id uuid.UUID) error
}

// creditCardService 信用卡業務邏輯實作
type creditCardService struct {
	repo repository.CreditCardRepository
}

// NewCreditCardService 建立新的信用卡 service
func NewCreditCardService(repo repository.CreditCardRepository) CreditCardService {
	return &creditCardService{
		repo: repo,
	}
}

// CreateCreditCard 建立新的信用卡
func (s *creditCardService) CreateCreditCard(input *models.CreateCreditCardInput) (*models.CreditCard, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 建立信用卡
	card, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create credit card: %w", err)
	}

	return card, nil
}

// GetCreditCard 取得信用卡
func (s *creditCardService) GetCreditCard(id uuid.UUID) (*models.CreditCard, error) {
	card, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card: %w", err)
	}

	return card, nil
}

// ListCreditCards 列出所有信用卡
func (s *creditCardService) ListCreditCards() ([]*models.CreditCard, error) {
	cards, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list credit cards: %w", err)
	}

	return cards, nil
}

// GetUpcomingBilling 取得即將到來的帳單日信用卡（未來 N 天內）
func (s *creditCardService) GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error) {
	cards, err := s.repo.GetUpcomingBilling(daysAhead)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming billing: %w", err)
	}

	return cards, nil
}

// GetUpcomingPayment 取得即將到來的繳款截止日信用卡（未來 N 天內）
func (s *creditCardService) GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error) {
	cards, err := s.repo.GetUpcomingPayment(daysAhead)
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming payment: %w", err)
	}

	return cards, nil
}

// UpdateCreditCard 更新信用卡
func (s *creditCardService) UpdateCreditCard(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 更新信用卡
	card, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update credit card: %w", err)
	}

	return card, nil
}

// DeleteCreditCard 刪除信用卡
func (s *creditCardService) DeleteCreditCard(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete credit card: %w", err)
	}

	return nil
}

