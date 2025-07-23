package routes

import (
	"github.com/gin-gonic/gin"
)

func Setup(app *gin.Engine) {
	// Define your routes here
	app.GET("/example", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// Add more routes as needed
}
