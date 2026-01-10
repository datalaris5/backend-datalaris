package controllers

import (
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/dto"
	"go-datalaris/models"
	"go-datalaris/services"
	"go-datalaris/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDashboardPesananTotalPenjualan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}
	db := config.DB
	layout := "2006-01-02"

	// convert string to time.Time
	dateFrom, err1 := time.Parse(layout, input.DateFrom)
	dateTo, err2 := time.Parse(layout, input.DateTo)
	if err1 != nil || err2 != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD", nil)
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64
	if marketplace.Name == constant.ShopeeConst {
		totalQuery := `
		SELECT COALESCE(SUM(total_harga_produk), 0)
    	FROM shopee_data_upload_pesanan_details
		WHERE waktu_pesanan_dibuat BETWEEN ? AND ?
		`
		currentArgs := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			totalQuery += " AND store_id = ?"
			currentArgs = append(currentArgs, input.StoreId)
		}

		// total penjualan periode sekarang
		db.Raw(totalQuery, currentArgs...).Scan(&current)

		// cari range periode sebelumnya
		days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
		periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
		periodBeforeTo := dateTo.AddDate(0, 0, -days)

		// total periode sebelumnya
		previousArgs := []interface{}{
			periodBeforeFrom,
			periodBeforeTo,
		}

		if input.StoreId != 0 {
			previousArgs = append(previousArgs, input.StoreId)
		}

		// total periode sebelumnya
		db.Raw(totalQuery, previousArgs...).Scan(&previous)

		// hitung persentase naik / turun
		changePercent := 0.0
		if previous > 0 {
			changePercent = ((current - previous) / previous) * 100
		}
		if changePercent > 0 {
			result.Trend = "Up"
		} else if changePercent == 0 {
			result.Trend = "Equal"
		} else {
			result.Trend = "Down"
		}

		result.Total = current
		result.Percent = changePercent

		sparklineQuery := `
		SELECT waktu_pesanan_dibuat AS tanggal, SUM(total_harga_produk) AS total
		FROM shopee_data_upload_pesanan_details
		WHERE waktu_pesanan_dibuat BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY waktu_pesanan_dibuat ORDER BY waktu_pesanan_dibuat ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananTotalPesanan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}
	db := config.DB
	layout := "2006-01-02"

	// convert string to time.Time
	dateFrom, err1 := time.Parse(layout, input.DateFrom)
	dateTo, err2 := time.Parse(layout, input.DateTo)
	if err1 != nil || err2 != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD", nil)
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64
	if marketplace.Name == constant.ShopeeConst {
		totalQuery := `
		SELECT COALESCE(SUM(jumlah), 0)
    	FROM shopee_data_upload_pesanan_details
		WHERE waktu_pesanan_dibuat BETWEEN ? AND ?
		`
		currentArgs := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			totalQuery += " AND store_id = ?"
			currentArgs = append(currentArgs, input.StoreId)
		}

		// total penjualan periode sekarang
		db.Raw(totalQuery, currentArgs...).Scan(&current)

		// cari range periode sebelumnya
		days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
		periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
		periodBeforeTo := dateTo.AddDate(0, 0, -days)

		// total periode sebelumnya
		previousArgs := []interface{}{
			periodBeforeFrom,
			periodBeforeTo,
		}

		if input.StoreId != 0 {
			previousArgs = append(previousArgs, input.StoreId)
		}

		// total periode sebelumnya
		db.Raw(totalQuery, previousArgs...).Scan(&previous)

		// hitung persentase naik / turun
		changePercent := 0.0
		if previous > 0 {
			changePercent = ((current - previous) / previous) * 100
		}
		if changePercent > 0 {
			result.Trend = "Up"
		} else if changePercent == 0 {
			result.Trend = "Equal"
		} else {
			result.Trend = "Down"
		}

		result.Total = current
		result.Percent = changePercent

		sparklineQuery := `
		SELECT waktu_pesanan_dibuat AS tanggal, SUM(jumlah) AS total
		FROM shopee_data_upload_pesanan_details
		WHERE waktu_pesanan_dibuat BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY waktu_pesanan_dibuat ORDER BY waktu_pesanan_dibuat ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananTopUsernamePembeli(c *gin.Context) {
	var result []dto.ResponseTopUsernamePembeliDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT username_pembeli AS username, COUNT(username_pembeli) AS total FROM shopee_data_upload_pesanan_details WHERE waktu_pesanan_dibuat BETWEEN ? AND ? `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY username_pembeli ORDER BY COUNT(username_pembeli) DESC"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananTopMetodePembayaran(c *gin.Context) {
	var result []dto.ResponseTopMetodePembayaranDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT metode_pembayaran, COUNT(metode_pembayaran) AS total FROM shopee_data_upload_pesanan_details WHERE waktu_pesanan_dibuat BETWEEN ? AND ? `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY metode_pembayaran ORDER BY COUNT(metode_pembayaran) DESC"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananStatusPesanan(c *gin.Context) {
	var result []dto.ResponseStatusPesananDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT 
		CASE
			WHEN status_pesanan ILIKE 'Pesanan Diterima%' 
				THEN 'Pesanan Diterima'
			ELSE status_pesanan
		END AS status_pesanan,
		COUNT(status_pesanan) AS total
		FROM shopee_data_upload_pesanan_details WHERE waktu_pesanan_dibuat BETWEEN ? AND ? `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += ` GROUP BY
				CASE
					WHEN status_pesanan ILIKE 'Pesanan Diterima%' 
						THEN 'Pesanan Diterima'
					ELSE status_pesanan
				END
			ORDER BY COUNT(status_pesanan) DESC`

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananOpsiPengiriman(c *gin.Context) {
	var result []dto.ResponseOpsiPengirimanDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT opsi_pengiriman, COUNT(opsi_pengiriman) AS total FROM shopee_data_upload_pesanan_details WHERE waktu_pesanan_dibuat BETWEEN ? AND ? `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY opsi_pengiriman ORDER BY COUNT(opsi_pengiriman) DESC"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananMetodePickup(c *gin.Context) {
	var result []dto.ResponseMetodePickupDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT counter_pickup AS metode_pickup, COUNT(counter_pickup) AS total FROM shopee_data_upload_pesanan_details WHERE waktu_pesanan_dibuat BETWEEN ? AND ? `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY counter_pickup ORDER BY COUNT(counter_pickup) DESC"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}

func GetDashboardPesananByDay(c *gin.Context) {
	var result []dto.ResponseTotalPesananByDayDashboardPesanan
	input, errBind := utils.BindJSON[dto.RequestDashboardPesanan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	db := config.DB

	if marketplace.Name == constant.ShopeeConst {

		query := `
		SELECT
    EXTRACT(ISODOW FROM waktu_pesanan_dibuat) AS no,
    TRIM(TO_CHAR(waktu_pesanan_dibuat, 'Day')) AS day,
    COUNT(*) AS total
	FROM shopee_data_upload_pesanan_details
	WHERE waktu_pesanan_dibuat BETWEEN ? AND ?  `

		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY no, day ORDER BY no;"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardPesananConst+constant.SuccessFetch, result)
}
