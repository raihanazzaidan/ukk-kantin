// this page used for middleware
package middleware

import (
	"backend_golang/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie not found", "status": false})
			c.Abort()
			return
		}

		// token validator
		token, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "status": false})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			expiration := int64(claims["exp"].(float64))
			if time.Now().Unix() > expiration {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired", "status": false})
				c.Abort()
				return
			}
			c.Set("user", claims["sub"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "status": false})
			c.Abort()
		}
	}
}
