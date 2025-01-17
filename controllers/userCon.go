package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"

	// "golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAllUser(c *gin.Context) {
	var user []models.User

	if err := setup.DB.Preload("Role").Find(&user).Order("username ASC").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func GetUserById(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := setup.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// func AddUser(c *gin.Context) {
// 	var input struct {
// 		Username string `gorm:"type:varchar(100)" json:"username" binding:"required"`
// 		Password string `gorm:"type:varchar(100)" json:"password" binding:"required,min=8"`
// 		RoleId   int64  `json:"role_id" binding:"required"`
// 	}
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error":   err.Error(),
// 			"valid":   false,
// 			"message": "Pastikan form sudah terisi dengan benar",
// 		})
// 		return
// 	}

// 	User := models.User{
// 		Username: input.Username,
// 		Password: input.Password,
// 		RoleId:   input.RoleId,
// 	}

// 	if err := setup.DB.Create(&User).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error":   err.Error(),
// 			"valid":   false,
// 			"message": "Gagal menambah User",
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"valid":   true,
// 		"message": "Sukses menambah User",
// 	})
// }

func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := setup.DB.Preload("Role").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
		return
	}

	var input struct {
		Username string `json:"username"`
		RoleId   int64  `json:"role_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"status":  false,
			"message": "Pastikan form sudah terisi dengan benar",
		})
		return
	}

	updates := map[string]interface{}{}
	if input.Username != "" {
		updates["username"] = input.Username
	}
	if input.RoleId != 0 {
		updates["role_id"] = input.RoleId
	}

	if err := setup.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengupdate user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Sukses mengupdate user",
	})
}

func ResetPassword(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := setup.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
		return
	}

	var input struct {
		Password        string `json:"password" binding:"required,min=8"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Password minimal 8 karakter",
		})
		return
	}

	if input.Password != input.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Password dan konfirmasi password tidak sama",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengenkripsi password",
		})
		return
	}

	if err := setup.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mereset password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Password berhasil diubah",
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := setup.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := setup.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
