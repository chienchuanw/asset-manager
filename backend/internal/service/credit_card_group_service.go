package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// CreditCardGroupService 信用卡群組業務邏輯介面
type CreditCardGroupService interface {
	CreateCreditCardGroup(input *models.CreateCreditCardGroupInput) (*models.CreditCardGroupWithCards, error)
	GetCreditCardGroup(id uuid.UUID) (*models.CreditCardGroupWithCards, error)
	ListCreditCardGroups() ([]*models.CreditCardGroupWithCards, error)
	UpdateCreditCardGroup(id uuid.UUID, input *models.UpdateCreditCardGroupInput) (*models.CreditCardGroup, error)
	DeleteCreditCardGroup(id uuid.UUID) error
	AddCardsToGroup(groupID uuid.UUID, input *models.AddCardsToGroupInput) error
	RemoveCardsFromGroup(groupID uuid.UUID, input *models.RemoveCardsFromGroupInput) error
}

// creditCardGroupService 信用卡群組業務邏輯實作
type creditCardGroupService struct {
	groupRepo repository.CreditCardGroupRepository
	cardRepo  repository.CreditCardRepository
}

// NewCreditCardGroupService 建立新的信用卡群組 service
func NewCreditCardGroupService(
	groupRepo repository.CreditCardGroupRepository,
	cardRepo repository.CreditCardRepository,
) CreditCardGroupService {
	return &creditCardGroupService{
		groupRepo: groupRepo,
		cardRepo:  cardRepo,
	}
}

// CreateCreditCardGroup 建立新的信用卡群組
func (s *creditCardGroupService) CreateCreditCardGroup(input *models.CreateCreditCardGroupInput) (*models.CreditCardGroupWithCards, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 將字串 ID 轉換為 UUID
	cardIDs := make([]uuid.UUID, len(input.CardIDs))
	for i, idStr := range input.CardIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid card ID: %w", err)
		}
		cardIDs[i] = id
	}

	// 驗證所有卡片存在且符合條件
	if _, err := s.validateCardsForGroup(cardIDs, input.IssuingBank, input.SharedCreditLimit); err != nil {
		return nil, err
	}

	// 建立群組
	group, err := s.groupRepo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create credit card group: %w", err)
	}

	// 將卡片加入群組
	if err := s.groupRepo.AddCardsToGroup(group.ID, cardIDs); err != nil {
		// 如果加入失敗,嘗試刪除已建立的群組
		_ = s.groupRepo.Delete(group.ID)
		return nil, fmt.Errorf("failed to add cards to group: %w", err)
	}

	// 取得完整的群組資料(包含卡片)
	groupWithCards, err := s.groupRepo.GetByID(group.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created group: %w", err)
	}

	return groupWithCards, nil
}

// GetCreditCardGroup 取得信用卡群組
func (s *creditCardGroupService) GetCreditCardGroup(id uuid.UUID) (*models.CreditCardGroupWithCards, error) {
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get credit card group: %w", err)
	}

	return group, nil
}

// ListCreditCardGroups 取得所有信用卡群組
func (s *creditCardGroupService) ListCreditCardGroups() ([]*models.CreditCardGroupWithCards, error) {
	groups, err := s.groupRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list credit card groups: %w", err)
	}

	return groups, nil
}

// UpdateCreditCardGroup 更新信用卡群組
func (s *creditCardGroupService) UpdateCreditCardGroup(id uuid.UUID, input *models.UpdateCreditCardGroupInput) (*models.CreditCardGroup, error) {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	// 如果更新共享額度,需要驗證群組內所有卡片的額度是否一致
	if input.SharedCreditLimit != nil {
		group, err := s.groupRepo.GetByID(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get group: %w", err)
		}

		// 驗證所有卡片的額度是否與新的共享額度一致
		for _, card := range group.Cards {
			if card.CreditLimit != *input.SharedCreditLimit {
				return nil, fmt.Errorf("card %s has different credit limit (%.2f), expected %.2f",
					card.CardName, card.CreditLimit, *input.SharedCreditLimit)
			}
		}
	}

	// 更新群組
	group, err := s.groupRepo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update credit card group: %w", err)
	}

	return group, nil
}

