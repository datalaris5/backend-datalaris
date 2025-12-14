package controllers

import (
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/dto"
	"go-datalaris/models"
	"go-datalaris/services"
	"go-datalaris/utils"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetDashboardTinjauanTotalPenjualan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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
		SELECT COALESCE(SUM(total_penjualan), 0)
    	FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?
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
		SELECT tanggal, SUM(total_penjualan) AS total
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanTotalPesanan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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

	var current, previous int64

	if marketplace.Name == constant.ShopeeConst {
		// total pesanan periode sekarang

		totalQuery := `
		SELECT COALESCE(SUM(total_pesanan), 0)
    	FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?
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
			changePercent = ((float64(current) - float64(previous)) / float64(previous)) * 100
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
		SELECT tanggal, SUM(total_pesanan) AS total
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanTotalPengunjung(c *gin.Context) {
	var result dto.ResponseHeaderDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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

	var current, previous int64

	if marketplace.Name == constant.ShopeeConst {
		// total pengunjung periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(total_pengunjung), 0)
    	FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?
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
			changePercent = ((float64(current) - float64(previous)) / float64(previous)) * 100
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
		SELECT tanggal, SUM(total_pengunjung) AS total
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanTrenPenjualan(c *gin.Context) {
	var result []dto.ResponseTrenPenjualanDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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
		WITH months AS (
			SELECT generate_series(
				date_trunc('year', ? :: date),
				date_trunc('year', ? :: date) + INTERVAL '11 months',
				INTERVAL '1 month'
			) AS month_start
		)
		SELECT 
			COALESCE(SUM(sd.total_penjualan), 0) AS total,
			TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
		FROM months m
		LEFT JOIN shopee_data_upload_details sd 
			ON date_trunc('month', sd.tanggal) = m.month_start
		`
		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY m.month_start ORDER BY m.month_start"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanTotalPesananInWeek(c *gin.Context) {
	var result []dto.ResponseTotalPenjualanInWeekDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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
		TRIM(TO_CHAR(tanggal, 'Day')) AS day,
		COALESCE(SUM(total_pesanan), 0) AS total
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? ::date - INTERVAL '6 days' AND ? ::date
		`
		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY day, tanggal ORDER BY tanggal"

		db.Raw(query, args...).Scan(&result)
		for i := range result {
			result[i].Day = utils.ChangeDayEnIn(result[i].Day)
		}
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanYoYGrowth(c *gin.Context) {

}

func GetDashboardTinjauanConvertionRate(c *gin.Context) {
	var result dto.ResponseHeaderDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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
		// total convertion rate periode sekarang
		totalQuery := `
		SELECT 
		CASE 
			WHEN SUM(total_pengunjung) = 0 THEN 0
			ELSE (SUM(pembeli)::float / SUM(total_pengunjung)::float) * 100
		END AS cr
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?
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

		changePercent := 0.0

		if previous > 0 {
			changePercent = ((current - previous) / previous) * 100
		}

		if changePercent > 0 {
			result.Trend = "Up"
		} else if changePercent < 0 {
			result.Trend = "Down"
		} else {
			result.Trend = "Equal"
		}

		result.Total = current
		result.Percent = math.Round(changePercent*100) / 100 // dua angka desimal

		sparklineQuery := `
		SELECT 
        tanggal,
        CASE 
            WHEN SUM(total_pengunjung) = 0 THEN 0
            ELSE (SUM(pembeli)::float / SUM(total_pengunjung)::float) * 100
        END AS total
    	FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}

func GetDashboardTinjauanBasketSize(c *gin.Context) {
	var result dto.ResponseHeaderDashboardTinjauan
	input, errBind := utils.BindJSON[dto.RequestDashboardTinjauan](c)
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
		// total penjualan periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(total_penjualan)::float/SUM(total_pesanan)::float, 0)
    	FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?
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
		SELECT tanggal, (SUM(total_penjualan)::float/SUM(total_pesanan)::float) AS total
		FROM shopee_data_upload_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparkline
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardTinjauanConst+constant.SuccessFetch, result)
}
