package handlers

import (
	"app-auth/pkg/database"
	"app-auth/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := database.AuthenticateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !authUser {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Kullanıcı adı veya şifre hatalı"})
		return
	}

	token, err := database.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulurken hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
