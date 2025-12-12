package controllers

import (
	"context"
	"encoding/csv"
	"fmt"
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/models"
	"go-datalaris/services"
	"go-datalaris/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

func UploadExcelShopeeTinjauan(c *gin.Context) {
	storeId := c.Param("id")
	// c.FormFile hanya mengembalikan 2 nilai: *multipart.FileHeader, error
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, constant.ErrorUploadFile, nil)
		return
	}

	_, err = services.GetWhereFirst[models.Store]("id = ?", storeId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.StoreConst+constant.ErrorNotFound, nil)
		return
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		utils.Error(c, http.StatusBadRequest, constant.ErrorInvalidFileType, nil)
		return
	}

	// Simpan file ke folder uploads
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCreateDir, nil)
		return
	}

	dstPath := filepath.Join("uploads", fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorSaveFile, nil)
		return
	}
	defer dst.Close()

	// Buka file dari FileHeader
	src, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorOpenFile, nil)
		return
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCopyFile, nil)
		return
	}

	details := readExcelShopeeTinjauan(dstPath)

	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err = services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {

		startDate, err := parseTanggalShopee(details[0][0])
		if err != nil {
			return fmt.Errorf("invalid first date: %w", err)
		}
		endDate, err := parseTanggalShopee(details[len(details)-1][0])
		if err != nil {
			return fmt.Errorf("invalid last date: %w", err)
		}

		// 1. Delete data existing untuk rentang tanggal upload
		_, err = tx.
			Where("store_id = ? AND tanggal BETWEEN ? AND ?", storeId, startDate, endDate).
			Delete(&models.ShopeeDataUploadTinjauanDetail{}).
			Rows()
		if err != nil {
			return err
		}

		// 2. Insert ulang row baru
		for i := range details {
			detail, err := ConvertDetailShopeeTinjauanFromRow(details[i], storeId)
			if err != nil {
				return err
			}

			_, err = services.Save[models.ShopeeDataUploadTinjauanDetail](*detail, tx)
			if err != nil {
				return err
			}
		}

		// 3. Insert history log
		history := models.HistoryDataUpload{
			StoreId:  utils.ParseUintParam(storeId),
			Filename: fileHeader.Filename,
			Status:   constant.Success,
		}
		_, err = services.Save[models.HistoryDataUpload](history, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.SuccessUploadFile, nil)
}

