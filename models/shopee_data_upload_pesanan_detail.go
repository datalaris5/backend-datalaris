package models

import (
	"go-datalaris/common"
	"time"

	"github.com/shopspring/decimal"
)

type ShopeeDataUploadPesananDetail struct {
	common.BaseModel

	StoreID                       uint            `gorm:"column:store_id"`
	NoPesanan                     string          `gorm:"column:no_pesanan;type:varchar(250)"`
	StatusPesanan                 string          `gorm:"column:status_pesanan;type:varchar(50)"`
	AlasanPembatalan              string          `gorm:"column:alasan_pembatalan;type:text"`
	StatusPembatalan              string          `gorm:"column:status_pembatalan;type:varchar(50)"`
	NoResi                        string          `gorm:"column:no_resi;type:varchar(250)"`
	OpsiPengiriman                string          `gorm:"column:opsi_pengiriman;type:varchar(250)"`
	CounterPickup                 string          `gorm:"column:counter_pickup;type:varchar(50)"`
	PesananHarusDikirimkanSebelum *time.Time      `gorm:"column:pesanan_harus_dikirimkan_sebelum"`
	WaktuPengirimanDiatur         *time.Time      `gorm:"column:waktu_pengiriman_diatur"`
	WaktuPesananDibuat            *time.Time      `gorm:"column:waktu_pesanan_dibuat"`
	WaktuPembayaranDilakukan      *time.Time      `gorm:"column:waktu_pembayaran_dilakukan"`
	MetodePembayaran              string          `gorm:"column:metode_pembayaran;type:varchar(100)"`
	SkuInduk                      string          `gorm:"column:sku_induk;type:varchar(100)"`
	NamaProduk                    string          `gorm:"column:nama_produk;type:text"`
	NomorReferensiSku             string          `gorm:"column:nomor_referensi_sku;type:varchar(250)"`
	NamaVariasi                   string          `gorm:"column:nama_variasi;type:varchar(250)"`
	HargaAwal                     decimal.Decimal `gorm:"column:harga_awal;type:numeric(18,2)"`
	HargaSetelahDiskon            decimal.Decimal `gorm:"column:harga_setelah_diskon;type:numeric(18,2)"`
	Jumlah                        int             `gorm:"column:jumlah"`
	ReturnedQuantity              int             `gorm:"column:returned_quantity"`
	TotalHargaProduk              decimal.Decimal `gorm:"column:total_harga_produk;type:numeric(18,2)"`
	TotalDiskon                   decimal.Decimal `gorm:"column:total_diskon;type:numeric(18,2)"`
	DiskonDariPenjual             decimal.Decimal `gorm:"column:diskon_dari_penjual;type:numeric(18,2)"`
	DiskonDariShopee              decimal.Decimal `gorm:"column:diskon_dari_shopee;type:numeric(18,2)"`
	BeratProduk                   string          `gorm:"column:berat_produk;type:varchar(50)"`
	JumlahProdukDipesan           int             `gorm:"column:jumlah_produk_dipesan"`
	TotalBerat                    string          `gorm:"column:total_berat;type:varchar(50)"`
	VoucherDitanggungPenjual      decimal.Decimal `gorm:"column:voucher_ditanggung_penjual;type:numeric(18,2)"`
	CashbackCoin                  decimal.Decimal `gorm:"column:cashback_coin;type:numeric(18,2)"`
	VoucherDitanggungShopee       decimal.Decimal `gorm:"column:voucher_ditanggung_shopee;type:numeric(18,2)"`
	PaketDiskon                   string          `gorm:"column:paket_diskon;type:varchar(10)"`
	PaketDiskonDariShopee         decimal.Decimal `gorm:"column:paket_diskon_dari_shopee;type:numeric(18,2)"`
	PaketDiskonDariPenjual        decimal.Decimal `gorm:"column:paket_diskon_dari_penjual;type:numeric(18,2)"`
	PotonganKoinShopee            decimal.Decimal `gorm:"column:potongan_koin_shopee;type:numeric(18,2)"`
	DiskonKartuKredit             decimal.Decimal `gorm:"column:diskon_kartu_kredit;type:numeric(18,2)"`
	OngkosKirimPembeli            decimal.Decimal `gorm:"column:ongkos_kirim_dibayar_oleh_pembeli;type:numeric(18,2)"`
	EstimasiPotonganOngkir        decimal.Decimal `gorm:"column:estimasi_potongan_biaya_pengiriman;type:numeric(18,2)"`
	OngkosKirimPengembalian       decimal.Decimal `gorm:"column:ongkos_kirim_pengembalian_barang;type:numeric(18,2)"`
	TotalPembayaran               decimal.Decimal `gorm:"column:total_pembayaran;type:numeric(18,2)"`
	PerkiraanOngkosKirim          decimal.Decimal `gorm:"column:perkiraan_ongkos_kirim;type:numeric(18,2)"`
	CatatanDariPembeli            string          `gorm:"column:catatan_dari_pembeli;type:text"`
	Catatan                       string          `gorm:"column:catatan;type:text"`
	UsernamePembeli               string          `gorm:"column:username_pembeli;type:varchar(250)"`
	NamaPenerima                  string          `gorm:"column:nama_penerima;type:varchar(250)"`
	NoTelepon                     string          `gorm:"column:no_telepon;type:varchar(50)"`
	AlamatPengiriman              string          `gorm:"column:alamat_pengiriman;type:varchar(250)"`
	Kota                          string          `gorm:"column:kota;type:varchar(250)"`
	Provinsi                      string          `gorm:"column:provinsi;type:varchar(250)"`
	WaktuPesananSelesai           *time.Time      `gorm:"column:waktu_pesanan_selesai"`
}

func (ShopeeDataUploadPesananDetail) TableName() string {
	return "shopee_data_upload_pesanan_details"
}
