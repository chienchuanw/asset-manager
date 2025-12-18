package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// SubscriptionService 訂閱業務邏輯介面
type SubscriptionService interface {
	CreateSubscription(input *models.CreateSubscriptionInput) (*models.Subscription, error)
	GetSubscription(id uuid.UUID) (*models.Subscription, error)
	ListSubscriptions(filters repository.SubscriptionFilters) ([]*models.Subscription, error)
	UpdateSubscription(id uuid.UUID, input *models.UpdateSubscriptionInput) (*models.Subscription, error)
	DeleteSubscription(id uuid.UUID) error
	CancelSubscription(id uuid.UUID, endDate time.Time) (*models.Subscription, error)
	GetDueBillings(date time.Time) ([]*models.Subscription, error)
	GetExpiringSoon(days int) ([]*models.Subscription, error)
}

// subscriptionService 訂閱業務邏輯實作
type subscriptionService struct {
	repo         repository.SubscriptionRepository
	categoryRepo repository.CategoryRepository
}

// NewSubscriptionService 建立新的訂閱 service
func NewSubscriptionService(
	repo repository.SubscriptionRepository,
	categoryRepo repository.CategoryRepository,
) SubscriptionService {
	return &subscriptionService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// CreateSubscription 建立新的訂閱
func (s *subscriptionService) CreateSubscription(input *models.CreateSubscriptionInput) (*models.Subscription, error) {
	// 驗證訂閱名稱
	if input.Name == "" {
		return nil, fmt.Errorf("subscription name is required")
	}

	if len(input.Name) > 255 {
		return nil, fmt.Errorf("subscription name must not exceed 255 characters")
	}

	// 驗證金額
	if input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 驗證計費週期
	if !input.BillingCycle.Validate() {
		return nil, fmt.Errorf("invalid billing cycle: %s", input.BillingCycle)
	}

	// 驗證扣款日
	if input.BillingDay < 1 || input.BillingDay > 31 {
		return nil, fmt.Errorf("billing day must be between 1 and 31")
	}

	// 驗證分類是否存在且為支出類型
	category, err := s.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	if category.Type != models.CashFlowTypeExpense {
		return nil, fmt.Errorf("subscription category must be expense type")
	}

	// 驗證日期
	if input.EndDate != nil && input.EndDate.Before(input.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// 驗證付款方式
	if !input.PaymentMethod.Validate() {
		return nil, fmt.Errorf("invalid payment method: %s", input.PaymentMethod)
	}

	// 驗證帳戶 ID（當付款方式需要帳戶時）
	if input.PaymentMethod.RequiresAccountID() && input.AccountID == nil {
		return nil, fmt.Errorf("account ID is required for payment method: %s", input.PaymentMethod)
	}

	// 呼叫 repository 建立訂閱
	subscription, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

// GetSubscription 取得單筆訂閱
func (s *subscriptionService) GetSubscription(id uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetByID(id)
}

// ListSubscriptions 取得訂閱列表
func (s *subscriptionService) ListSubscriptions(filters repository.SubscriptionFilters) ([]*models.Subscription, error) {
	// 驗證篩選條件
	if filters.Status != nil && !filters.Status.Validate() {
		return nil, fmt.Errorf("invalid subscription status filter: %s", *filters.Status)
	}

	return s.repo.List(filters)
}

// UpdateSubscription 更新訂閱
func (s *subscriptionService) UpdateSubscription(id uuid.UUID, input *models.UpdateSubscriptionInput) (*models.Subscription, error) {
	// 驗證訂閱是否存在
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// 驗證訂閱名稱
	if input.Name != nil {
		if *input.Name == "" {
			return nil, fmt.Errorf("subscription name cannot be empty")
		}
		if len(*input.Name) > 255 {
			return nil, fmt.Errorf("subscription name must not exceed 255 characters")
		}
	}

	// 驗證金額
	if input.Amount != nil && *input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 驗證計費週期
	if input.BillingCycle != nil && !input.BillingCycle.Validate() {
		return nil, fmt.Errorf("invalid billing cycle: %s", *input.BillingCycle)
	}

	// 驗證扣款日
	if input.BillingDay != nil {
		if *input.BillingDay < 1 || *input.BillingDay > 31 {
			return nil, fmt.Errorf("billing day must be between 1 and 31")
		}
	}

	// 驗證分類
	if input.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
		if category.Type != models.CashFlowTypeExpense {
			return nil, fmt.Errorf("subscription category must be expense type")
		}
	}

	// 驗證日期
	if input.EndDate != nil && input.EndDate.Before(existing.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// 驗證付款方式
	if input.PaymentMethod != nil {
		if !input.PaymentMethod.Validate() {
			return nil, fmt.Errorf("invalid payment method: %s", *input.PaymentMethod)
		}

		// 驗證帳戶 ID（當付款方式需要帳戶時）
		if input.PaymentMethod.RequiresAccountID() && input.AccountID == nil {
			return nil, fmt.Errorf("account ID is required for payment method: %s", *input.PaymentMethod)
		}
	}

	// 呼叫 repository 更新訂閱
	subscription, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return subscription, nil
}

// DeleteSubscription 刪除訂閱
func (s *subscriptionService) DeleteSubscription(id uuid.UUID) error {
	// 驗證訂閱是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// 呼叫 repository 刪除訂閱
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// CancelSubscription 取消訂閱
func (s *subscriptionService) CancelSubscription(id uuid.UUID, endDate time.Time) (*models.Subscription, error) {
	// 驗證訂閱是否存在
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// 檢查訂閱是否已經取消
	if existing.Status == models.SubscriptionStatusCancelled {
		return nil, fmt.Errorf("subscription is already cancelled")
	}

	// 驗證結束日期
	if endDate.Before(existing.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// 準備更新資料
	status := models.SubscriptionStatusCancelled
	updateInput := &models.UpdateSubscriptionInput{
		Status:  &status,
		EndDate: &endDate,
	}

	// 呼叫 repository 更新訂閱
	subscription, err := s.repo.Update(id, updateInput)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel subscription: %w", err)
	}

	return subscription, nil
}

// GetDueBillings 取得指定日期需要扣款的訂閱
func (s *subscriptionService) GetDueBillings(date time.Time) ([]*models.Subscription, error) {
	return s.repo.GetDueBillings(date)
}

// GetExpiringSoon 取得即將到期的訂閱
func (s *subscriptionService) GetExpiringSoon(days int) ([]*models.Subscription, error) {
	if days <= 0 {
		return nil, fmt.Errorf("days must be greater than zero")
	}

	return s.repo.GetExpiringSoon(days)
}

