package main

import (
	"sso/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the application
	app := gin.Default()

	// Set up routes
	routes.Setup(app)
}
