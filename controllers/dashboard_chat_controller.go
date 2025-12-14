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

func GetDashboardChatJumlahChat(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		SELECT COALESCE(SUM(jumlah_chat), 0)
		FROM shopee_data_upload_chat_details
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
		SELECT tanggal, SUM(jumlah_chat) AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatChatDibalas(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		// total chat dibalas periode sekarang

		totalQuery := `
		SELECT COALESCE(SUM(chat_dibalas), 0)
    	FROM shopee_data_upload_chat_details
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
		SELECT tanggal, SUM(chat_dibalas) AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatTotalPembeli(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		// total pembeli periode sekarang
		totalQuery := `
		SELECT COALESCE(SUM(total_pembeli), 0)
    	FROM shopee_data_upload_chat_details
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
		SELECT tanggal, SUM(total_pembeli) AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatEstPenjualan(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		SELECT COALESCE(SUM(penjualan), 0)
    	FROM shopee_data_upload_chat_details
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
		SELECT tanggal, SUM(penjualan) AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatConvertionRate(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		SELECT
		CASE 
			WHEN COALESCE(SUM(chat_dibalas), 0) = 0 THEN 0
			ELSE (COALESCE(SUM(total_pembeli), 0)::float 
				/ COALESCE(SUM(chat_dibalas), 0)::float) * 100
		END AS current
		FROM shopee_data_upload_chat_details
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
        CASE 
            WHEN SUM(chat_dibalas) = 0 THEN 0
        	ELSE (SUM(total_pembeli)::float / SUM(chat_dibalas)::float) * 100
        END AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatPersentaseChat(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		CASE 
			WHEN COALESCE(SUM(jumlah_chat), 0) = 0 THEN 0
			ELSE (COALESCE(SUM(chat_dibalas), 0)::float 
				/ COALESCE(SUM(jumlah_chat), 0)::float) * 100
		END AS current
		FROM shopee_data_upload_chat_details
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
        CASE 
            WHEN SUM(jumlah_chat) = 0 THEN 0
        	ELSE (SUM(chat_dibalas)::float / SUM(jumlah_chat)::float) * 100
        END AS total
		FROM shopee_data_upload_chat_details
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardRataRataWaktuRespon(c *gin.Context) {
	var result dto.ResponseHeaderDashboardChatTimeDuration
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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

	var current, previous string

	if marketplace.Name == constant.ShopeeConst {
		// total penjualan periode sekarang
		totalQuery := `
		SELECT
		COALESCE(
			SUM(
				COALESCE(chat_dibalas, 0) * COALESCE(waktu_respon_rata_rata, INTERVAL '0')
			),
			INTERVAL '0'
		) AS current
		FROM shopee_data_upload_chat_details
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

		previousArgs := []interface{}{
			periodBeforeFrom,
			periodBeforeTo,
		}

		if input.StoreId != 0 {
			previousArgs = append(previousArgs, input.StoreId)
		}

		// total periode sebelumnya
		db.Raw(totalQuery, previousArgs...).Scan(&previous)

		currentDur, _ := utils.ParseIntervalToDuration(current)

		previousDur, _ := utils.ParseIntervalToDuration(previous)

		// hitung persentase naik / turun
		currentSec := currentDur.Seconds()
		previousSec := previousDur.Seconds()

		changePercent := 0.0
		if previousSec > 0 {
			changePercent = ((currentSec - previousSec) / previousSec) * 100
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
			SUM(
				COALESCE(chat_dibalas, 0) 
				* COALESCE(waktu_respon_rata_rata, INTERVAL '0')
			),
			INTERVAL '0'
    	) AS total
		FROM shopee_data_upload_chat_details
		WHERE tanggal BETWEEN ? AND ?`

		if input.StoreId != 0 {
			sparklineQuery += " AND store_id = ?"
		}

		sparklineQuery += " GROUP BY tanggal ORDER BY tanggal ASC"

		var sparkline []dto.ResponseSparklineTimeDuration
		// detail untuk sparkline
		db.Raw(sparklineQuery, currentArgs...).Scan(&sparkline)
		result.Sparkline = sparkline
	}

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatTotalJumlahChat(c *gin.Context) {
	var result []dto.ResponseTotalJumlahChatDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
			COALESCE(SUM(sd.jumlah_chat), 0) AS jumlah_chat,
			TRIM(TO_CHAR(m.month_start, 'FMMonth YYYY')) AS month
		FROM months m
		LEFT JOIN shopee_data_upload_chat_details sd 
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

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}

func GetDashboardChatRataRataWaktuResponInWeek(c *gin.Context) {
	var result []dto.ResponseAvgWaktuResponInWeekDashboardChat
	input, errBind := utils.BindJSON[dto.RequestDashboardChat](c)
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
		WITH days AS (
			SELECT generate_series(
				date_trunc('week', ?::date),
				date_trunc('week', ?::date) + INTERVAL '6 days',
				INTERVAL '1 day'
			)::date AS tanggal
		)
		SELECT
			TRIM(TO_CHAR(d.tanggal, 'Day')) AS day,
			COALESCE(
				SUM(
					COALESCE(s.chat_dibalas, 0)
					* COALESCE(s.waktu_respon_rata_rata, INTERVAL '0')
				),
				INTERVAL '0'
			) AS total
		FROM days d
		LEFT JOIN shopee_data_upload_chat_details s
			ON s.tanggal::date = d.tanggal
		`

		args := []interface{}{
			input.DateFrom,
			input.DateFrom,
		}

		if input.StoreId != 0 {
			query += " AND s.store_id = ?"
			args = append(args, input.StoreId)
		}
		query += " GROUP BY d.tanggal ORDER BY d.tanggal"

		db.Raw(query, args...).Scan(&result)

		for i := range result {
			result[i].Day = utils.ChangeDayEnIn(result[i].Day)
		}
	}

	utils.Success(c, constant.DashboardChatConst+constant.SuccessFetch, result)
}
