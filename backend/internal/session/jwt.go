package session

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	identityKey  = "user_id"
	jwtSecret    = "secret key"   // 在生產環境中應該從環境變數讀取
	tokenTimeout = time.Hour * 24 // 24 小時
	cookieName   = "jwt"
)

var (
	ErrMissingToken = errors.New("缺少 JWT token")
	ErrInvalidToken = errors.New("無效的 JWT token")
)

// JWT 用戶結構
type JWTUser struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// JWT Claims 結構
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(user *JWTUser) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTimeout)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "SSO Zone",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken 驗證 JWT token
func ValidateToken(tokenString string) (*JWTUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return &JWTUser{
			ID:    claims.UserID,
			Email: claims.Email,
			Name:  claims.Name,
		}, nil
	}

	return nil, ErrInvalidToken
}

// GetTokenFromRequest 從請求中獲取 JWT token
func GetTokenFromRequest(c *gin.Context) string {
	// 1. 從 Cookie 中獲取
	if token, err := c.Cookie(cookieName); err == nil && token != "" {
		return token
	}

	// 2. 從 Authorization header 中獲取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 3. 從 query 參數中獲取
	if token := c.Query("token"); token != "" {
		return token
	}

	return ""
}

// SetTokenCookie 設置 JWT Cookie
func SetTokenCookie(c *gin.Context, token string) {
	c.SetCookie(
		cookieName,                  // cookie 名稱
		token,                       // cookie 值
		int(tokenTimeout.Seconds()), // 過期時間（秒）
		"/",                         // 路徑
		"",                          // 域名
		false,                       // secure（在生產環境應設為 true）
		true,                        // httpOnly
	)
}

// ClearTokenCookie 清除 JWT Cookie
func ClearTokenCookie(c *gin.Context) {
	c.SetCookie(
		cookieName, // cookie 名稱
		"",         // 空值
		-1,         // 立即過期
		"/",        // 路徑
		"",         // 域名
		false,      // secure
		true,       // httpOnly
	)
}

// JWTMiddleware JWT 中間件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := GetTokenFromRequest(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供認證 token"})
			c.Abort()
			return
		}

		user, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證 token"})
			c.Abort()
			return
		}

		// 將用戶信息存儲到上下文中
		c.Set("user", user)
		c.Next()
	}
}

// GetJWTUser 從 gin.Context 中獲取當前已認證的用戶
func GetJWTUser(c *gin.Context) (*JWTUser, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	jwtUser, ok := user.(*JWTUser)
	return jwtUser, ok
}

// GenerateTokenAndSetCookie 生成 JWT token 並設置到 Cookie 中
func GenerateTokenAndSetCookie(c *gin.Context, user *JWTUser) (string, error) {
	token, err := GenerateToken(user)
	if err != nil {
		return "", err
	}

	SetTokenCookie(c, token)
	return token, nil
}