func UploadCsvShopeeIklan(c *gin.Context) {
	storeId := c.Param("id")
	// c.FormFile hanya mengembalikan 2 nilai: *multipart.FileHeader, error
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, constant.ErrorUploadFile, nil)
		return
	}

	_, err = services.GetWhereFirst[models.Store]("id = ?", storeId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.StoreConst+constant.ErrorNotFound, nil)
		return
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fmt.Println(ext)
	if ext != ".csv" {
		utils.Error(c, http.StatusBadRequest, constant.ErrorInvalidFileType, nil)
		return
	}

	// Simpan file ke folder uploads
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCreateDir, nil)
		return
	}

	dstPath := filepath.Join("uploads", fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorSaveFile, nil)
		return
	}
	defer dst.Close()

	// Buka file dari FileHeader
	src, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorOpenFile, nil)
		return
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCopyFile, nil)
		return
	}

	details, dateFrom, dateTo := readCsvShopeeIklan(dstPath)

	if !dateFrom.Equal(dateTo) {
		utils.Error(c, http.StatusBadRequest, "Range tanggal tidak valid, harus 1 hari", nil)
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err = services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {

		// 1. Delete data existing untuk tanggal upload
		_, err = tx.
			Where("store_id = ? AND tanggal = ?", storeId, dateFrom).
			Delete(&models.ShopeeDataUploadIklanDetail{}).
			Rows()
		if err != nil {
			return err
		}

		// 2. Insert ulang row baru
		for i := range details {
			detail, err := ConvertDetailShopeeIklanFromRow(details[i], storeId, dateFrom)
			if err != nil {
				return err
			}

			_, err = services.Save[models.ShopeeDataUploadIklanDetail](*detail, tx)
			if err != nil {
				return err
			}
		}

		// 3. Insert history log
		history := models.HistoryDataUpload{
			StoreId:  utils.ParseUintParam(storeId),
			Filename: fileHeader.Filename,
			Status:   constant.Success,
		}
		_, err = services.Save[models.HistoryDataUpload](history, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.SuccessUploadFile, nil)
}

func UploadExcelShopeeChat(c *gin.Context) {
	storeId := c.Param("id")
	// c.FormFile hanya mengembalikan 2 nilai: *multipart.FileHeader, error
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, constant.ErrorUploadFile, nil)
		return
	}

	_, err = services.GetWhereFirst[models.Store]("id = ?", storeId)
	if err == gorm.ErrRecordNotFound {
		utils.Error(c, http.StatusNotFound, constant.StoreConst+constant.ErrorNotFound, nil)
		return
	}

	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		utils.Error(c, http.StatusBadRequest, constant.ErrorInvalidFileType, nil)
		return
	}

	// Simpan file ke folder uploads
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCreateDir, nil)
		return
	}

	dstPath := filepath.Join("uploads", fileHeader.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorSaveFile, nil)
		return
	}
	defer dst.Close()

	// Buka file dari FileHeader
	src, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorOpenFile, nil)
		return
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, constant.ErrorCopyFile, nil)
		return
	}

	details := readExcelShopeeChat(dstPath)

	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err = services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {

		startDate, err := parseTanggalShopee(details[0][0])
		if err != nil {
			return fmt.Errorf("invalid first date: %w", err)
		}
		endDate, err := parseTanggalShopee(details[len(details)-1][0])
		if err != nil {
			return fmt.Errorf("invalid last date: %w", err)
		}

		// 1. Delete data existing untuk rentang tanggal upload
		_, err = tx.
			Where("store_id = ? AND tanggal BETWEEN ? AND ?", storeId, startDate, endDate).
			Delete(&models.ShopeeDataUploadChatDetail{}).
			Rows()
		if err != nil {
			return err
		}

		// 2. Insert ulang row baru
		for i := range details {
			detail, err := ConvertDetailShopeeChatFromRow(details[i], storeId)
			if err != nil {
				return err
			}

			_, err = services.Save[models.ShopeeDataUploadChatDetail](*detail, tx)
			if err != nil {
				return err
			}
		}

		// 3. Insert history log
		history := models.HistoryDataUpload{
			StoreId:  utils.ParseUintParam(storeId),
			Filename: fileHeader.Filename,
			Status:   constant.Success,
		}
		_, err = services.Save[models.HistoryDataUpload](history, tx)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.SuccessUploadFile, nil)
}

func readExcelShopeeTinjauan(path string) [][]string {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println("Error baca excel:", err)
		return nil
	}

	sheet := "Pesanan Siap Dikirim"

	// Cek apakah sheet ada
	sheets := f.GetSheetList()
	found := false
	for _, s := range sheets {
		if s == sheet {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Sheet", sheet, "tidak ditemukan")
		return nil
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println("Gagal baca rows:", err)
		return nil
	}

	var details [][]string

	for i := range rows {
		if i < 4 {
			continue
		} else {
			details = append(details, rows[i])
		}
	}

	return details

}

func ConvertDetailShopeeTinjauanFromRow(row []string, storeId string) (*models.ShopeeDataUploadTinjauanDetail, error) {
	// pastikan minimal ada kolom
	if len(row) < 16 {
		return nil, fmt.Errorf("jumlah kolom summary tidak sesuai")
	}

	// Parse tanggal
	tanggal, err := parseTanggalShopee(row[0]) // contoh format Shopee dd/mm/yyyy
	if err != nil {
		return nil, err
	}

	return &models.ShopeeDataUploadTinjauanDetail{
		StoreID:                  utils.ParseUintParam(storeId),
		Tanggal:                  tanggal,
		TotalPenjualan:           decimal.RequireFromString(CleanNumber(row[1])),
		TotalPesanan:             toInt(row[2]),
		PenjualanPerPesanan:      decimal.RequireFromString(CleanNumber(row[3])),
		ProdukKlik:               toInt(row[4]),
		TotalPengunjung:          toInt(row[5]),
		TingkatKonversiHarian:    decimal.RequireFromString(CleanNumber(row[6])),
		PesananDibatalkan:        toInt(row[7]),
		PenjualanDibatalkan:      toInt(row[8]),
		PesananDikembalikan:      toInt(row[9]),
		PenjualanDikembalikan:    toInt(row[10]),
		Pembeli:                  toInt(row[11]),
		TotalPembeliBaru:         toInt(row[12]),
		TotalPembeliSaatIni:      toInt(row[13]),
		TotalPotensiPembeli:      toInt(row[14]),
		TingkatPembelianBerulang: decimal.RequireFromString(CleanNumber(row[15])),
	}, nil
}

