package models

type User struct {
	BaseModel
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"-"`
	TenantID *uint   `json:"tenant_id"`
	RoleID   *uint   `json:"role_id"`
	Role     *Role   `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Tenant   *Tenant `json:"tenant,omitempty"`
	IsActive bool    `gorm:"type:bool;default:true" json:"is_active"`
}
