package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type oauthHandler struct{}

func (h *oauthHandler) Authorize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth authorize (stub)"})
}

func (h *oauthHandler) Token(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth token (stub)"})
}

func (h *oauthHandler) Callback(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "oauth callback (stub)"})
}
