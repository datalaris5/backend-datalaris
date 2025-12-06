package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RequestRole struct {
	ID       uint    `json:"id"`
	TenantID *uint   `json:"tenant_id"`
	Status   *string `json:"status"`
}

type RequestUser struct {
	ID       uint    `json:"id"`
	TenantID *uint   `json:"tenant_id"`
	Status   *string `json:"status"`
}

type RequestTenantId struct {
	TenantID *uint `json:"tenant_id"`
}

type RequestDashboardTinjauan struct {
	StoreId  uint   `json:"store_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
