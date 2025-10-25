package repository

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSettingsRepository_GetByKey 測試取得單一設定
func TestSettingsRepository_GetByKey(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// 取得預設設定
	setting, err := repo.GetByKey("discord_enabled")
	require.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, "discord_enabled", setting.Key)
	assert.Equal(t, "false", setting.Value)
}

// TestSettingsRepository_GetByKey_NotFound 測試取得不存在的設定
func TestSettingsRepository_GetByKey_NotFound(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// 取得不存在的設定
	setting, err := repo.GetByKey("non_existent_key")
	assert.Error(t, err)
	assert.Nil(t, setting)
}

// TestSettingsRepository_GetAll 測試取得所有設定
func TestSettingsRepository_GetAll(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// 取得所有設定
	settings, err := repo.GetAll()
	require.NoError(t, err)
	assert.NotNil(t, settings)
	assert.GreaterOrEqual(t, len(settings), 7) // 至少有 7 個預設設定
}

// TestSettingsRepository_Update 測試更新設定
func TestSettingsRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// 更新設定
	input := &models.UpdateSettingInput{
		Value: "https://discord.com/api/webhooks/test",
	}
	setting, err := repo.Update("discord_webhook_url", input)
	require.NoError(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, "discord_webhook_url", setting.Key)
	assert.Equal(t, "https://discord.com/api/webhooks/test", setting.Value)

	// 驗證更新成功
	updated, err := repo.GetByKey("discord_webhook_url")
	require.NoError(t, err)
	assert.Equal(t, "https://discord.com/api/webhooks/test", updated.Value)
}

// TestSettingsRepository_Update_NotFound 測試更新不存在的設定
func TestSettingsRepository_Update_NotFound(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// 更新不存在的設定
	input := &models.UpdateSettingInput{
		Value: "test_value",
	}
	setting, err := repo.Update("non_existent_key", input)
	assert.Error(t, err)
	assert.Nil(t, setting)
}

