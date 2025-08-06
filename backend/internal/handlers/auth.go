package handlers

import (
	"net/http"
	"sso/internal/services"
	"sso/internal/session"
	"time"

	"github.com/gin-gonic/gin"
)

type authHandler struct{}

// NewAuthHandler 創建新的認證處理器
func NewAuthHandler() (*authHandler, error) {
	return &authHandler{}, nil
}

// GetJWTMiddleware 返回 JWT 中間件
func (h *authHandler) GetJWTMiddleware() gin.HandlerFunc {
	return session.JWTMiddleware()
}

func (h *authHandler) Login(c *gin.Context) {
	reqData := struct {
		Email    string `form:"email" json:"email" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}{}
	if err := c.ShouldBind(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供 email 和密碼"})
		return
	}

	// 使用用戶服務進行認證
	user, err := services.UserService.Login(reqData.Email, reqData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 創建 JWT 用戶結構
	jwtUser := &session.JWTUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	// 生成 token 並設置 Cookie
	tokenString, err := session.GenerateTokenAndSetCookie(c, jwtUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT 生成失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   tokenString,
		"expire":  time.Now().Add(time.Hour * 24).Format(time.RFC3339),
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

func (h *authHandler) Logout(c *gin.Context) {
	// 清除 JWT Cookie
	session.ClearTokenCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// RefreshToken 刷新 JWT token
func (h *authHandler) RefreshToken(c *gin.Context) {
	// 獲取當前用戶
	user, exists := session.GetJWTUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未認證的用戶"})
		return
	}

	// 生成新的 token 並設置 Cookie
	tokenString, err := session.GenerateTokenAndSetCookie(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 刷新失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Token 刷新成功",
		"token":   tokenString,
		"expire":  time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	})
}

func (h *authHandler) Register(c *gin.Context) {
	reqData := struct {
		Email           string `form:"email" binding:"required"`
		Password        string `form:"password" binding:"required"`
		PasswordConfirm string `form:"password_confirm" binding:"required,eqfield=Password"`
		Name            string `form:"name" binding:"required"`
		AvatarBase64    string `form:"avatar_base64"`
	}{}
	if err := c.ShouldBind(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "資料不完整，或密碼不一致"})
		return
	}

	_, err := services.UserService.Register(
		reqData.Email, reqData.Password,
		reqData.Name, reqData.AvatarBase64,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "register success"})
}

func (h *authHandler) GoogleLogin(c *gin.Context) {
	// TODO: 實際的 Google OAuth 流程
	// 1. 生成 state 參數防止 CSRF 攻擊
	// 2. 構建 Google OAuth 授權 URL
	// 3. 重定向用戶到 Google 授權頁面

	// 暫時返回應該重定向的 URL
	c.JSON(http.StatusOK, gin.H{
		"message":      "請重定向到 Google OAuth",
		"redirect_url": "https://accounts.google.com/oauth/authorize?...", // 實際應用中需要構建完整的 URL
		"note":         "實際實現需要配置 Google OAuth 客戶端",
	})
}

func (h *authHandler) GithubLogin(c *gin.Context) {
	// TODO: 實際的 GitHub OAuth 流程
	// 1. 生成 state 參數防止 CSRF 攻擊
	// 2. 構建 GitHub OAuth 授權 URL
	// 3. 重定向用戶到 GitHub 授權頁面

	// 暫時返回應該重定向的 URL
	c.JSON(http.StatusOK, gin.H{
		"message":      "請重定向到 GitHub OAuth",
		"redirect_url": "https://github.com/login/oauth/authorize?...", // 實際應用中需要構建完整的 URL
		"note":         "實際實現需要配置 GitHub OAuth 客戶端",
	})
}

func (h *authHandler) ExternalCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "google" && provider != "github" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不支援的第三方登入提供者"})
		return
	}

	// TODO: 實際的 OAuth 流程處理
	// 這裡應該：
	// 1. 驗證 OAuth 回調參數（code, state 等）
	// 2. 使用 code 換取 access token
	// 3. 使用 access token 獲取用戶信息

	// 暫時使用模擬數據示範如何為第三方登入用戶發放 JWT
	// 在實際應用中，這些數據應該從 OAuth 提供者獲取
	mockEmail := "user@example.com"
	mockName := "第三方用戶"
	mockProviderUserID := "12345"

	// 使用用戶服務處理第三方登入或註冊
	user, err := services.UserService.LoginOrRegisterWithOAuth(mockEmail, mockName, provider, mockProviderUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 創建 JWT 用戶結構
	jwtUser := &session.JWTUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	// 使用新的函數生成 token 並設置 Cookie
	tokenString, err := session.GenerateTokenAndSetCookie(c, jwtUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT 生成失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "第三方登入成功",
		"token":   tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// GetProfile 獲取當前已認證用戶的資料
func (h *authHandler) GetProfile(c *gin.Context) {
	// 從 JWT 中獲取用戶信息
	jwtUser, exists := session.GetJWTUser(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未認證的用戶"})
		return
	}

	// 返回 JWT 中的用戶信息
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    jwtUser.ID,
			"email": jwtUser.Email,
			"name":  jwtUser.Name,
		},
	})
}
