package models

import (
	"go-datalaris/common"
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeDataUploadChatDetail struct {
	common.BaseModel

	StoreID                           uint            `gorm:"column:store_id"`
	Tanggal                           time.Time       `gorm:"column:tanggal"`
	Pengunjung                        int             `gorm:"column:pengunjung"`
	JumlahChat                        int             `gorm:"column:jumlah_chat"`
	PengunjungBertanya                int             `gorm:"column:pengunjung_bertanya"`
	PertanyaanDiajukan                decimal.Decimal `gorm:"column:pertanyaan_diajukan;type:numeric(18,2)"`
	ChatDibalas                       int             `gorm:"column:chat_dibalas"`
	ChatBelumDibalas                  int             `gorm:"column:chat_belum_dibalas"`
	WaktuResponRataRata               time.Duration   `gorm:"column:waktu_respon_rata_rata"`
	Csat                              decimal.Decimal `gorm:"column:csat;type:numeric(18,2)"`
	WaktuResponChatPertama            time.Duration   `gorm:"column:waktu_respon_chat_pertama"`
	PresentaseChatDibalas             decimal.Decimal `gorm:"column:presentase_chat_dibalas;type:numeric(18,2)"`
	TingkatKonversiJumlahChatDirespon decimal.Decimal `gorm:"column:tingkat_konversi_jumlah_chat_direspon;type:numeric(18,2)"`
	TotalPembeli                      int             `gorm:"column:total_pembeli"`
	TotalPesanan                      int             `gorm:"column:total_pesanan"`
	Produk                            int             `gorm:"column:produk"`
	Penjualan                         decimal.Decimal `gorm:"column:penjualan;type:numeric(18,2)"`
	TingkatKonversiChatDibalas        decimal.Decimal `gorm:"column:tingkat_konversi_chat_dibalas;type:numeric(18,2)"`
}

func (ShopeeDataUploadChatDetail) TableName() string {
	return "shopee_data_upload_chat_details"
}
