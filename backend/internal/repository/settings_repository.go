package repository

import (
	"database/sql"
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// SettingsRepository 設定資料存取介面
type SettingsRepository interface {
	// GetByKey 根據 key 取得設定
	GetByKey(key string) (*models.Setting, error)
	// GetAll 取得所有設定
	GetAll() ([]*models.Setting, error)
	// Update 更新設定
	Update(key string, input *models.UpdateSettingInput) (*models.Setting, error)
}

type settingsRepository struct {
	db *sql.DB
}

// NewSettingsRepository 建立設定資料存取實例
func NewSettingsRepository(db *sql.DB) SettingsRepository {
	return &settingsRepository{db: db}
}

// GetByKey 根據 key 取得設定
func (r *settingsRepository) GetByKey(key string) (*models.Setting, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM settings
		WHERE key = $1
	`

	var setting models.Setting
	err := r.db.QueryRow(query, key).Scan(
		&setting.ID,
		&setting.Key,
		&setting.Value,
		&setting.Description,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("setting not found: %s", key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get setting: %w", err)
	}

	return &setting, nil
}

// GetAll 取得所有設定
func (r *settingsRepository) GetAll() ([]*models.Setting, error) {
	query := `
		SELECT id, key, value, description, created_at, updated_at
		FROM settings
		ORDER BY key
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all settings: %w", err)
	}
	defer rows.Close()

	var settings []*models.Setting
	for rows.Next() {
		var setting models.Setting
		err := rows.Scan(
			&setting.ID,
			&setting.Key,
			&setting.Value,
			&setting.Description,
			&setting.CreatedAt,
			&setting.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, &setting)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}

// Update 更新設定
func (r *settingsRepository) Update(key string, input *models.UpdateSettingInput) (*models.Setting, error) {
	query := `
		UPDATE settings
		SET value = $1, updated_at = CURRENT_TIMESTAMP
		WHERE key = $2
		RETURNING id, key, value, description, created_at, updated_at
	`

	var setting models.Setting
	err := r.db.QueryRow(query, input.Value, key).Scan(
		&setting.ID,
		&setting.Key,
		&setting.Value,
		&setting.Description,
		&setting.CreatedAt,
		&setting.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("setting not found: %s", key)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	return &setting, nil
}

