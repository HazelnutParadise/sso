package main

import (
	"net/http"
	"sso/internal/env"
	"sso/internal/routes"
	"sso/internal/sql"

	"github.com/gin-gonic/gin"
)

// CORS 中間件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000") // 前端域名
		c.Header("Access-Control-Allow-Credentials", "true")             // 允許發送 Cookie
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	sql.ConnectAndInitDB()

	gin.SetMode(env.GIN_MODE)

	// Initialize the application
	app := gin.Default()

	// 添加 CORS 中間件
	app.Use(corsMiddleware())

	// Set up routes
	routes.Setup(app)

	// 啟動伺服器
	port := "8080" // 可以從環境變數讀取
	app.Run(":" + port)
}
