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
