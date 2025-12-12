package controllers

import (
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/dto"
	"go-datalaris/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetDashboardIklanPenjualanIklan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
    SELECT COALESCE(SUM(omzet_penjualan), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
    SELECT COALESCE(SUM(omzet_penjualan), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT tanggal, SUM(omzet_penjualan) AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanBiayaIklan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
    SELECT COALESCE(SUM(biaya), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
    SELECT COALESCE(SUM(biaya), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT tanggal, SUM(biaya) AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanROAS(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(biaya) = 0 THEN 0
        ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float)
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(biaya) = 0 THEN 0
        ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float)
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT 
		tanggal,
        CASE 
            WHEN SUM(biaya) = 0 THEN 0
            ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float)
        END AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanConvertionRateIklan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(jumlah_klik) = 0 THEN 0
        ELSE (SUM(konversi)::float / SUM(jumlah_klik)::float) * 100
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(jumlah_klik) = 0 THEN 0
        ELSE (SUM(konversi)::float / SUM(jumlah_klik)::float) * 100
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT 
		tanggal,
        CASE 
            WHEN SUM(jumlah_klik) = 0 THEN 0
            ELSE (SUM(konversi)::float / SUM(jumlah_klik)::float) * 100
        END AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanPresentaseKlik(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(dilihat) = 0 THEN 0
        ELSE (SUM(jumlah_klik)::float / SUM(dilihat)::float) * 100
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
	SELECT
    CASE 
        WHEN SUM(dilihat) = 0 THEN 0
        ELSE (SUM(jumlah_klik)::float / SUM(dilihat)::float) * 100
    END AS cr
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT 
		tanggal,
        CASE 
            WHEN SUM(dilihat) = 0 THEN 0
            ELSE (SUM(jumlah_klik)::float / SUM(dilihat)::float) * 100
        END AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanDilihat(c *gin.Context) {
	var result dto.ResponseHeaderDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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

	var current, previous float64

	// total penjualan periode sekarang
	db.Raw(`
    SELECT COALESCE(SUM(dilihat), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, dateFrom, dateTo).Scan(&current)

	// cari range periode sebelumnya
	days := int(dateTo.Sub(dateFrom).Hours()/24) + 1
	periodBeforeFrom := dateFrom.AddDate(0, 0, -days)
	periodBeforeTo := dateTo.AddDate(0, 0, -days)

	// total periode sebelumnya
	db.Raw(`
    SELECT COALESCE(SUM(dilihat), 0)
    FROM shopee_data_upload_iklan_details
    WHERE store_id = ? AND tanggal BETWEEN ? AND ?
`, input.StoreId, periodBeforeFrom, periodBeforeTo).Scan(&previous)

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

	var sparkline []dto.ResponseSparkline
	// detail untuk sparkline
	db.Raw(`
	SELECT tanggal, SUM(dilihat) AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTotalPenjulandanBiaya(c *gin.Context) {
	var result []dto.ResponseTotalPenjualanDanBiayaDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}
	db := config.DB
	db.Raw(`
	WITH months AS (
		SELECT generate_series(
			date_trunc('year', ? :: date),
			date_trunc('year', ? :: date) + INTERVAL '11 months',
			INTERVAL '1 month'
		) AS month_start
	)
	SELECT 
		COALESCE(SUM(sd.biaya), 0) AS biaya,
		COALESCE(SUM(sd.omzet_penjualan), 0) AS penjualan,
		TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
	FROM months m
	LEFT JOIN shopee_data_upload_iklan_details sd 
		ON date_trunc('month', sd.tanggal) = m.month_start AND store_id = ?
	GROUP BY m.month_start
	ORDER BY m.month_start;
	`, input.DateFrom, input.DateTo, input.StoreId).Scan(&result)

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTotalROAS(c *gin.Context) {
	var result []dto.ResponseTotalROASDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}
	db := config.DB
	db.Raw(`
	WITH months AS (
		SELECT generate_series(
			date_trunc('year', ? :: date),
			date_trunc('year', ? :: date) + INTERVAL '11 months',
			INTERVAL '1 month'
		) AS month_start
	)
	SELECT 
		 CASE 
        	WHEN COALESCE(SUM(sd.biaya), 0) = 0 THEN 0
        	ELSE (COALESCE(SUM(sd.omzet_penjualan), 0)::float / COALESCE(SUM(sd.biaya), 0)::float)
    	END AS roas,
		TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
	FROM months m
	LEFT JOIN shopee_data_upload_iklan_details sd 
		ON date_trunc('month', sd.tanggal) = m.month_start AND store_id = ?
	GROUP BY m.month_start
	ORDER BY m.month_start;
	`, input.DateFrom, input.DateTo, input.StoreId).Scan(&result)

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTopProduct(c *gin.Context) {
	var result []dto.ResponseTopProductDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
	if errBind != nil {
		utils.Error(c, http.StatusBadRequest, "Invalid input", errBind.Error())
		return
	}
	db := config.DB
	db.Raw(`
	SELECT CONCAT('[', s.name, '] ', d.nama_iklan) AS nama_iklan, 
	CASE 
		WHEN SUM(biaya) = 0 THEN 0
		ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float)
	END AS roas,
	CASE 
		WHEN SUM(jumlah_klik) = 0 THEN 0
		ELSE (SUM(konversi)::float / SUM(jumlah_klik)::float) * 100
	END AS convertion_rate,
	COALESCE(SUM(biaya), 0) as biaya,
	COALESCE(SUM(omzet_penjualan), 0) as penjualan
	FROM shopee_data_upload_iklan_details d
	JOIN stores s ON s.id = d.store_id
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY concat('[', s.name, '] ', d.nama_iklan)
	ORDER BY penjualan DESC
	LIMIT 10;
	`, input.StoreId, input.DateFrom, input.DateTo).Scan(&result)

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}
