package service

import (
	"errors"
	"os"

	"github.com/chienchuanw/asset-manager/internal/auth"
)

// AuthService 處理身份驗證相關的業務邏輯
type AuthService struct {
	username string
	password string
}

// NewAuthService 建立新的 AuthService 實例
// 從環境變數讀取 AUTH_USERNAME 和 AUTH_PASSWORD
func NewAuthService() *AuthService {
	return &AuthService{
		username: os.Getenv("AUTH_USERNAME"),
		password: os.Getenv("AUTH_PASSWORD"),
	}
}

// Login 驗證使用者帳號密碼並返回 JWT token
// 參數:
//   - username: 使用者名稱
//   - password: 密碼
// 返回:
//   - string: JWT token (登入成功時)
//   - error: 錯誤訊息 (登入失敗時)
func (s *AuthService) Login(username, password string) (string, error) {
	// 驗證輸入不為空
	if username == "" || password == "" {
		return "", errors.New("username and password are required")
	}

	// 驗證環境變數是否設定
	if s.username == "" || s.password == "" {
		return "", errors.New("authentication configuration is missing")
	}

	// 驗證帳號密碼
	if username != s.username || password != s.password {
		return "", errors.New("invalid username or password")
	}

	// 生成 JWT token
	token, err := auth.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil
}

