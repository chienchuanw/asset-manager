package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定義 JWT payload 的結構
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
// 參數:
//   - username: 使用者名稱
// 返回:
//   - string: JWT token 字串
//   - error: 錯誤訊息
func GenerateToken(username string) (string, error) {
	// 從環境變數取得 JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET environment variable is not set")
	}

	// 設定 token 有效期限為 24 小時
	expirationTime := time.Now().Add(24 * time.Hour)

	// 建立 claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 使用 HS256 演算法建立 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 簽署 token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 驗證 JWT token 並返回 claims
// 參數:
//   - tokenString: JWT token 字串
// 返回:
//   - *Claims: 解析後的 claims
//   - error: 錯誤訊息
func ValidateToken(tokenString string) (*Claims, error) {
	// 從環境變數取得 JWT secret
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}

	// 解析 token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽署方法是否為 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 驗證 token 是否有效
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

