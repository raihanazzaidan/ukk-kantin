package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllStan(c *gin.Context) {
	var stan []models.Stan

	if err := setup.DB.Preload("User").Find(&stan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stan})
}

func GetStanById(c *gin.Context) {
	id := c.Param("id")
	var stan models.Stan

	if err := setup.DB.Preload("User").First(&stan, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stan})
}

func AddStan(c *gin.Context) {
	var stan models.Stan

	if err := c.ShouldBindJSON(&stan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data input tidak valid"})
		return
	}

	newStan := models.Stan{
		NamaStan:    stan.NamaStan,
		NamaPemilik: stan.NamaPemilik,
		Telp:        stan.Telp,
		UserId:      stan.UserId,
	}

	if err := setup.DB.Create(&newStan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan pengguna baru"})
		return
	}

	setup.DB.Preload("User").First(&newStan, newStan.Id)

	c.JSON(http.StatusCreated, gin.H{"data": newStan})
}

func UpdateStan(c *gin.Context) {
	id := c.Param("id")

	var stan models.Stan
	if err := setup.DB.Preload("User").First(&stan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
		return
	}

	var input struct {
		NamaStan    string `json:"nama_stan"`
		NamaPemilik string `json:"nama_pemilik"`
		Telp        string `json:"telp"`
		UserId      int64  `json:"user_id"`
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
	if input.NamaStan != "" {
		updates["nama_stan"] = input.NamaStan
	}
	if input.NamaPemilik != "" {
		updates["nama_pemilik"] = input.NamaPemilik
	}
	if input.Telp != "" {
		updates["telp"] = input.Telp
	}
	if input.UserId != 0 {
		updates["user_id"] = input.UserId
	}

	if err := setup.DB.Model(&stan).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengupdate stan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Sukses mengupdate stan",
	})
}

func DeleteStan(c *gin.Context) {
	id := c.Param("id")
	var stan models.Stan

	if err := setup.DB.First(&stan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stan tidak ditemukan"})
		return
	}

	if err := setup.DB.Delete(&stan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus stan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stan berhasil dihapus"})
}
// func getMenu(c *gin.Context) {
// 	userID, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
// 		return
// 	}

// 	var user models.User
// 	if err := setup.DB.First(&user, userID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
// 		return
// 	}

// 	var Menu []models.Menu

// 	if user.RoleId == 1 {
// 		if err := setup.DB.Where("user_id = ?", user.Id).Order("nama_makanan ASC").Find(&Menu).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
// 			return
// 		}
// 	} else if user.RoleId == 2 {
// 		if err := setup.DB.Order("nama_makanan ASC").Find(&Menu).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
// 			return
// 		}
// 	} else {
// 		c.JSON(http.StatusForbidden, gin.H{
// 			"error":  "Role tidak valid",
// 			"status": false,
// 		})
// 		return
// 	}

// }
