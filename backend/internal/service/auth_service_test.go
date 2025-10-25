package service

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuthService_Login_Success 測試正確的帳號密碼登入成功
func TestAuthService_Login_Success(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入
	token, err := authService.Login("admin", "admin123")

	// 驗證結果
	require.NoError(t, err, "正確的帳號密碼應該登入成功")
	assert.NotEmpty(t, token, "應該返回 JWT token")
}

// TestAuthService_Login_WrongUsername 測試錯誤的帳號登入失敗
func TestAuthService_Login_WrongUsername(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入（錯誤的帳號）
	token, err := authService.Login("wronguser", "admin123")

	// 驗證結果
	assert.Error(t, err, "錯誤的帳號應該登入失敗")
	assert.Empty(t, token, "不應該返回 token")
	assert.Contains(t, err.Error(), "invalid", "錯誤訊息應該包含 invalid")
}

// TestAuthService_Login_WrongPassword 測試錯誤的密碼登入失敗
func TestAuthService_Login_WrongPassword(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入（錯誤的密碼）
	token, err := authService.Login("admin", "wrongpassword")

	// 驗證結果
	assert.Error(t, err, "錯誤的密碼應該登入失敗")
	assert.Empty(t, token, "不應該返回 token")
	assert.Contains(t, err.Error(), "invalid", "錯誤訊息應該包含 invalid")
}

// TestAuthService_Login_EmptyUsername 測試空白帳號登入失敗
func TestAuthService_Login_EmptyUsername(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入（空白帳號）
	token, err := authService.Login("", "admin123")

	// 驗證結果
	assert.Error(t, err, "空白帳號應該登入失敗")
	assert.Empty(t, token, "不應該返回 token")
}

// TestAuthService_Login_EmptyPassword 測試空白密碼登入失敗
func TestAuthService_Login_EmptyPassword(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入（空白密碼）
	token, err := authService.Login("admin", "")

	// 驗證結果
	assert.Error(t, err, "空白密碼應該登入失敗")
	assert.Empty(t, token, "不應該返回 token")
}

// TestAuthService_Login_MissingEnvVars 測試缺少環境變數
func TestAuthService_Login_MissingEnvVars(t *testing.T) {
	// 確保環境變數不存在
	os.Unsetenv("AUTH_USERNAME")
	os.Unsetenv("AUTH_PASSWORD")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 建立 AuthService
	authService := NewAuthService()

	// 執行登入
	token, err := authService.Login("admin", "admin123")

	// 驗證結果
	assert.Error(t, err, "缺少環境變數應該登入失敗")
	assert.Empty(t, token, "不應該返回 token")
}

