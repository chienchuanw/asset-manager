package service

import (
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCashFlowRepository 現金流 repository 的 mock
type MockCashFlowRepository struct {
	mock.Mock
}

func (m *MockCashFlowRepository) Create(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) GetByID(id uuid.UUID) (*models.CashFlow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) GetAll(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) Update(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashFlowRepository) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CashFlowSummary), args.Error(1)
}

// MockInstallmentRepository 分期 repository 的 mock
type MockInstallmentRepository struct {
	mock.Mock
}

func (m *MockInstallmentRepository) Create(input *models.CreateInstallmentInput) (*models.Installment, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Installment), args.Error(1)
}

func (m *MockInstallmentRepository) GetByID(id uuid.UUID) (*models.Installment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Installment), args.Error(1)
}

func (m *MockInstallmentRepository) List(filters repository.InstallmentFilters) ([]*models.Installment, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Installment), args.Error(1)
}

func (m *MockInstallmentRepository) Update(id uuid.UUID, input *models.UpdateInstallmentInput) (*models.Installment, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Installment), args.Error(1)
}

func (m *MockInstallmentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockInstallmentRepository) GetDueBillings(date time.Time) ([]*models.Installment, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Installment), args.Error(1)
}

func (m *MockInstallmentRepository) GetCompletingSoon(remainingCount int) ([]*models.Installment, error) {
	args := m.Called(remainingCount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Installment), args.Error(1)
}

// MockCategoryRepository 分類 repository 的 mock
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) GetByID(id uuid.UUID) (*models.CashFlowCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) Create(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) GetAll(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	args := m.Called(flowType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) Update(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

