package routes

import (
	"sso/internal/handlers"

	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	app.GET("/.well-known/openid-configuration",
		handlers.OAuthProviderHandlers.WellKnown_OpenIDConfiguration)

	apiGp := app.Group("/api")
	{
		apiAuthGp := apiGp.Group("/auth")
		{
			// 一般帳密登入/登出
			apiAuthGp.POST("/login", handlers.AuthHandlers.Login)
			apiAuthGp.POST("/logout", handlers.AuthHandlers.Logout)
			apiAuthGp.POST("/register", handlers.AuthHandlers.Register)

			// 外部 OAuth 登入
			apiAuthExternalGp := apiAuthGp.Group("/external")
			{
				apiAuthExternalGp.GET("/google", handlers.AuthHandlers.GoogleLogin)
				apiAuthExternalGp.GET("/github", handlers.AuthHandlers.GithubLogin)
				apiAuthExternalGp.GET("/callback/:provider", handlers.AuthHandlers.ExternalCallback)
			}
		}

		apiOauthGp := apiGp.Group("/oauth")
		{
			apiOauthGp.GET("/jwks.json", handlers.OAuthProviderHandlers.JWKS)
			apiOauthGp.GET("/authorize", handlers.OAuthProviderHandlers.Authorize)
			apiOauthGp.POST("/token", handlers.OAuthProviderHandlers.Token)
			apiOauthGp.GET("/userinfo", handlers.OAuthProviderHandlers.UserInfo)
			apiOauthGp.POST("/logout", handlers.OAuthProviderHandlers.Logout)
		}
	}

	// 透過 JWT 獲取使用者資訊
	// app.GET("/user/profile", handlers.AuthHandlers.GetProfile)

	// OAuth 流程

	app.GET("/oauth/authorize", handlers.OAuthProviderHandlers.Authorize)
	app.POST("/oauth/token", handlers.OAuthProviderHandlers.Token)
}
