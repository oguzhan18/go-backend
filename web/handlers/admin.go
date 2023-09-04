package handlers

import (
	"app-auth/pkg/database"
	"app-auth/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	claims, _ := c.Get("claims")
	userClaims, ok := claims.(*models.UserClaims)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Yetkisiz erişim"})
		return
	}

	isAdmin := false
	for _, role := range userClaims.Roles {
		if role == "admin" {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Yetkisiz erişim"})
		return
	}

	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := database.InsertUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı başarıyla oluşturuldu."})
}
