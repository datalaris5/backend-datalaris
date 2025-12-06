package models

type Store struct {
	BaseModel
	Name          string `json:"name"`
	TenantID      uint   `json:"tenant_id"`
	MarketplaceID uint   `json:"marketplace_id"`
	IsActive      bool   `gorm:"type:bool;default:true" json:"is_active"`
	IsDeleted     bool   `gorm:"type:bool;default:false" json:"is_deleted"`
}
