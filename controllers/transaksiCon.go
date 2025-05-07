package controllers
import (
	middleware "backend_golang/middlewares"
	"backend_golang/models"
	"backend_golang/setup"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetTransaksi(c *gin.Context) {
	var Transaksi models.Transaksi

	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	if err := setup.DB.Find(&Transaksi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get stan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": Transaksi})
}

func GetTransaksiById(c *gin.Context) {
	id := c.Param("id")
	var Transaksi []models.Transaksi
	if err := setup.DB.
		Preload("Status").
		Preload("Siswa").
		Preload("Stan").
		First(&Transaksi, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": Transaksi})
}

func CreateTransaksi(c *gin.Context) {

	var input struct {
		StanId   int64 `json:"stan_id" binding:"required`
		SiswaId  int64 `json:"siswa_id" binding:"required`
		StatusId int64 `json:"status"  binding:"required`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if (input.StatusId){
	// 	input.StatusId = "1"
	// }

	// StandIdInt, err := strconv.ParseInt(input.StanId, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stand_id format"})
	// 	return
	// }
	// SiswaIdInt, err := strconv.ParseInt(input.SiswaId, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid siswa_id format"})
	// 	return
	// }
	// StatusIdInt, err := strconv.ParseInt(input.StatusId, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status_id format"})
	// 	return
	// }

	// Tran := models.Transaksi{
	// 	StanId: StandIdInt,
	// 	SiswaId: SiswaIdInt,
	// 	StatusId: StatusIdInt,
	// }

	Tran := models.Transaksi{
		StanId:   input.StanId,
		SiswaId:  input.SiswaId,
		StatusId: input.StatusId,
	}

	if err := setup.DB.Create(&Tran).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaksi created successfully"})
}

func UpdateStatusTransaksi(c *gin.Context) {
	id := c.Param("id")
	var TransCheck models.Transaksi

	// Cek role admin
	middleware.Admin(c)
	if c.IsAborted() {
		return
	}

	if err := setup.DB.First(&TransCheck, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		StatusId string `json:"stat" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	StatusIdInt, err := strconv.ParseInt(input.StatusId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format status_id tidak valid"})
		return
	}

	// Update status transaksi
	if err := setup.DB.Model(&TransCheck).Update("status_id", StatusIdInt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status transaksi berhasil diubah"})
}

// GetTransaksiByBulan melihat semua transaksi dalam satu bulan
func GetTransaksiByBulan(c *gin.Context) {
	// Ambil parameter bulan dan tahun
	bulan := c.Param("bulan")
	tahun := c.Param("tahun")

	// Validasi input
	bulanInt, err := strconv.Atoi(bulan)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan tidak valid"})
		return
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Hitung tanggal awal dan akhir bulan
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.Local)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Cari detail transaksi dalam rentang tanggal
	var details []models.Detail
	if err := setup.DB.
		Preload("Transaksi").
		Preload("Transaksi.Siswa").
		Preload("Menu").
		Joins("JOIN transaksis ON transaksis.id = details.transaksi_id").
		Where("transaksis.tanggal BETWEEN ? AND ?", awalBulan, akhirBulan).
		Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	// Kelompokkan detail berdasarkan transaksi
	transaksiMap := make(map[int64][]models.Detail)
	for _, d := range details {
		transaksiMap[d.TransaksiId] = append(transaksiMap[d.TransaksiId], d)
	}

	// Format response
	var response []gin.H
	for transaksiId, details := range transaksiMap {
		var total float64
		var items []gin.H
		for _, d := range details {
			total += d.HargaBeli
			items = append(items, gin.H{
				"menu":     d.Menu.NamaMakanan,
				"harga":    d.Menu.Harga,
				"qty":      d.Qty,
				"subtotal": d.HargaBeli,
			})
		}

		// Ambil data siswa dari transaksi pertama
		siswaData := gin.H{}
		if len(details) > 0 && details[0].Transaksi.SiswaId != 0 {
			siswaData = gin.H{
				"id":         details[0].Transaksi.Siswa.Id,
				"nama_siswa": details[0].Transaksi.Siswa.NamaSiswa,
				"alamat":     details[0].Transaksi.Siswa.Alamat,
				"telp":       details[0].Transaksi.Siswa.Telp,
			}
		}

		response = append(response, gin.H{
			"id":      transaksiId,
			"tanggal": details[0].Transaksi.Tanggal.Format("02-01-2006"),
			"total":   total,
			"siswa":   siswaData,
			"items":   items,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"bulan": bulanInt,
		"tahun": tahunInt,
		"data":  response,
	})
}

// GetRekapBulanan melihat rekap pendapatan per bulan
func GetRekapBulanan(c *gin.Context) {
	// Ambil parameter tahun
	tahun := c.Param("tahun")

	// Validasi input
	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Buat query untuk menghitung total per bulan
	query := `
		SELECT 
			MONTH(transaksis.tanggal) as bulan,
			COUNT(DISTINCT details.transaksi_id) as jumlah_transaksi,
			SUM(details.harga_beli) as total_pendapatan
		FROM details
		JOIN transaksis ON transaksis.id = details.transaksi_id
		WHERE YEAR(transaksis.tanggal) = ?
		GROUP BY MONTH(transaksis.tanggal)
		ORDER BY bulan
	`

	var rekap []struct {
		Bulan           int     `json:"bulan"`
		JumlahTransaksi int     `json:"jumlah_transaksi"`
		TotalPendapatan float64 `json:"total_pendapatan"`
	}

	if err := setup.DB.Raw(query, tahunInt).Scan(&rekap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data rekap"})
		return
	}

	// Format nama bulan
	namaBulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	var response []gin.H
	for _, r := range rekap {
		response = append(response, gin.H{
			"bulan":            namaBulan[r.Bulan-1],
			"jumlah_transaksi": r.JumlahTransaksi,
			"total_pendapatan": r.TotalPendapatan,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"tahun": tahunInt,
		"data":  response,
	})
}

// CreateTransaksiWithDetail membuat transaksi dan detail dalam satu request
func CreateTransaksiWithDetail(c *gin.Context) {
	var input struct {
		StanId   int64 `json:"stan_id" binding:"required"`
		SiswaId  int64 `json:"siswa_id" binding:"required"`
		StatusId int64 `json:"status" binding:"required"`
		Items    []struct {
			MenuId int64 `json:"menu_id" binding:"required"`
			Qty    int64 `json:"qty" binding:"required"`
		} `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mulai transaksi database
	tx := setup.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	// Buat transaksi
	transaksi := models.Transaksi{
		StanId:   input.StanId,
		SiswaId:  input.SiswaId,
		StatusId: input.StatusId,
	}

	if err := tx.Create(&transaksi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat transaksi"})
		return
	}

	// Buat detail untuk setiap item
	var total float64
	var details []models.Detail
	for _, item := range input.Items {
		// Ambil data menu
		var menu models.Menu
		if err := tx.First(&menu, item.MenuId).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Menu tidak ditemukan"})
			return
		}

		// Cek diskon yang berlaku untuk menu ini
		var menuDiskon models.MenuDiskon
		today := time.Now()
		if err := tx.
			Preload("Diskon").
			Joins("JOIN diskons ON diskons.id = menu_diskons.diskon_id").
			Where("menu_diskons.menu_id = ? AND diskons.tanggal_awal <= ? AND diskons.tanggal_akhir >= ?",
				item.MenuId, today, today).
			First(&menuDiskon).Error; err == nil {
			// Ada diskon
			presentaseDiskon, err := strconv.ParseInt(menuDiskon.Diskon.PresentaseDiskon, 10, 64)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Format diskon tidak valid"})
				return
			}

			hargaDiskon := menu.Harga * (1 - float64(presentaseDiskon)/100)
			total += hargaDiskon * float64(item.Qty)

			detail := models.Detail{
				TransaksiId: transaksi.Id,
				MenuId:      item.MenuId,
				Qty:         item.Qty,
				HargaBeli:   hargaDiskon * float64(item.Qty),
				DiskonId:    menuDiskon.DiskonId,
			}
			details = append(details, detail)
		} else {
			// Tidak ada diskon
			total += menu.Harga * float64(item.Qty)

			detail := models.Detail{
				TransaksiId: transaksi.Id,
				MenuId:      item.MenuId,
				Qty:         item.Qty,
				HargaBeli:   menu.Harga * float64(item.Qty),
				DiskonId:    0, // Set ke 0 untuk menandakan tidak ada diskon
			}
			details = append(details, detail)
		}
	}

	// Simpan semua detail
	for _, detail := range details {
		if err := tx.Create(&detail).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat detail transaksi"})
			return
		}
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Transaksi berhasil dibuat",
		"data": gin.H{
			"transaksi": transaksi,
			"details":   details,
			"total":     total,
		},
	})
}

// GetHistoriTransaksiSiswa melihat histori transaksi siswa
func GetHistoriTransaksiSiswa(c *gin.Context) {
	// Ambil siswa_id dari parameter
	siswaId := c.Param("siswa_id")
	bulan := c.Param("bulan")
	tahun := c.Param("tahun")

	// Konversi parameter ke int64
	siswaIdInt, err := strconv.ParseInt(siswaId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format siswa_id tidak valid"})
		return
	}

	bulanInt, err := strconv.Atoi(bulan)
	if err != nil || bulanInt < 1 || bulanInt > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bulan tidak valid"})
		return
	}

	tahunInt, err := strconv.Atoi(tahun)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tahun tidak valid"})
		return
	}

	// Hitung tanggal awal dan akhir bulan
	awalBulan := time.Date(tahunInt, time.Month(bulanInt), 1, 0, 0, 0, 0, time.Local)
	akhirBulan := awalBulan.AddDate(0, 1, -1)

	// Cari detail transaksi berdasarkan siswa_id dan bulan
	var details []models.Detail
	if err := setup.DB.
		Preload("Menu").
		Preload("Transaksi").
		Preload("Transaksi.Stan").
		Preload("Transaksi.Status").
		Joins("JOIN transaksis ON transaksis.id = details.transaksi_id").
		Where("transaksis.siswa_id = ? AND transaksis.tanggal BETWEEN ? AND ?",
			siswaIdInt, awalBulan, akhirBulan).
		Order("transaksis.tanggal DESC").
		Find(&details).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data transaksi"})
		return
	}

	// Kelompokkan detail berdasarkan transaksi
	transaksiMap := make(map[int64][]models.Detail)
	for _, d := range details {
		transaksiMap[d.TransaksiId] = append(transaksiMap[d.TransaksiId], d)
	}

	// Format response
	var response []gin.H
	for transaksiId, details := range transaksiMap {
		var total float64
		var items []gin.H
		for _, d := range details {
			total += d.HargaBeli
			items = append(items, gin.H{
				"menu":     d.Menu.NamaMakanan,
				"harga":    d.Menu.Harga,
				"qty":      d.Qty,
				"subtotal": d.HargaBeli,
			})
		}

		response = append(response, gin.H{
			"id":      transaksiId,
			"tanggal": details[0].Transaksi.Tanggal.Format("02-01-2006 15:04:05"),
			"stan":    details[0].Transaksi.Stan.NamaStan,
			"status":  details[0].Transaksi.Status.Status,
			"total":   total,
			"items":   items,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"siswa_id": siswaIdInt,
		"bulan":    bulanInt,
		"tahun":    tahunInt,
		"data":     response,
	})
}