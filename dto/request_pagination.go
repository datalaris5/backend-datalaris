package dto

type PaginationRequest struct {
	Page     int     `json:"page" binding:"required"`
	Limit    int     `json:"limit" binding:"required"`
	Search   *string `json:"search"`
	TenantID *uint   `json:"tenant_id"`
}
