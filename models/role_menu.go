package models

type RoleMenu struct {
	BaseModel
	RoleID    uint `json:"role_id"`
	MenuID    uint `json:"menu_id"`
	CanCreate bool `json:"can_create"`
	CanRead   bool `json:"can_read"`
	CanUpdate bool `json:"can_update"`
	CanDelete bool `json:"can_delete"`

	Role Role `json:"-"`
	Menu Menu `json:"-"`
}
