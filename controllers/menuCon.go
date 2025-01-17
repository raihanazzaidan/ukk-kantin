package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllMenu(c *gin.Context) {
	var menu []models.Menu

	if err := setup.DB.Preload("Jenis").Preload("Stan").Find(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func GetMenuByStanId(c *gin.Context) {
	id := c.Param("stan_id")
	var menu []models.Menu

	if err := setup.DB.Preload("Jenis").Preload("Stan").Where("stan_id = ?", id).Find(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"error":  err.Error(),
		})
		return
	}

	if len(menu) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Tidak ada menu yang ditemukan untuk Stan ID ini",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"data":   menu,
	})
}

func AddMenu(c *gin.Context) {
	var menu models.Menu

	NamaMakanan := c.PostForm("nama_makanan")
	if NamaMakanan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nama makanan tidak boleh kosong",
		})
		return
	}
	menu.NamaMakanan = NamaMakanan

	Harga := c.PostForm("harga")
	if Harga == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Harga tidak boleh kosong",
		})
		return
	}

	hargaFloat, err := strconv.ParseFloat(Harga, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Harga harus berupa angka",
		})
		return
	}
	menu.Harga = hargaFloat

	JenisId := c.PostForm("jenis_id")
	if JenisId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nomor telepon tidak boleh kosong",
		})
		return
	}

	jenisIdInt, err := strconv.ParseInt(JenisId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Jenis ID harus berupa angka",
		})
		return
	}
	menu.JenisId = jenisIdInt

	Foto, err := c.FormFile("foto")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(Foto.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Format file tidak didukung. Gunakan JPG, JPEG, atau PNG",
			})
			return
		}

		if Foto.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Ukuran file maksimal 5MB",
			})
			return
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%d_%s", timestamp, Foto.Filename)
		uploadPath := "public/uploads/foto-menu/" + filename

		if err := os.MkdirAll("public/uploads/foto-menu", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal membuat direktori upload",
			})
			return
		}

		if err := c.SaveUploadedFile(Foto, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal menyimpan foto",
			})
			return
		}
		menu.Foto = uploadPath
	}

	Deskripsi := c.PostForm("deskripsi")
	if Deskripsi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Deskripsi tidak boleh kosong",
		})
		return
	}
	menu.Deskripsi = Deskripsi

	StanId := c.PostForm("stan_id")
	if StanId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Id Stan tidak boleh kosong",
		})
		return
	}

	StanIdInt, err := strconv.ParseInt(StanId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Jenis ID harus berupa angka",
		})
		return
	}
	menu.StanId = StanIdInt

	tx := setup.DB.Begin()
	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal menyimpan data menu: " + err.Error(),
		})
		return
	}
	tx.Commit()

	setup.DB.Preload("Jenis").Preload("Stan").First(&menu, menu.Id)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Data menu berhasil ditambahkan",
		"data":    menu,
	})
}

func UpdateMenu(c *gin.Context) {
	id := c.Param("id")

	var menu models.Menu
	if err := setup.DB.Preload("Jenis").Preload("Stan").First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Data tidak ditemukan",
		})
		return
	}

	// Menyimpan foto lama
	oldFoto := menu.Foto

	// Mengambil data form
	NamaMakanan := c.PostForm("nama_makanan")
	if NamaMakanan != "" {
		menu.NamaMakanan = NamaMakanan
	}

	Harga := c.PostForm("harga")
	if Harga != "" {
		hargaFloat, err := strconv.ParseFloat(Harga, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Harga harus berupa angka",
			})
			return
		}
		menu.Harga = hargaFloat
	}

	JenisId := c.PostForm("jenis_id")
	if JenisId != "" {
		jenisIdInt, err := strconv.ParseInt(JenisId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Jenis ID harus berupa angka",
			})
			return
		}
		menu.JenisId = jenisIdInt
	}

	Foto, err := c.FormFile("foto")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(Foto.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Format file tidak didukung. Gunakan JPG, JPEG, atau PNG",
			})
			return
		}

		if Foto.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Ukuran file maksimal 5MB",
			})
			return
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%d_%s", timestamp, Foto.Filename)
		uploadPath := "public/uploads/foto-menu/" + filename

		if err := os.MkdirAll("public/uploads/foto-menu", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal membuat direktori upload",
			})
			return
		}

		if err := c.SaveUploadedFile(Foto, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal menyimpan foto",
			})
			return
		}

		if oldFoto != "" {
			os.Remove(oldFoto)
		}
		menu.Foto = uploadPath
	}

	Deskripsi := c.PostForm("deskripsi")
	if Deskripsi != "" {
		menu.Deskripsi = Deskripsi
	}

	StanId := c.PostForm("stan_id")
	if StanId != "" {
		stanIdInt, err := strconv.ParseInt(StanId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Stan ID harus berupa angka",
			})
			return
		}
		menu.StanId = stanIdInt
	}

	tx := setup.DB.Begin()
	if err := tx.Save(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengupdate menu: " + err.Error(),
		})
		return
	}
	tx.Commit()

	setup.DB.Preload("Jenis").Preload("Stan").First(&menu, menu.Id)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Menu berhasil diupdate",
		"data":    menu,
	})
}

func DeleteMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu

	if err := setup.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
		return
	}

	if err := setup.DB.Delete(&menu).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus menu"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus"})
}
