package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"backend_golang/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// func Register(c *gin.Context) {
// 	var input struct {
// 		Username string `json:"username" binding:"required"`
// 		Password string `json:"password" binding:"required,min=8"`
// 		RoleId   int64  `json:"role_id"`
// 	}

// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
// 		return
// 	}

// 	user := models.User{
// 		Username: input.Username,
// 		Password: string(hashedPassword),
// 		RoleId:   input.RoleId,
// 	}

// 	if err := setup.DB.Create(&user).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

// 	if input.RoleId == 1 {
// 		stan := models.Stan{
// 			NamaPemilik: input.Username,
// 			UserId:      user.Id,
// 		}
// 		if err := setup.DB.Create(&stan).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stan"})
// 			return
// 		}
// 	} else if input.RoleId == 2 {
// 		siswa := models.Siswa{
// 			NamaSiswa: input.Username,
// 			UserId:    user.Id,
// 		}
// 		if err := setup.DB.Create(&siswa).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create siswa"})
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
// }

func RegisterSiswa(c *gin.Context) {
	var siswa models.Siswa

	NamaSiswa := c.PostForm("nama_siswa")
	if NamaSiswa == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nama siswa tidak boleh kosong",
		})
		return
	}
	siswa.NamaSiswa = NamaSiswa

	Alamat := c.PostForm("alamat")
	if Alamat == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Alamat tidak boleh kosong",
		})
		return
	}
	siswa.Alamat = Alamat

	Telp := c.PostForm("telp")
	if Telp == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nomor telepon tidak boleh kosong",
		})
		return
	}
	siswa.Telp = Telp

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
		uploadPath := "public/uploads/foto-siswa/" + filename

		if err := os.MkdirAll("public/uploads/foto-siswa", 0755); err != nil {
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
		siswa.Foto = uploadPath
	}

	Password := c.PostForm("password")
	if Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Password tidak boleh kosong",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: siswa.NamaSiswa,
		Password: string(hashedPassword),
		RoleId:   2,
	}

	if err := setup.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal membuat user: " + err.Error(),
		})
		return
	}
	siswa.UserId = user.Id // Mengaitkan UserId dengan siswa

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

func RegisterStan(c *gin.Context) {
	var stan models.Stan

	NamaStan := c.PostForm("nama_stan")
	if NamaStan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nama stan tidak boleh kosong",
		})
		return
	}
	stan.NamaStan = NamaStan

	NamaPemilik := c.PostForm("nama_pemilik")
	if NamaPemilik == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nama Pemilik tidak boleh kosong",
		})
		return
	}
	stan.NamaPemilik = NamaPemilik

	Telp := c.PostForm("telp")
	if Telp == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Nomor telepon tidak boleh kosong",
		})
		return
	}
	stan.Telp = Telp

	Password := c.PostForm("password")
	if Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Password tidak boleh kosong",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: stan.NamaPemilik,
		Password: string(hashedPassword),
		RoleId:   1,
	}

	if err := setup.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal membuat user: " + err.Error(),
		})
		return
	}
	stan.UserId = user.Id

	tx := setup.DB.Begin()
	if err := tx.Create(&stan).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Gagal menyimpan data stan: " + err.Error(),
		})
		return
	}
	tx.Commit()

	setup.DB.Preload("User").First(&stan, stan.Id)

	c.JSON(http.StatusCreated, gin.H{
		"status":  true,
		"message": "Data stan berhasil ditambahkan",
		"data":    stan,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Username   string `json:"username" binding:"required"`
		Password   string `json:"password" binding:"required"`
		RememberMe bool   `json:"remember_me"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := setup.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// if err := setup.DB.Where("password = ?", input.Password).First(&user).Error; err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
	// 	return
	// }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	var tokenDuration time.Duration
	if input.RememberMe {
		tokenDuration = 7 * 24 * time.Hour
	} else {
		tokenDuration = 24 * time.Hour
	}

	tokenString, err := utils.GenerateJWT(uint(user.Id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to generate token", "authenticated": false})
		return
	}

	c.SetCookie("Authorization", tokenString, int(tokenDuration.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"Username":      user.Username,
		"Role":          user.RoleId,
		"authenticated": true,
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	var user models.User
	if err := setup.DB.Preload("Role").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found", "status": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	// Kirim respon logout sukses
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
