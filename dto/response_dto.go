package dto

import "time"

type ResponseLov struct {
	Id    uint   `json:"id"`
	Value string `json:"value"`
}

type ResponseFile struct {
	Metadata map[string]interface{} `json:"metadata"`
	FilePath []string               `json:"file_paths"`
}

type ResponseHeaderDashboardTinjauan struct {
	Total     any                 `json:"total"`
	Percent   float64             `json:"percent"`
	Trend     string              `json:"trend"`
	Sparkline []ResponseSparkline `json:"sparkline"`
}

type ResponseHeaderDashboardIklan struct {
	Total     any                 `json:"total"`
	Percent   float64             `json:"percent"`
	Trend     string              `json:"trend"`
	Sparkline []ResponseSparkline `json:"sparkline"`
}

type ResponseSparkline struct {
	Tanggal time.Time `json:"tanggal"`
	Total   float64   `json:"total"`
}

type ResponseTrenPenjualanDashboardTinjauan struct {
	Total float64 `json:"total"`
	Month string  `json:"month"`
}

type ResponseTotalPenjualanInWeekDashboardTinjauan struct {
	Total float64 `json:"total"`
	Day   string  `json:"day"`
}

type ResponseTotalPenjualanDanBiayaDashboardIklan struct {
	Biaya     float64 `json:"biaya"`
	Penjualan float64 `json:"penjualan"`
	Month     string  `json:"month"`
}

type ResponseTotalROASDashboardIklan struct {
	Roas  float64 `json:"roas"`
	Month string  `json:"month"`
}

type ResponseTopProductDashboardIklan struct {
	NamaIklan      string  `json:"nama_iklan"`
	Roas           float64 `json:"roas"`
	ConvertionRate float64 `json:"convertion_rate"`
	Biaya          float64 `json:"biaya"`
	Penjualan      float64 `json:"penjualan"`
}
