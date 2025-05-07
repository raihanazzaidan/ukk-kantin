package controllers

import (
	"net/http"
	"time"

	middleware "backend_golang/middlewares"
	"backend_golang/models"
	"backend_golang/setup"

	"github.com/gin-gonic/gin"
)

func GetAllDiskon(c *gin.Context) {
	var diskon []models.Diskon

	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	if err := setup.DB.Find(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": diskon})
}

func GetDiskonById(c *gin.Context) {
	var diskon models.Diskon
	id := c.Param("id")
	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}
	if err := setup.DB.First(&diskon, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": diskon})
}

func AddDiskon(c *gin.Context) {
	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	var input struct {
		NamaDiskon       string `json:"nama_diskon" binding:"required"`
		PresentaseDiskon string `json:"presentase_diskon" binding:"required"`
		TanggalAwal      string `json:"tanggal_awal" binding:"required"`
		TanggalAkhir     string `json:"tanggal_akhir" binding:"required"`
		MenuId           int64  `json:"menu_id" binding:"required"`
		StanId           int64  `json:"stan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse tanggal
	tanggalAwal, err := time.Parse("2006-01-02", input.TanggalAwal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_awal tidak valid. Gunakan format YYYY-MM-DD"})
		return
	}

	tanggalAkhir, err := time.Parse("2006-01-02", input.TanggalAkhir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_akhir tidak valid. Gunakan format YYYY-MM-DD"})
		return
	}

	// Mulai transaksi database
	tx := setup.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	// Buat diskon
	diskon := models.Diskon{
		NamaDiskon:       input.NamaDiskon,
		PresentaseDiskon: input.PresentaseDiskon,
		TanggalAwal:      tanggalAwal,
		TanggalAkhir:     tanggalAkhir,
		StanId:           input.StanId,
	}

	if err := tx.Create(&diskon).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat diskon"})
		return
	}

	// Buat relasi menu_diskon
	menuDiskon := models.MenuDiskon{
		MenuId:   input.MenuId,
		DiskonId: diskon.Id,
		StanId:   input.StanId,
	}

	if err := tx.Create(&menuDiskon).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghubungkan diskon dengan menu"})
		return
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan diskon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Diskon berhasil dibuat",
		"data": gin.H{
			"diskon":      diskon,
			"menu_diskon": menuDiskon,
		},
	})
}

func UpdateDiskon(c *gin.Context) {
	id := c.Param("id")
	var diskon models.Diskon

	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	if err := setup.DB.First(&diskon, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		NamaDiskon       string `json:"nama_diskon" binding:"required"`
		PresentaseDiskon string `json:"presentase_diskon" binding:"required"`
		TanggalAwal      string `json:"tanggal_awal" binding:"required"`
		TanggalAkhir     string `json:"tanggal_akhir" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse tanggal
	tanggalAwal, err := time.Parse("2006-01-02", input.TanggalAwal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_awal tidak valid. Gunakan format YYYY-MM-DD"})
		return
	}

	tanggalAkhir, err := time.Parse("2006-01-02", input.TanggalAkhir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal_akhir tidak valid. Gunakan format YYYY-MM-DD"})
		return
	}

	diskon.NamaDiskon = input.NamaDiskon
	diskon.PresentaseDiskon = input.PresentaseDiskon
	diskon.TanggalAwal = tanggalAwal
	diskon.TanggalAkhir = tanggalAkhir

	if err := setup.DB.Save(&diskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diskon berhasil diupdate"})
}

func DeleteDiskon(c *gin.Context) {
	var diskon models.Diskon
	id := c.Param("id")

	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	if err := setup.DB.Delete(&diskon, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diskon deleted successfully"})
}

func AddDiskonToMenu(c *gin.Context) {
	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	var input struct {
		MenuId   int64 `json:"menu_id" binding:"required"`
		DiskonId int64 `json:"diskon_id" binding:"required"`
		StanId   int64 `json:"stan_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	menuDiskon := models.MenuDiskon{
		MenuId:   input.MenuId,
		DiskonId: input.DiskonId,
		StanId:   input.StanId,
	}

	if err := setup.DB.Create(&menuDiskon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Diskon berhasil ditambahkan ke menu"})
}
