package routes

import (
	"sso/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// 一般帳密登入/登出
	app.POST("/login", handlers.AuthHandlers.Login)
	app.POST("/logout", handlers.AuthHandlers.Logout)
	// 透過 JWT 獲取使用者資訊
	// app.GET("/user/profile", handlers.AuthHandlers.GetProfile)

	// OAuth 流程
	app.GET("/oauth/authorize", handlers.OAuthHandlers.Authorize)
	app.POST("/oauth/token", handlers.OAuthHandlers.Token)
	app.GET("/oauth/callback", handlers.OAuthHandlers.Callback)
}
