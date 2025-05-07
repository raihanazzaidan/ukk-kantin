package middleware

import (
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"

	"github.com/gin-gonic/gin"
)

func checkRole(c *gin.Context, requiredRole int64) {
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	// Cari user berdasarkan ID
	var user models.User
	if err := setup.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var id = user.RoleId

	// Cek role user
	if id != requiredRole {
		c.JSON(http.StatusForbidden, gin.H{
			"error":         "Restricted Action",
		})
		c.Abort()
		return
	}

	c.Next()
}

func Admin(c *gin.Context) {
	checkRole(c, 1)
}

func Siswa(c *gin.Context) {
	checkRole(c, 2)
}