func readCsvShopeeIklan(path string) ([][]string, time.Time, time.Time) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error baca csv:", err)
		return nil, time.Time{}, time.Time{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Gagal baca csv:", err)
		return nil, time.Time{}, time.Time{}
	}

	var details [][]string
	var dateFrom, dateTo time.Time

	for i := range rows {
		if i == 5 {
			raw := rows[i][1] // kolom ke-1
			fmt.Println(raw)
			dateFrom, dateTo, _ = parseDateRange(raw)
			continue
		}

		// skip 7 baris header
		if i < 7 {
			continue
		}
		details = append(details, rows[i])
	}

	return details, dateFrom, dateTo
}

func ConvertDetailShopeeIklanFromRow(row []string, storeId string, tanggal time.Time) (*models.ShopeeDataUploadIklanDetail, error) {
	// minimal 28 kolom
	if len(row) < 28 {
		return nil, fmt.Errorf("jumlah kolom iklan tidak sesuai")
	}

	return &models.ShopeeDataUploadIklanDetail{
		StoreID:                  utils.ParseUintParam(storeId),
		Tanggal:                  tanggal,
		NamaIklan:                row[1],
		Status:                   row[2],
		JenisIklan:               row[3],
		KodeProduk:               row[4],
		TampilanIklan:            row[5],
		ModeBidding:              row[6],
		PenempatanIklan:          row[7],
		TanggalMulai:             row[8],
		TanggalSelesai:           row[9],
		Dilihat:                  toInt(row[10]),
		JumlahKlik:               toInt(row[11]),
		PresentaseKlik:           decimal.RequireFromString(CleanNumber(row[12])),
		Konversi:                 toInt(row[13]),
		KonversiLangsung:         toInt(row[14]),
		TingkatKonversi:          decimal.RequireFromString(CleanNumber(row[15])),
		BiayaPerKonversi:         decimal.RequireFromString(CleanNumber(row[16])),
		BiayaPerKonversiLangsung: decimal.RequireFromString(CleanNumber(row[17])),
		ProdukTerjual:            toInt(row[18]),
		TerjualLangsung:          toInt(row[19]),
		OmzetPenjualan:           decimal.RequireFromString(CleanNumber(row[20])),
		GVMLangsung:              decimal.RequireFromString(CleanNumber(row[21])),
		Biaya:                    decimal.RequireFromString(CleanNumber(row[22])),
		EfektifitasIklan:         decimal.RequireFromString(CleanNumber(row[23])),
		EfektifitasLangsung:      decimal.RequireFromString(CleanNumber(row[24])),
		ACOS:                     decimal.RequireFromString(CleanNumber(row[25])),
		ACOSLangsung:             decimal.RequireFromString(CleanNumber(row[26])),
		JumlahProdukDilihat:      toInt(row[27]),
		JumlahProdukDiklik:       toInt(row[28]),
		PresentaseProdukDiklik:   decimal.RequireFromString(CleanNumber(row[29])),
	}, nil
}

