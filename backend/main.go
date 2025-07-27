package main

import (
	"sso/internal/routes"
	"sso/internal/sql"

	"github.com/gin-gonic/gin"
)

func main() {
	sql.ConnectAndInitDB()

	// Initialize the application
	app := gin.Default()

	// Set up routes
	routes.Setup(app)
}
