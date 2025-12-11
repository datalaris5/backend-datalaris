package models

import (
	"go-datalaris/common"
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeDataUploadIklanDetail struct {
	common.BaseModel
	StoreID                  uint            `gorm:"column:store_id"`
	Tanggal                  time.Time       `gorm:"column:tanggal"`
	NamaIklan                string          `gorm:"column:nama_iklan"`
	Status                   string          `gorm:"column:status"`
	JenisIklan               string          `gorm:"column:jenis_iklan"`
	KodeProduk               string          `gorm:"column:kode_produk"`
	TampilanIklan            string          `gorm:"column:tampilan_iklan"`
	ModeBidding              string          `gorm:"column:mode_bidding"`
	PenempatanIklan          string          `gorm:"column:penempatan_iklan"`
	TanggalMulai             string          `gorm:"column:tanggal_mulai"`
	TanggalSelesai           string          `gorm:"column:tanggal_selesai"`
	Dilihat                  int             `gorm:"column:dilihat"`
	JumlahKlik               int             `gorm:"column:jumlah_klik"`
	PresentaseKlik           decimal.Decimal `gorm:"column:presentase_klik;type:numeric(5,2)"`
	Konversi                 int             `gorm:"column:konversi"`
	KonversiLangsung         int             `gorm:"column:konversi_langsung"`
	TingkatKonversi          decimal.Decimal `gorm:"column:tingkat_konversi;type:numeric(5,2)"`
	BiayaPerKonversi         decimal.Decimal `gorm:"column:biaya_per_konversi;type:numeric(18,2)"`
	BiayaPerKonversiLangsung decimal.Decimal `gorm:"column:biaya_per_konversi_langsung;type:numeric(18,2)"`
	ProdukTerjual            int             `gorm:"column:produk_terjual"`
	TerjualLangsung          int             `gorm:"column:terjual_langsung"`
	OmzetPenjualan           decimal.Decimal `gorm:"column:omzet_penjualan;type:numeric(18,2)"`
	GVMLangsung              decimal.Decimal `gorm:"column:gvm_langsung;type:numeric(18,2)"`
	Biaya                    decimal.Decimal `gorm:"column:biaya;type:numeric(18,2)"`
	EfektifitasIklan         decimal.Decimal `gorm:"column:efektifitas_iklan;type:numeric(5,2)"`
	EfektifitasLangsung      decimal.Decimal `gorm:"column:efektifitas_langsung;type:numeric(5,2)"`
	ACOS                     decimal.Decimal `gorm:"column:acos;type:numeric(5,2)"`
	ACOSLangsung             decimal.Decimal `gorm:"column:acos_langsung;type:numeric(5,2)"`
	JumlahProdukDilihat      int             `gorm:"column:jumlah_produk_dilihat"`
	JumlahProdukDiklik       int             `gorm:"column:jumlah_produk_diklik"`
	PresentaseProdukDiklik   decimal.Decimal `gorm:"column:presentase_produk_diklik;type:numeric(5,2)"`
}

func (ShopeeDataUploadIklanDetail) TableName() string {
	return "shopee_data_upload_iklan_details"
}
