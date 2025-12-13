package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required"`
	Password   string `json:"password" binding:"required"`
	TenantName string `json:"tenant_name"`
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

type RequestDashboardIklan struct {
	StoreId  uint   `json:"store_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type RequestDashboardChat struct {
	StoreId  uint   `json:"store_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