// DeleteCreditCardGroup 刪除信用卡群組
func (s *creditCardGroupService) DeleteCreditCardGroup(id uuid.UUID) error {
	// 取得群組資料
	group, err := s.groupRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// 將群組內的卡片移除群組關聯
	if len(group.Cards) > 0 {
		cardIDs := make([]uuid.UUID, len(group.Cards))
		for i, card := range group.Cards {
			cardIDs[i] = card.ID
		}
		if err := s.groupRepo.RemoveCardsFromGroup(cardIDs); err != nil {
			return fmt.Errorf("failed to remove cards from group: %w", err)
		}
	}

	// 刪除群組
	if err := s.groupRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete credit card group: %w", err)
	}

	return nil
}

// AddCardsToGroup 新增卡片到群組
func (s *creditCardGroupService) AddCardsToGroup(groupID uuid.UUID, input *models.AddCardsToGroupInput) error {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	// 取得群組資料
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// 將字串 ID 轉換為 UUID
	cardIDs := make([]uuid.UUID, len(input.CardIDs))
	for i, idStr := range input.CardIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid card ID: %w", err)
		}
		cardIDs[i] = id
	}

	// 驗證所有卡片存在且符合條件
	_, err = s.validateCardsForGroup(cardIDs, group.IssuingBank, group.SharedCreditLimit)
	if err != nil {
		return err
	}

	// 將卡片加入群組
	if err := s.groupRepo.AddCardsToGroup(groupID, cardIDs); err != nil {
		return fmt.Errorf("failed to add cards to group: %w", err)
	}

	return nil
}

// RemoveCardsFromGroup 從群組移除卡片
func (s *creditCardGroupService) RemoveCardsFromGroup(groupID uuid.UUID, input *models.RemoveCardsFromGroupInput) error {
	// 驗證輸入資料
	if err := input.Validate(); err != nil {
		return fmt.Errorf("invalid input: %w", err)
	}

	// 取得群組資料
	group, err := s.groupRepo.GetByID(groupID)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// 將字串 ID 轉換為 UUID
	cardIDs := make([]uuid.UUID, len(input.CardIDs))
	for i, idStr := range input.CardIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return fmt.Errorf("invalid card ID: %w", err)
		}
		cardIDs[i] = id
	}

	// 驗證要移除的卡片是否都在群組內
	for _, cardID := range cardIDs {
		found := false
		for _, card := range group.Cards {
			if card.ID == cardID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("card %s is not in the group", cardID)
		}
	}

	// 從群組移除卡片
	if err := s.groupRepo.RemoveCardsFromGroup(cardIDs); err != nil {
		return fmt.Errorf("failed to remove cards from group: %w", err)
	}

	return nil
}

// validateCardsForGroup 驗證卡片是否符合群組條件
func (s *creditCardGroupService) validateCardsForGroup(cardIDs []uuid.UUID, issuingBank string, sharedCreditLimit float64) ([]*models.CreditCard, error) {
	cards := make([]*models.CreditCard, 0, len(cardIDs))

	for _, cardID := range cardIDs {
		card, err := s.cardRepo.GetByID(cardID)
		if err != nil {
			return nil, fmt.Errorf("card %s not found: %w", cardID, err)
		}

		// 驗證卡片尚未屬於任何群組
		if card.GroupID != nil {
			return nil, fmt.Errorf("card %s already belongs to a group", card.CardName)
		}

		// 驗證發卡銀行一致
		if card.IssuingBank != issuingBank {
			return nil, fmt.Errorf("card %s has different issuing bank (%s), expected %s",
				card.CardName, card.IssuingBank, issuingBank)
		}

		// 驗證信用額度一致
		if card.CreditLimit != sharedCreditLimit {
			return nil, fmt.Errorf("card %s has different credit limit (%.2f), expected %.2f",
				card.CardName, card.CreditLimit, sharedCreditLimit)
		}

		cards = append(cards, card)
	}

	return cards, nil
}

