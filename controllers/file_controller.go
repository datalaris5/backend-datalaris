package controllers

import (
	"context"
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

func UploadExcel(c *gin.Context) {
	tenantID := utils.GetTenantId(c)
	// c.FormFile hanya mengembalikan 2 nilai: *multipart.FileHeader, error
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, constant.ErrorUploadFile, nil)
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

	summaries, details := readExcelShopee(dstPath)

	userID, ok := c.Get("user_id")
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "User ID not found", nil)
		return
	}

	ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)

	db := config.DB
	err = services.WithTransaction(db.WithContext(ctx), func(tx *gorm.DB) error {
		summary, err := ConvertSummaryFromRow(summaries, tenantID)
		_, err = services.Save[models.ShopeeDataUploadSummary](*summary, tx)
		if err != nil {
			return err
		}

		for i := range details {
			detail, err := ConvertDetailFromRow(details[i], tenantID)
			_, err = services.Save[models.ShopeeDataUploadDetail](*detail, tx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Transaction failed", utils.ParseDBError(err))
		return
	}

	utils.Success(c, constant.SuccessUploadFile, nil)
}

func readExcelShopee(path string) ([]string, [][]string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println("Error baca excel:", err)
		return nil, nil
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
		return nil, nil
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		fmt.Println("Gagal baca rows:", err)
		return nil, nil
	}

	var summaries []string
	var details [][]string

	for i := range rows {
		if i == 0 || i == 2 || i == 3 {
			continue
		} else if i == 1 {
			summaries = rows[i]
		} else {
			details = append(details, rows[i])
		}
	}

	return summaries, details

}

func ConvertSummaryFromRow(row []string, tenantID any) (*models.ShopeeDataUploadSummary, error) {
	// pastikan minimal ada kolom
	if len(row) < 16 {
		return nil, fmt.Errorf("jumlah kolom summary tidak sesuai")
	}

	return &models.ShopeeDataUploadSummary{
		TenantID:                 tenantID.(uint),
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

func ConvertDetailFromRow(row []string, tenantID any) (*models.ShopeeDataUploadDetail, error) {
	// pastikan minimal ada kolom
	if len(row) < 16 {
		return nil, fmt.Errorf("jumlah kolom summary tidak sesuai")
	}

	// Parse tanggal
	tanggal, err := parseTanggalShopee(row[0]) // contoh format Shopee dd/mm/yyyy
	if err != nil {
		return nil, err
	}

	return &models.ShopeeDataUploadDetail{
		TenantID:                 tenantID.(uint),
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

func toInt(s string) int {
	i, _ := strconv.Atoi(strings.ReplaceAll(s, ",", ""))
	return i
}

func CleanNumber(s string) string {
	s = strings.TrimSpace(s)

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
