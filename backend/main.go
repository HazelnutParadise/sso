package main

import (
	"sso/internal/env"
	"sso/internal/routes"
	"sso/internal/sql"

	"github.com/gin-gonic/gin"
)

func main() {
	sql.ConnectAndInitDB()

	gin.SetMode(env.GIN_MODE)

	// Initialize the application
	app := gin.Default()

	// Set up routes
	routes.Setup(app)
}
