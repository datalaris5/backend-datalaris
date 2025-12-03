package utils

import "go-datalaris/dto"

func PaginationResponse(data interface{}, total int64, page, limit int) dto.PaginationResponse {
	return dto.PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int((total + int64(limit) - 1) / int64(limit)),
	}
}
