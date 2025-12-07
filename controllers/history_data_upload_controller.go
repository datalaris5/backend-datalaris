package controllers

import (
	"go-datalaris/config"
	"go-datalaris/constant"
	"go-datalaris/models"
	"go-datalaris/utils"

	"github.com/gin-gonic/gin"
)

func GetHistoryDataUpload(c *gin.Context) {
	tenantID := utils.GetTenantId(c)
	var result []models.HistoryDataUpload
	db := config.DB
	db.Raw(`
	SELECT h.*
	FROM history_data_uploads h
	JOIN stores s on s.id = h.store_id
	JOIN tenants t on t.id = s.tenant_id
	WHERE t.ID = ? AND h.created_at >= CURRENT_DATE - INTERVAL '30 days'
	ORDER BY h.created_at DESC;`, tenantID).Scan(&result)

	utils.Success(c, constant.HistoryDataUpload+constant.SuccessFetch, result)
}
