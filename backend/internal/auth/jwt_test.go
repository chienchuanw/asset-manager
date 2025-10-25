package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken 測試生成 JWT token
func TestGenerateToken(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	username := "testuser"

	// 生成 token
	tokenString, err := GenerateToken(username)

	// 驗證結果
	require.NoError(t, err, "生成 token 不應該發生錯誤")
	assert.NotEmpty(t, tokenString, "token 不應該是空字串")

	// 驗證 token 格式是否正確（應該包含三個部分，用 . 分隔）
	// JWT 格式: header.payload.signature
	assert.Contains(t, tokenString, ".", "token 應該包含 . 分隔符")
}

// TestValidateToken_ValidToken 測試驗證有效的 JWT token
func TestValidateToken_ValidToken(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	username := "testuser"

	// 先生成一個有效的 token
	tokenString, err := GenerateToken(username)
	require.NoError(t, err)

	// 驗證 token
	claims, err := ValidateToken(tokenString)

	// 驗證結果
	require.NoError(t, err, "驗證有效 token 不應該發生錯誤")
	assert.NotNil(t, claims, "claims 不應該是 nil")
	assert.Equal(t, username, claims.Username, "username 應該一致")
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "token 應該還沒過期")
}

// TestValidateToken_ExpiredToken 測試驗證過期的 JWT token
func TestValidateToken_ExpiredToken(t *testing.T) {
	// 設定測試用的 JWT secret
	secret := "test-secret-key-for-testing"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// 手動建立一個已過期的 token（過期時間設為 1 秒前）
	claims := &Claims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	require.NoError(t, err)

	// 驗證過期的 token
	_, err = ValidateToken(tokenString)

	// 驗證結果：應該要有錯誤
	assert.Error(t, err, "驗證過期 token 應該要有錯誤")
}

// TestValidateToken_InvalidToken 測試驗證無效的 JWT token
func TestValidateToken_InvalidToken(t *testing.T) {
	// 設定測試用的 JWT secret
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	testCases := []struct {
		name        string
		tokenString string
	}{
		{
			name:        "空字串",
			tokenString: "",
		},
		{
			name:        "隨機字串",
			tokenString: "this-is-not-a-valid-jwt-token",
		},
		{
			name:        "格式錯誤的 token",
			tokenString: "header.payload",
		},
		{
			name:        "使用錯誤 secret 簽署的 token",
			tokenString: generateTokenWithSecret("testuser", "wrong-secret"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 驗證無效的 token
			_, err := ValidateToken(tc.tokenString)

			// 驗證結果：應該要有錯誤
			assert.Error(t, err, "驗證無效 token 應該要有錯誤")
		})
	}
}

// TestValidateToken_MissingSecret 測試缺少 JWT_SECRET 環境變數
func TestValidateToken_MissingSecret(t *testing.T) {
	// 確保 JWT_SECRET 環境變數不存在
	os.Unsetenv("JWT_SECRET")

	// 嘗試驗證 token
	_, err := ValidateToken("any-token")

	// 驗證結果：應該要有錯誤
	assert.Error(t, err, "缺少 JWT_SECRET 應該要有錯誤")
	assert.Contains(t, err.Error(), "JWT_SECRET", "錯誤訊息應該提到 JWT_SECRET")
}

// generateTokenWithSecret 使用指定的 secret 生成 token（測試輔助函式）
func generateTokenWithSecret(username, secret string) string {
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

