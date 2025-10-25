package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 處理身份驗證相關的 HTTP 請求
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 建立新的 AuthHandler 實例
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// LoginRequest 登入請求的結構
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登入成功的回應結構
type LoginResponse struct {
	Message string `json:"message"`
}

// UserResponse 使用者資訊的回應結構
type UserResponse struct {
	Username string `json:"username"`
}

// Login 處理登入請求
// @Summary 使用者登入
// @Description 驗證使用者帳號密碼並返回 JWT token (存在 httpOnly cookie)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登入資訊"
// @Success 200 {object} APIResponse[LoginResponse]
// @Failure 400 {object} APIResponse[any]
// @Failure 401 {object} APIResponse[any]
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	// 綁定並驗證請求 body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// 呼叫 service 進行登入驗證
	token, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error: &APIError{
				Code:    "LOGIN_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 設定 httpOnly cookie
	c.SetCookie(
		"token",           // cookie name
		token,             // cookie value
		24*60*60,          // maxAge (24 hours in seconds)
		"/",               // path
		"",                // domain (empty = current domain)
		false,             // secure (set to true in production with HTTPS)
		true,              // httpOnly
	)

	// 返回成功訊息
	c.JSON(http.StatusOK, APIResponse{
		Data: LoginResponse{
			Message: "Login successful",
		},
	})
}

// Logout 處理登出請求
// @Summary 使用者登出
// @Description 清除 JWT token cookie
// @Tags auth
// @Produce json
// @Success 200 {object} APIResponse[LoginResponse]
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 清除 cookie (設定 MaxAge 為 -1)
	c.SetCookie(
		"token",           // cookie name
		"",                // cookie value (empty)
		-1,                // maxAge (-1 = delete cookie)
		"/",               // path
		"",                // domain
		false,             // secure
		true,              // httpOnly
	)

	// 返回成功訊息
	c.JSON(http.StatusOK, APIResponse{
		Data: LoginResponse{
			Message: "Logout successful",
		},
	})
}

// GetCurrentUser 取得當前登入使用者的資訊
// @Summary 取得當前使用者
// @Description 取得當前登入使用者的資訊 (需要驗證)
// @Tags auth
// @Produce json
// @Success 200 {object} APIResponse[UserResponse]
// @Failure 401 {object} APIResponse[any]
// @Router /api/auth/me [get]
// @Security BearerAuth
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	// 從 context 取得使用者名稱 (由 AuthMiddleware 設定)
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Error: &APIError{
				Code:    "UNAUTHORIZED",
				Message: "User not authenticated",
			},
		})
		return
	}

	// 返回使用者資訊
	c.JSON(http.StatusOK, APIResponse{
		Data: UserResponse{
			Username: username.(string),
		},
	})
}

