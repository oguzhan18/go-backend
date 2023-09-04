package middlewares

import (
	"app-auth/pkg/database"
	"app-auth/pkg/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header eksik"})
			c.Abort()
			return
		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz Authorization header formatı"})
			c.Abort()
			return
		}

		tokenString := authHeaderParts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Geçersiz token imza metodu")
			}
			return []byte(database.GetSecretKey()), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var userRoles []string
			roles, roleExists := claims["roles"]
			if roleExists {
				rolesStr, ok := roles.([]interface{})
				if ok {
					for _, role := range rolesStr {
						userRoles = append(userRoles, fmt.Sprint(role))
					}
				}
			}
			userClaims := models.UserClaims{
				Username: claims["username"].(string),
				Roles:    userRoles,
			}
			c.Set("claims", userClaims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Geçersiz token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
