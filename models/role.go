package models

type Role struct {
	BaseModel
	Name      string     `json:"name"`
	TenantID  *uint      `gorm:"index" json:"tenant_id"`
	RoleType  string     `json:"role_type" binding:"required,oneof=GLOBAL_SUPERADMIN TENANT_SUPERADMIN TENANT_ADMIN"`
	IsActive  bool       `gorm:"type:bool;default:true" json:"is_active"`
	Users     []User     `json:"users" gorm:"foreignKey:RoleID"`
	RoleMenus []RoleMenu `json:"role_menus" gorm:"foreignKey:RoleID"` // ⬅️ tambahkan ini
}
