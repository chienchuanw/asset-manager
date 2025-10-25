package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)
}

// TestAuthMiddleware_ValidToken 測試有效 token 通過驗證
func TestAuthMiddleware_ValidToken(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 生成有效的 token
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err)

	// 建立測試 router
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		// 從 context 取得使用者名稱
		username, exists := c.Get("username")
		assert.True(t, exists, "username 應該存在於 context 中")
		assert.Equal(t, "testuser", username, "username 應該正確")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 建立測試請求
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code, "應該返回 200 OK")
}

// TestAuthMiddleware_MissingToken 測試缺少 token 被拒絕
func TestAuthMiddleware_MissingToken(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 建立測試 router
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 建立測試請求（不帶 token）
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusUnauthorized, w.Code, "應該返回 401 Unauthorized")
}

// TestAuthMiddleware_InvalidToken 測試無效 token 被拒絕
func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 建立測試 router
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	testCases := []struct {
		name  string
		token string
	}{
		{
			name:  "隨機字串",
			token: "invalid-token-string",
		},
		{
			name:  "空字串",
			token: "",
		},
		{
			name:  "格式錯誤",
			token: "header.payload",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 建立測試請求
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: tc.token,
			})
			w := httptest.NewRecorder()

			// 執行請求
			router.ServeHTTP(w, req)

			// 驗證結果
			assert.Equal(t, http.StatusUnauthorized, w.Code, "應該返回 401 Unauthorized")
		})
	}
}

// TestAuthMiddleware_ExpiredToken 測試過期 token 被拒絕
func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// 設定測試環境變數
	secret := "test-secret-key"
	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// 手動建立一個已過期的 token
	claims := &auth.Claims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// 建立測試 router
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 建立測試請求
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusUnauthorized, w.Code, "應該返回 401 Unauthorized")
}

// TestAuthMiddleware_WrongSecret 測試使用錯誤 secret 簽署的 token 被拒絕
func TestAuthMiddleware_WrongSecret(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("JWT_SECRET", "correct-secret")
	defer os.Unsetenv("JWT_SECRET")

	// 使用錯誤的 secret 生成 token
	wrongSecret := "wrong-secret"
	claims := &auth.Claims{
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(wrongSecret))
	assert.NoError(t, err)

	// 建立測試 router
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// 建立測試請求
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusUnauthorized, w.Code, "應該返回 401 Unauthorized")
}

