package handlers

import (
	"net/http"
	"sso/internal/services"

	"github.com/gin-gonic/gin"
)

type authHandler struct{}

func (h *authHandler) Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	if email == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, err := services.UserService.Login(email, password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "login success", "data": user})
}

func (h *authHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout success (stub)"})
}
