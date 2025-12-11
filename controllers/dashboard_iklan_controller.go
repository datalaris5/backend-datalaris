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
        ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float) * 100
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
        ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float) * 100
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
            ELSE (SUM(omzet_penjualan)::float / SUM(biaya)::float) * 100
        END AS total
	FROM shopee_data_upload_iklan_details
	WHERE store_id = ? AND tanggal BETWEEN ? AND ?
	GROUP BY tanggal
	ORDER BY tanggal ASC;`, input.StoreId, dateFrom, dateTo).Scan(&sparkline)

	result.Sparkline = sparkline

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardIklanConvertionRateIklan(c *gin.Context) {

}

func GetDashboardIklanPresentaseKlik(c *gin.Context) {

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
