package models

type Marketplace struct {
	BaseModel
	Name      string `json:"name"`
	IsActive  bool   `gorm:"type:bool;default:true" json:"is_active"`
	IsDeleted bool   `gorm:"type:bool;default:false" json:"is_deleted"`
}
