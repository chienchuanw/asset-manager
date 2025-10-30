package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// BankAccountService 銀行帳戶業務邏輯介面
type BankAccountService interface {
	CreateBankAccount(input *models.CreateBankAccountInput) (*models.BankAccount, error)
	GetBankAccount(id uuid.UUID) (*models.BankAccount, error)
	ListBankAccounts(currency *models.Currency) ([]*models.BankAccount, error)
	UpdateBankAccount(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error)
	DeleteBankAccount(id uuid.UUID) error
}

// bankAccountService 銀行帳戶業務邏輯實作
type bankAccountService struct {
	repo repository.BankAccountRepository
}

// NewBankAccountService 建立新的銀行帳戶 service
func NewBankAccountService(repo repository.BankAccountRepository) BankAccountService {
	return &bankAccountService{
		repo: repo,
	}
}

// CreateBankAccount 建立新的銀行帳戶
func (s *bankAccountService) CreateBankAccount(input *models.CreateBankAccountInput) (*models.BankAccount, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 建立銀行帳戶
	account, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	return account, nil
}

// GetBankAccount 取得銀行帳戶
func (s *bankAccountService) GetBankAccount(id uuid.UUID) (*models.BankAccount, error) {
	account, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bank account: %w", err)
	}

	return account, nil
}

// ListBankAccounts 列出所有銀行帳戶（可選擇性篩選幣別）
func (s *bankAccountService) ListBankAccounts(currency *models.Currency) ([]*models.BankAccount, error) {
	accounts, err := s.repo.GetAll(currency)
	if err != nil {
		return nil, fmt.Errorf("failed to list bank accounts: %w", err)
	}

	return accounts, nil
}

// UpdateBankAccount 更新銀行帳戶
func (s *bankAccountService) UpdateBankAccount(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 更新銀行帳戶
	account, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update bank account: %w", err)
	}

	return account, nil
}

// DeleteBankAccount 刪除銀行帳戶
func (s *bankAccountService) DeleteBankAccount(id uuid.UUID) error {
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete bank account: %w", err)
	}

	return nil
}