func readExcelShopeeChat(path string) [][]string {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println("Error baca excel:", err)
		return nil
	}

	sheet := "Grafik Kriteria"

	// Cek apakah sheet ada
	sheets := f.GetSheetList()
	found := false
	for _, s := range sheets {
		if s == sheet {
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Sheet", sheet, "tidak ditemukan")
		return nil
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println("Gagal baca rows:", err)
		return nil
	}

	var details [][]string

	for i := range rows {
		if i < 1 {
			continue
		} else {
			details = append(details, rows[i])
		}
	}

	return details

}

func ConvertDetailShopeeChatFromRow(row []string, storeId string) (*models.ShopeeDataUploadChatDetail, error) {

	// jumlah kolom minimum
	if len(row) < 17 {
		return nil, fmt.Errorf("jumlah kolom chat detail tidak sesuai")
	}

	// Parse tanggal (format Shopee: dd/mm/yyyy)
	tanggal, err := parseTanggalShopee(row[0])
	if err != nil {
		return nil, err
	}

	waktuResponRata, err := ParseIntervalToDuration(row[7])
	if err != nil {
		return nil, err
	}

	waktuResponPertama, err := ParseIntervalToDuration(row[9])
	if err != nil {
		return nil, err
	}

	return &models.ShopeeDataUploadChatDetail{
		StoreID:                           utils.ParseUintParam(storeId),
		Tanggal:                           tanggal,
		Pengunjung:                        toInt(row[1]),
		JumlahChat:                        toInt(row[2]),
		PengunjungBertanya:                toInt(row[3]),
		PertanyaanDiajukan:                decimal.RequireFromString(CleanNumber(row[4])),
		ChatDibalas:                       toInt(row[5]),
		ChatBelumDibalas:                  toInt(row[6]),
		WaktuResponRataRata:               waktuResponRata,
		Csat:                              decimal.RequireFromString(CleanNumber(row[8])),
		WaktuResponChatPertama:            waktuResponPertama,
		PresentaseChatDibalas:             decimal.RequireFromString(CleanNumber(row[10])),
		TingkatKonversiJumlahChatDirespon: decimal.RequireFromString(CleanNumber(row[11])),
		TotalPembeli:                      toInt(row[12]),
		TotalPesanan:                      toInt(row[13]),
		Produk:                            toInt(row[14]),
		Penjualan:                         decimal.RequireFromString(CleanNumber(row[15])),
		TingkatKonversiChatDibalas:        decimal.RequireFromString(CleanNumber(row[16])),
	}, nil
}

func parseDateRange(raw string) (time.Time, time.Time, error) {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("format tidak valid")
	}

	dateFromStr := strings.TrimSpace(parts[0])
	dateToStr := strings.TrimSpace(parts[1])

	layout := "02/01/2006"

	dateFrom, err1 := time.Parse(layout, dateFromStr)
	dateTo, err2 := time.Parse(layout, dateToStr)

	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("gagal parse tanggal")
	}

	return dateFrom, dateTo, nil
}

func toInt(s string) int {
	i, _ := strconv.Atoi(strings.ReplaceAll(s, ",", ""))
	return i
}

func CleanNumber(s string) string {
	s = strings.TrimSpace(s)

	// kalau "-"
	if s == "-" || s == "" {
		return "0"
	}

	if strings.Contains(s, "%") {
		s = strings.ReplaceAll(s, "%", "")
		s = strings.ReplaceAll(s, ",", ".")
		return s
	}

	// normal number
	s = strings.ReplaceAll(s, ".", "") // ribuan
	s = strings.ReplaceAll(s, ",", ".")
	return s
}

func parseTanggalShopee(s string) (time.Time, error) {
	formats := []string{
		"02/01/2006",
		"2006-01-02",
		"02-01-2006",
		"02/01/2006 15:04",
		"2006/01/02",
	}

	for _, f := range formats {
		if t, err := time.Parse(f, strings.TrimSpace(s)); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("format tanggal tidak didukung: %s", s)
}

func ParseIntervalToDuration(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, nil
	}

	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("format interval tidak valid: %s", s)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("interval jam tidak valid: %s", s)
	}

	mins, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("interval menit tidak valid: %s", s)
	}

	secs, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf("interval detik tidak valid: %s", s)
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(mins)*time.Minute +
		time.Duration(secs)*time.Second, nil
}
