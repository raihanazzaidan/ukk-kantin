package controllers

import (
	"net/http"

	"backend_golang/models"
	"backend_golang/setup"

	"github.com/gin-gonic/gin"
)

func AddDiskon(c *gin.Context) {
	var diskon models.Diskon
	if err := c.ShouldBindJSON(&diskon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi tanggal
	if diskon.TanggalAwal.After(diskon.TanggalAkhir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal awal tidak boleh lebih besar dari tanggal akhir"})
		return
	}

	if err := setup.DB.Create(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": diskon})
}

func GetAllDiskon(c *gin.Context) {
	var diskons []models.Diskon
	if err := setup.DB.Preload("Stan").Find(&diskons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": diskons})
}

func GetDiskonById(c *gin.Context) {
	id := c.Param("id")
	var diskon models.Diskon

	if err := setup.DB.Preload("Stan").First(&diskon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diskon tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": diskon})
}

func UpdateDiskon(c *gin.Context) {
	id := c.Param("id")
	var diskon models.Diskon

	if err := setup.DB.First(&diskon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diskon tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&diskon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := setup.DB.Save(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": diskon})
}

func DeleteDiskon(c *gin.Context) {
	id := c.Param("id")
	var diskon models.Diskon

	if err := setup.DB.First(&diskon, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diskon tidak ditemukan"})
		return
	}

	if err := setup.DB.Delete(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diskon berhasil dihapus"})
}
