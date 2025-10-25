package middleware

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/auth"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware 驗證 JWT token 的 middleware
// 從 cookie 中讀取 token，驗證後將使用者資訊存入 context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 cookie 中取得 token
		token, err := c.Cookie("token")
		if err != nil {
			// 沒有 token，返回 401 Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Authentication required",
				},
			})
			c.Abort()
			return
		}

		// 驗證 token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			// token 無效，返回 401 Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"code":    "INVALID_TOKEN",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// 將使用者資訊存入 context
		c.Set("username", claims.Username)

		// 繼續處理請求
		c.Next()
	}
}

