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

func GetAllSiswa(c *gin.Context) {
	var Siswa []models.Siswa

	if err := setup.DB.Preload("User").Find(&Siswa).Order("name ASC").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": Siswa})
}

func GetSiswaById(c *gin.Context) {
	id := c.Param("id")
	var siswa models.Siswa

	if err := setup.DB.Preload("User").First(&siswa, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": siswa})
}

func AddSiswa(c *gin.Context) {
	var siswa models.Siswa

	namaSiswa := c.PostForm("nama_siswa")
	if namaSiswa == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nama siswa tidak boleh kosong",
		})
		return
	}
	siswa.NamaSiswa = namaSiswa

	alamat := c.PostForm("alamat")
	if alamat == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Alamat tidak boleh kosong",
		})
		return
	}
	siswa.Alamat = alamat

	telp := c.PostForm("telp")
	if telp == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nomor telepon tidak boleh kosong",
		})
		return
	}
	siswa.Telp = telp

	userId := c.PostForm("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User ID tidak boleh kosong",
		})
		return
	}

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "User ID harus berupa angka",
		})
		return
	}
	siswa.UserId = userIdInt

	// Proses upload foto
	foto, err := c.FormFile("foto")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(foto.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Format file tidak didukung. Gunakan JPG, JPEG, atau PNG",
			})
			return
		}

		if foto.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Ukuran file maksimal 5MB",
			})
			return
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%d_%s", timestamp, foto.Filename)
		uploadPath := "public/uploads/foto-siswa/" + filename

		if err := os.MkdirAll("public/uploads/foto-siswa", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal membuat direktori upload",
			})
			return
		}

		if err := c.SaveUploadedFile(foto, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal menyimpan foto",
			})
			return
		}
		siswa.Foto = uploadPath
	}

	tx := setup.DB.Begin()
	if err := tx.Create(&siswa).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal menyimpan data siswa: " + err.Error(),
		})
		return
	}
	tx.Commit()

	setup.DB.Preload("User").First(&siswa, siswa.Id)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Data siswa berhasil ditambahkan",
		"data":    siswa,
	})
}

func UpdateSiswa(c *gin.Context) {
	id := c.Param("id")
	var siswa models.Siswa

	if err := setup.DB.Preload("User").First(&siswa, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  false,
			"message": "Data siswa tidak ditemukan",
		})
		return
	}

	if namaSiswa := c.PostForm("nama_siswa"); namaSiswa != "" {
		siswa.NamaSiswa = namaSiswa
	}

	if alamat := c.PostForm("alamat"); alamat != "" {
		siswa.Alamat = alamat
	}

	if telp := c.PostForm("telp"); telp != "" {
		siswa.Telp = telp
	}

	if userId := c.PostForm("user_id"); userId != "" {
		userIdInt, err := strconv.ParseInt(userId, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "User ID harus berupa angka",
			})
			return
		}
		siswa.UserId = userIdInt
	}

	foto, err := c.FormFile("foto")
	if err == nil {
		ext := strings.ToLower(filepath.Ext(foto.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Format file tidak didukung. Gunakan JPG, JPEG, atau PNG",
			})
			return
		}

		if foto.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Ukuran file maksimal 5MB",
			})
			return
		}

		if siswa.Foto != "" {
			if err := os.Remove(siswa.Foto); err != nil {
				fmt.Printf("Error menghapus foto lama: %v\n", err)
			}
		}

		timestamp := time.Now().Unix()
		filename := fmt.Sprintf("%d_%s", timestamp, foto.Filename)
		uploadPath := "public/uploads/foto-siswa/" + filename

		if err := os.MkdirAll("public/uploads/foto-siswa", 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal membuat direktori upload",
			})
			return
		}

		if err := c.SaveUploadedFile(foto, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Gagal menyimpan foto",
			})
			return
		}
		siswa.Foto = uploadPath
	}

	tx := setup.DB.Begin()
	if err := tx.Save(&siswa).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal mengupdate data siswa: " + err.Error(),
		})
		return
	}
	tx.Commit()

	setup.DB.Preload("User").First(&siswa, siswa.Id)

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Data siswa berhasil diupdate",
		"data":    siswa,
	})
}

func DeteleSiswa(c *gin.Context) {
	id := c.Param("id")
	var siswa models.Siswa

	if err := setup.DB.First(&siswa, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pengguna tidak ditemukan"})
		return
	}

	if err := setup.DB.Delete(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pengguna"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil dihapus"})
}
