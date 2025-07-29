package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authHandler struct{}

func (h *authHandler) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "login success (stub)"})
}

func (h *authHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout success (stub)"})
}
