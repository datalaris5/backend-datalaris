package models

import (
	"go-datalaris/common"
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeDataUploadDetail struct {
	common.BaseModel
	StoreID                  uint            `gorm:"column:store_id"`
	Tanggal                  time.Time       `gorm:"column:tanggal"`
	TotalPenjualan           decimal.Decimal `gorm:"column:total_penjualan;type:numeric(18,2)"`
	TotalPesanan             int             `gorm:"column:total_pesanan"`
	PenjualanPerPesanan      decimal.Decimal `gorm:"column:penjualan_per_pesanan;type:numeric(18,2)"`
	ProdukKlik               int             `gorm:"column:produk_klik"`
	TotalPengunjung          int             `gorm:"column:total_pengunjung"`
	TingkatKonversiHarian    decimal.Decimal `gorm:"column:tingkat_konversi_harian;type:numeric(5,2)"`
	PesananDibatalkan        int             `gorm:"column:pesanan_dibatalkan"`
	PenjualanDibatalkan      int             `gorm:"column:penjualan_dibatalkan"`
	PesananDikembalikan      int             `gorm:"column:pesanan_dikembalikan"`
	PenjualanDikembalikan    int             `gorm:"column:penjualan_dikembalikan"`
	Pembeli                  int             `gorm:"column:pembeli"`
	TotalPembeliBaru         int             `gorm:"column:total_pembeli_baru"`
	TotalPembeliSaatIni      int             `gorm:"column:total_pembeli_saat_ini"`
	TotalPotensiPembeli      int             `gorm:"column:total_potensi_pembeli"`
	TingkatPembelianBerulang decimal.Decimal `gorm:"column:tingkat_pembelian_berulang;type:numeric(5,2)"`
}

func (ShopeeDataUploadDetail) TableName() string {
	return "shopee_data_upload_details"
}
