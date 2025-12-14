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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(omzet_penjualan), 0)
    	FROM shopee_data_upload_iklan_details
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
		SELECT tanggal, SUM(omzet_penjualan) AS total
		FROM shopee_data_upload_iklan_details
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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(biaya), 0)
    	FROM shopee_data_upload_iklan_details
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
		SELECT tanggal, SUM(biaya) AS total
		FROM shopee_data_upload_iklan_details
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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT
		COALESCE(
			SUM(omzet_penjualan)::float
			/ NULLIF(SUM(biaya), 0)::float,
			0
    	) AS cr
		FROM shopee_data_upload_iklan_details
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
		SELECT 
		tanggal,
        COALESCE(
			SUM(omzet_penjualan)::float
			/ NULLIF(SUM(biaya), 0)::float,
			0
    	) AS total
		FROM shopee_data_upload_iklan_details
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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang

		totalQuery := `
		SELECT
		COALESCE(
			SUM(konversi)::float
			/ NULLIF(SUM(jumlah_klik), 0)::float,
			0
    	) AS cr
		FROM shopee_data_upload_iklan_details
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
		SELECT 
		tanggal,
        COALESCE(
			SUM(konversi)::float
			/ NULLIF(SUM(jumlah_klik), 0)::float,
			0
    	) AS total
		FROM shopee_data_upload_iklan_details
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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT
		COALESCE(
			SUM(jumlah_klik)::float
			/ NULLIF(SUM(dilihat), 0)::float,
			0
    	) AS cr
		FROM shopee_data_upload_iklan_details
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
		SELECT 
		tanggal,
        COALESCE(
			SUM(jumlah_klik)::float
			/ NULLIF(SUM(dilihat), 0)::float,
			0
    	) AS total
		FROM shopee_data_upload_iklan_details
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

	marketplace, err := services.GetWhereFirst[models.Marketplace]("id = ?", input.MarketplaceId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.MarketplaceConst+constant.ErrorNotFound, nil)
		return
	}

	var current, previous float64

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(dilihat), 0)
    	FROM shopee_data_upload_iklan_details
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
		SELECT tanggal, SUM(dilihat) AS total
		FROM shopee_data_upload_iklan_details
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

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTotalPenjulandanBiaya(c *gin.Context) {
	var result []dto.ResponseTotalPenjualanDanBiayaDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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
			COALESCE(SUM(sd.biaya), 0) AS biaya,
			COALESCE(SUM(sd.omzet_penjualan), 0) AS penjualan,
			TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
		FROM months m
		LEFT JOIN shopee_data_upload_iklan_details sd 
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

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTotalROAS(c *gin.Context) {
	var result []dto.ResponseTotalROASDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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
			COALESCE(
				SUM(omzet_penjualan)::float
				/ NULLIF(SUM(biaya), 0)::float,
				0
    		) AS roas,
			TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
		FROM months m
		LEFT JOIN shopee_data_upload_iklan_details sd 
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

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}

func GetDashboardTopProduct(c *gin.Context) {
	var result []dto.ResponseTopProductDashboardIklan
	input, errBind := utils.BindJSON[dto.RequestDashboardIklan](c)
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
		SELECT CONCAT('[', s.name, '] ', d.nama_iklan) AS nama_iklan, 
		COALESCE(
			SUM(omzet_penjualan)::float
			/ NULLIF(SUM(biaya), 0)::float,
			0
    	) AS roas,
		COALESCE(
			SUM(konversi)::float
			/ NULLIF(SUM(jumlah_klik), 0)::float,
			0
    	) AS convertion_rate,
		COALESCE(SUM(biaya), 0) as biaya,
		COALESCE(SUM(omzet_penjualan), 0) as penjualan
		FROM shopee_data_upload_iklan_details d
		JOIN stores s ON s.id = d.store_id
		WHERE tanggal BETWEEN ? AND ?
		`
		args := []interface{}{
			input.DateFrom,
			input.DateTo,
		}

		if input.StoreId != 0 {
			query += " AND store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY concat('[', s.name, '] ', d.nama_iklan)"

		db.Raw(query, args...).Scan(&result)
	}

	utils.Success(c, constant.DashboardIklanConst+constant.SuccessFetch, result)
}
