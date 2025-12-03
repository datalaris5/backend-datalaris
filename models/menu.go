package models

type Menu struct {
	BaseModel
	Name     string `json:"name"`
	Path     string `json:"path"`
	TenantID uint   `gorm:"index;not null" json:"tenant_id"`
}
