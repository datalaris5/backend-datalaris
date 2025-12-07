package models

import (
	"go-datalaris/common"
)

type HistoryDataUpload struct {
	common.BaseModel
	StoreId  uint   `json:"store_id"`
	Filename string `json:"filename"`
	Status   string `json:"status"`
}
