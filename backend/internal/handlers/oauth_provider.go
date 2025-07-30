package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type oauthProviderHandler struct{}

func (h *oauthProviderHandler) WellKnown_OpenIDConfiguration(c *gin.Context) {
	// todo: 實作 OpenID Connect 的 .well-known/openid-configuration
	c.JSON(http.StatusOK, gin.H{
		"issuer":                   "https://example.com",
		"authorization_endpoint":   "/api/oauth/authorize",
		"token_endpoint":           "/api/oauth/token",
		"userinfo_endpoint":        "/api/oauth/userinfo",
		"end_session_endpoint":     "/api/oauth/logout",
		"jwks_uri":                 "/api/oauth/jwks.json",
		"response_types_supported": []string{"code", "token"},
	})
}

func (h *oauthProviderHandler) JWKS(c *gin.Context) {
	// todo: 實作 JWKS
	// 假設這是你的 RSA 公鑰 (模擬)
	jwk := map[string]any{
		"kty": "RSA",
		"kid": "1",
		"use": "sig",
		"alg": "RS256",
		"n":   "base64-modulus",
		"e":   "AQAB",
	}
	c.JSON(http.StatusOK, map[string]any{
		"keys": []any{jwk},
	})
}

func (h *oauthProviderHandler) Authorize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth authorize (stub)"})
}

func (h *oauthProviderHandler) Token(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth token (stub)"})
}

func (h *oauthProviderHandler) UserInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth user info (stub)"})
}
func (h *oauthProviderHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth logout (stub)"})
}
