package handlers

import (
	"net/http"
	"sso/internal/services"

	"github.com/gin-gonic/gin"
)

type authHandler struct{}

func (h *authHandler) Login(c *gin.Context) {
	reqData := struct {
		Email    string `form:"email" binding:"required"`
		Password string `form:"password" binding:"required"`
	}{}
	if err := c.ShouldBind(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供 email 和密碼"})
		return
	}

	user, err := services.UserService.Login(reqData.Email, reqData.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// todo: 發配JWT

	c.JSON(http.StatusOK, gin.H{"message": "login success", "data": user})
}

func (h *authHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "logout success (stub)"})
}

func (h *authHandler) Register(c *gin.Context) {
	reqData := struct {
		Email           string `form:"email" binding:"required"`
		Password        string `form:"password" binding:"required"`
		PasswordConfirm string `form:"password_confirm" binding:"required,eqfield=Password"`
		Name            string `form:"name" binding:"required"`
		AvatarBase64    string `form:"avatar_base64"`
	}{}
	if err := c.ShouldBind(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "資料不完整，或密碼不一致"})
		return
	}

	_, err := services.UserService.Register(
		reqData.Email, reqData.Password,
		reqData.Name, reqData.AvatarBase64,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "register success"})
}
