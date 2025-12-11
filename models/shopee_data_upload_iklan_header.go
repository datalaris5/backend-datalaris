package models

import (
	"go-datalaris/common"
	"time"
)

type ShopeeDataUploadIklanHeader struct {
	common.BaseModel
	StoreID  uint      `gorm:"column:store_id"`
	Filename string    `gorm:"column:filename"`
	DateFrom time.Time `gorm:"column:date_from"`
	DateTo   time.Time `gorm:"column:date_to"`
}

func (ShopeeDataUploadIklanHeader) TableName() string {
	return "shopee_data_upload_iklan_headers"
}
