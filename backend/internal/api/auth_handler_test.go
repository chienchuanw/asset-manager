package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/middleware"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)
}

// setupAuthTestRouter 設定測試用的 router
func setupAuthTestRouter(authHandler *AuthHandler) *gin.Engine {
	router := gin.New()

	// 不需要驗證的路由
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// 需要驗證的路由
	protectedGroup := router.Group("/api/auth")
	protectedGroup.Use(middleware.AuthMiddleware())
	{
		protectedGroup.GET("/me", authHandler.GetCurrentUser)
	}

	return router
}

// TestAuthHandler_Login_Success 測試成功登入
func TestAuthHandler_Login_Success(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 handler
	authService := service.NewAuthService()
	authHandler := NewAuthHandler(authService)
	router := setupAuthTestRouter(authHandler)

	// 建立請求 body
	loginReq := LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	body, _ := json.Marshal(loginReq)

	// 建立測試請求
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code, "應該返回 200 OK")

	// 驗證 response body
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Nil(t, response.Error, "不應該有錯誤")
	assert.NotNil(t, response.Data, "應該有 data")

	// 驗證 cookie 是否設定
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies, "應該設定 cookie")

	var tokenCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenCookie = cookie
			break
		}
	}
	assert.NotNil(t, tokenCookie, "應該有 token cookie")
	assert.NotEmpty(t, tokenCookie.Value, "token cookie 不應該是空的")
	assert.True(t, tokenCookie.HttpOnly, "token cookie 應該是 HttpOnly")
	assert.Equal(t, "/", tokenCookie.Path, "cookie path 應該是 /")
}

// TestAuthHandler_Login_Failure 測試登入失敗
func TestAuthHandler_Login_Failure(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 handler
	authService := service.NewAuthService()
	authHandler := NewAuthHandler(authService)
	router := setupAuthTestRouter(authHandler)

	testCases := []struct {
		name           string
		username       string
		password       string
		expectedStatus int
	}{
		{
			name:           "錯誤的帳號",
			username:       "wronguser",
			password:       "admin123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "錯誤的密碼",
			username:       "admin",
			password:       "wrongpassword",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "空白帳號",
			username:       "",
			password:       "admin123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "空白密碼",
			username:       "admin",
			password:       "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 建立請求 body
			loginReq := LoginRequest{
				Username: tc.username,
				Password: tc.password,
			}
			body, _ := json.Marshal(loginReq)

			// 建立測試請求
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 執行請求
			router.ServeHTTP(w, req)

			// 驗證結果
			assert.Equal(t, tc.expectedStatus, w.Code, "應該返回正確的 HTTP 狀態碼")

			// 驗證 response body
			var response APIResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.NotNil(t, response.Error, "應該有錯誤")
		})
	}
}

// TestAuthHandler_Logout 測試登出
func TestAuthHandler_Logout(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer os.Unsetenv("JWT_SECRET")

	// 建立 handler
	authService := service.NewAuthService()
	authHandler := NewAuthHandler(authService)
	router := setupAuthTestRouter(authHandler)

	// 建立測試請求
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code, "應該返回 200 OK")

	// 驗證 cookie 是否被清除
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies, "應該設定 cookie")

	var tokenCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "token" {
			tokenCookie = cookie
			break
		}
	}
	assert.NotNil(t, tokenCookie, "應該有 token cookie")
	assert.Empty(t, tokenCookie.Value, "token cookie 應該是空的")
	assert.Equal(t, -1, tokenCookie.MaxAge, "MaxAge 應該是 -1 (刪除 cookie)")
}

// TestAuthHandler_GetCurrentUser 測試取得當前使用者
func TestAuthHandler_GetCurrentUser(t *testing.T) {
	// 設定測試環境變數
	os.Setenv("AUTH_USERNAME", "admin")
	os.Setenv("AUTH_PASSWORD", "admin123")
	os.Setenv("JWT_SECRET", "test-secret-key")
	defer func() {
		os.Unsetenv("AUTH_USERNAME")
		os.Unsetenv("AUTH_PASSWORD")
		os.Unsetenv("JWT_SECRET")
	}()

	// 建立 handler
	authService := service.NewAuthService()
	authHandler := NewAuthHandler(authService)
	router := setupAuthTestRouter(authHandler)

	// 先登入取得 token
	loginReq := LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	body, _ := json.Marshal(loginReq)
	loginReqHTTP := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(body))
	loginReqHTTP.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReqHTTP)

	// 取得 token cookie
	var tokenCookie *http.Cookie
	for _, cookie := range loginW.Result().Cookies() {
		if cookie.Name == "token" {
			tokenCookie = cookie
			break
		}
	}
	require.NotNil(t, tokenCookie, "應該有 token cookie")

	// 使用 token 呼叫 /me
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(tokenCookie)
	w := httptest.NewRecorder()

	// 執行請求
	router.ServeHTTP(w, req)

	// 驗證結果
	assert.Equal(t, http.StatusOK, w.Code, "應該返回 200 OK")

	// 驗證 response body
	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Nil(t, response.Error, "不應該有錯誤")
	assert.NotNil(t, response.Data, "應該有 data")

	// 驗證使用者資訊
	dataMap, ok := response.Data.(map[string]interface{})
	require.True(t, ok, "data 應該是 map")
	assert.Equal(t, "admin", dataMap["username"], "username 應該是 admin")
}

