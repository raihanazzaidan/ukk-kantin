package controllers

import (
	"backend_golang/models"
	"backend_golang/setup"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
)

func GetDetail(c *gin.Context) {

	var detail []models.Detail

	if err := setup.DB.First(&detail).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func GetDetailById(c *gin.Context) {

	id := c.Param("id")
	var detail []models.Detail

	if err := setup.DB.First(&detail, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"data": detail})
}

func CreateDetail(c *gin.Context) {
	var input struct {
		TransaksiId string `json:"transaksi_id" binding:"required"`
		MenuId      string `json:"menu_id" binding:"required"`
		Qty         string `json:"qty" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	TransaksiIdInt, err := strconv.ParseInt(input.TransaksiId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaksi_id format"})
		return
	}
	MenuIdInt, err := strconv.ParseInt(input.MenuId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu_id format"})
		return
	}
	QtyInt, err := strconv.ParseInt(input.Qty, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid qty format"})
		return
	}

	// Ambil data menu
	var Menu models.Menu
	if err := setup.DB.First(&Menu, input.MenuId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Menu not found"})
		return
	}

	// Ambil data diskon yang berlaku hari ini
	var diskon models.Diskon
	today := time.Now().Format("2006-01-02")
	if err := setup.DB.Where("tanggal_mulai <= ? AND tanggal_selesai >= ?", today, today).First(&diskon).Error; err != nil {
		// Jika tidak ada diskon, gunakan harga normal
		Final := Menu.Harga * float64(QtyInt)
		detail := models.Detail{
			TransaksiId: TransaksiIdInt,
			MenuId:      MenuIdInt,
			Qty:         QtyInt,
			HargaBeli:   Final,
		}

		if err := setup.DB.Create(&detail).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Detail Successfully Created", "data": detail})
		return
	}

	// Konversi PresentaseDiskon dari string ke int
	presentaseDiskon, err := strconv.ParseInt(diskon.PresentaseDiskon, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid diskon percentage"})
		return
	}

	// Hitung harga dengan diskon
	hargaDiskon := Menu.Harga * (1 - float64(presentaseDiskon)/100)
	Final := hargaDiskon * float64(QtyInt)

	detail := models.Detail{
		TransaksiId: TransaksiIdInt,
		MenuId:      MenuIdInt,
		Qty:         QtyInt,
		HargaBeli:   Final,
	}

	if err := setup.DB.Create(&detail).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail Successfully Created",
		"data":    detail,
		"diskon": gin.H{
			"id":                   diskon.Id,
			"presentase":           diskon.PresentaseDiskon,
			"harga_awal":           Menu.Harga,
			"harga_setelah_diskon": hargaDiskon,
		},
	})
}

func PrintNota(c *gin.Context) {
	// Ambil ID transaksi dari parameter
	transaksiID := c.Param("id")

	// Cari detail transaksi berdasarkan ID transaksi
	var details []models.Detail
	if err := setup.DB.
		Preload("Menu").
		Preload("Transaksi").
		Preload("Transaksi.Siswa").
		Preload("Diskon").
		Where("transaksi_id = ?", transaksiID).
		Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get transaction details"})
		return
	}

	if len(details) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Buat PDF baru
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "NOTA PEMBELIAN")
	pdf.Ln(20)

	// Informasi transaksi
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("No. Transaksi: %d", details[0].TransaksiId))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Tanggal: %s", details[0].Transaksi.Tanggal.Format("02-01-2006 15:04:05")))
	pdf.Ln(10)
	pdf.Cell(40, 10, fmt.Sprintf("Nama Pembeli: %s", details[0].Transaksi.Siswa.NamaSiswa))
	pdf.Ln(20)

	// Header tabel
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(40, 10, "Menu")
	pdf.Cell(30, 10, "Harga")
	pdf.Cell(20, 10, "Qty")
	pdf.Cell(30, 10, "Diskon")
	pdf.Cell(40, 10, "Subtotal")
	pdf.Ln(10)

	// Isi tabel
	pdf.SetFont("Arial", "", 12)
	var total float64
	for _, detail := range details {
		pdf.Cell(40, 10, detail.Menu.NamaMakanan)
		pdf.Cell(30, 10, fmt.Sprintf("Rp %.2f", detail.Menu.Harga))
		pdf.Cell(20, 10, fmt.Sprintf("%d", detail.Qty))

		// Tampilkan diskon jika ada
		if detail.DiskonId != 0 {
			presentaseDiskon, err := strconv.ParseFloat(detail.Diskon.PresentaseDiskon, 64)
			if err != nil {
				pdf.Cell(30, 10, "-") // Jika konversi gagal, tampilkan "-"
			} else {
				pdf.Cell(30, 10, fmt.Sprintf("%.0f%%", presentaseDiskon))
			}
		} else {
			pdf.Cell(30, 10, "-")
		}

		subtotal := detail.HargaBeli
		pdf.Cell(40, 10, fmt.Sprintf("Rp %.2f", subtotal))
		pdf.Ln(10)
		total += subtotal
	}

	// Total
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(120, 10, "TOTAL")
	pdf.Cell(40, 10, fmt.Sprintf("Rp %.2f", total))

	// Simpan PDF ke buffer
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate PDF"})
		return
	}

	// Set header untuk download PDF
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=nota_%s.pdf", transaksiID))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(buf.Len()))

	// Kirim PDF
	c.Data(http.StatusOK, "application/pdf", buf.Bytes())
}
