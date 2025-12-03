package models

type Tenant struct {
	BaseModel
	Name        string `gorm:"unique;not null" json:"name"`
	Description string `json:"description"`
	TenantKey   string `json:"tenant_key"`
}
